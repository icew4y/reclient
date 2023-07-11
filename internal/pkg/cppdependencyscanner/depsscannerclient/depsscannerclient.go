// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package depsscannerclient implements the cppdependencyscanner.DepsScanner with gRPC
// to perform the same tasks on a remote service instead of internally with cgo.
package depsscannerclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bazelbuild/remote-apis-sdks/go/pkg/command"
	"github.com/bazelbuild/remote-apis-sdks/go/pkg/outerr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"team/foundry-x/re-client/internal/pkg/ipc"

	pb "team/foundry-x/re-client/api/cppscandeps"

	log "github.com/golang/glog"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Executor can run commands and retrieve their outputs.
type executor interface {
	ExecuteInBackground(ctx context.Context, cmd *command.Command, oe outerr.OutErr, ch chan *command.Result) error
}

// Implements the outerr.OutErr interface to log stdout and stderr from a background process.
// Note that the dependency scanner service has its own logging and this output is expected to be
// empty except when a failure occurs.
type depsScannerLogger struct{}

// Log stdout lines from dependency scanner service to reproxy.INFO
func (l *depsScannerLogger) WriteOut(s []byte) {
	log.Infof("depsscannerclient stdout: %s", string(s[:]))
}

// Log stderr lines from dependency scanner service to reproxy.ERROR
func (l *depsScannerLogger) WriteErr(s []byte) {
	log.Errorf("depsscannerclient stderr: %s", string(s[:]))
}

const (
	// gRPC C++ library does not support named pipes on Windows.
	// Named pipes are supported on Linux and Mac but for consistency we will use TCP sockets for all.
	localhost = "127.0.0.1"
	// 60 seconds to try and shut down gracefully, 10 seconds to hard kill
	shutdownTimeout = 60 * time.Second
)

// DepsScannerClient wraps the dependency scanner gRPC client.
type DepsScannerClient struct {
	ctx               context.Context
	terminate         context.CancelFunc
	address           string
	executable        string
	cacheDir          string
	ignoredPluginsMap map[string]bool
	cacheFileMaxMb    int
	useDepsCache      bool
	logDir            string
	client            pb.CPPDepsScannerClient
	executor          executor
	oe                outerr.OutErr
	ch                chan *command.Result
	serviceRestarted  time.Time
	m                 sync.Mutex
}

var connect = func(ctx context.Context, address string) (pb.CPPDepsScannerClient, error) {
	conn, err := ipc.DialContextWithBlock(ctx, address)
	if err != nil {
		return nil, err
	}

	client := pb.NewCPPDepsScannerClient(conn)
	return client, nil
}

// TODO (b/258275137): make this configurable and move somewhere more appropriate when reconnect logic is implemented.
var connTimeout = 30 * time.Second

// New creates new DepsScannerClient.
func New(ctx context.Context, executor executor, cacheDir string, cacheFileMaxMb int, ignoredPlugins []string, useDepsCache bool, logDir string, depsScannerAddress, proxyServerAddress string) (*DepsScannerClient, error) {
	log.Infof("Connecting to remote dependency scanner: %v", depsScannerAddress)

	ignoredPluginsMap := map[string]bool{}
	for _, plugin := range ignoredPlugins {
		ignoredPluginsMap[plugin] = true
	}
	client := &DepsScannerClient{
		address:           depsScannerAddress,
		executor:          executor,
		cacheDir:          cacheDir,
		logDir:            logDir,
		ignoredPluginsMap: ignoredPluginsMap,
		cacheFileMaxMb:    cacheFileMaxMb,
		useDepsCache:      useDepsCache,
		// TODO (b/260707840): context shouldn't be a member variable. Pass in as function variable elsewhere and remote this.
		ctx: ctx,
	}

	if strings.HasPrefix(depsScannerAddress, "exec://") {
		executable := depsScannerAddress[7:]
		addr, err := buildAddress(proxyServerAddress, findOpenPort)
		if err != nil {
			return nil, fmt.Errorf("Failed to build address for dependency scanner: %w", err)
		}
		client.address = addr
		if err := client.startService(ctx, executable); err != nil {
			return nil, fmt.Errorf("Failed to start dependency scanner: %w", err)
		}
	}

	connTimeoutCtx, _ := context.WithTimeout(ctx, connTimeout)
	for {
		select {
		case <-connTimeoutCtx.Done():
			client.Close()
			return nil, fmt.Errorf("Failed to connect to dependency scanner service after %v seconds", connTimeout.Seconds())
		case err := <-client.ch:
			if err != nil {
				return nil, fmt.Errorf("%v terminated during startup: %w", client.executable, err)
			}
			continue
		default:
			tctx, _ := context.WithTimeout(connTimeoutCtx, 50*time.Millisecond)
			c, err := connect(tctx, client.address)
			if err != nil {
				log.Infof("Failed to connect to dependency scanner, it may not be started yet: %v", err)
				continue
			}
			log.Infof("Connected to dependency scanner service on %v", client.address)
			client.client = c
			return client, nil
		}
	}
}

// buildAddress generates an address for the depsscanner process to listen on.
// If reproxy is on UNIX and using a UDS address then it will append .despscan.sock
// to that address. Otherwise a random TCP port will be chosen.
func buildAddress(proxyServerAddress string, openPortFunc func() (int, error)) (string, error) {
	address := proxyServerAddress
	if ipc.GrpcCxxSupportsUDS && strings.HasPrefix(address, "unix://") {
		return address + ".despscan.sock", nil
	}
	if strings.HasPrefix(address, "unix://") || strings.HasPrefix(address, "pipe://") {
		address = fmt.Sprintf("%s:0", localhost)
	}
	base, _, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("failed to find base address: %w", err)
	}
	port, err := openPortFunc()
	if err != nil {
		return "", fmt.Errorf("failed to find open port: %w", err)
	}
	return fmt.Sprintf("%s:%d", base, port), nil
}

// findOpenPort finds an open port by resolving 127.0.0.1:0 which the kernel resolves
// to an unused free port. It then opens a tcp server and client on that port and sends 1
// byte before closing to ensure the kernel sets the port to TIME-WAIT to prevent non
// subprocesses from listening on this port for 60s.
// Inspired by https://github.com/Yelp/ephemeral-port-reserve
func findOpenPort() (int, error) {
	// Listen on port 0 because by convention kernels will resolve it to an unsused free port
	// see https://linux.die.net/man/7/ip
	// and https://learn.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-bind
	// Other systems may have different behaviour
	// (https://daniel.haxx.se/blog/2014/10/25/pretending-port-zero-is-a-normal-one/)
	// in that case a uds socket should be used since only windows does not support them.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(localhost)})
	if err != nil {
		return 0, err
	}
	// Ignore error as l will be already closed unless there is an error before l.Close() below.
	defer l.Close()
	// Wait a maximum of 1 second for the byte to be sent on the acquired port
	l.SetDeadline(time.Now().Add(time.Second))
	resolvedAddr, ok := l.Addr().(*net.TCPAddr)
	if !ok || resolvedAddr == nil {
		return 0, fmt.Errorf("Failed to resolve %s:0 to an open port", localhost)
	}
	lAddr := resolvedAddr.String()
	port := resolvedAddr.Port
	if port == 0 {
		return 0, fmt.Errorf("Kernel did not resolve %s:0 to an open port", localhost)
	}
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer close(errCh)
		defer func() {
			errCh <- err
		}()
		// Accept 1 connection from the socket.
		var ac net.Conn
		ac, err = l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				err = nil
			}
			return
		}
		// Ignore error as ac will be already closed unless there is an
		// error before ac.Close() below.
		defer ac.Close()
		// Read from the connection until it's closed by the remote peer.
		if _, err = io.ReadAll(ac); err != nil {
			return
		}
		// Close the socket.
		if err = ac.Close(); err != nil {
			return
		}
	}()
	// Dial the socket that we just opened.
	c, err := net.Dial("tcp", lAddr)
	if err != nil {
		return 0, err
	}
	// Ignore error as c will be already closed unless there is an error before c.Close() below.
	defer c.Close()
	// Write some arbitrary bytes to it.
	// If the port is not written to then the kernel will close and free the port immediately.
	if _, err := c.Write([]byte("x")); err != nil {
		return 0, err
	}
	// Close the connection to the socket now that we're done.
	if err := c.Close(); err != nil {
		return 0, err
	}
	// Close the listener now that we know we're done sending anything to it.
	if err := l.Close(); err != nil {
		return 0, err
	}
	if err := <-errCh; err != nil {
		return 0, err
	}
	// The port should be in TIME-WAIT state now
	return port, nil
}

// Close implements DepsScanner.Close.
// It cleanly disconnects from the remote service and releases resource associated with
// DepsScannerClient.
func (ds *DepsScannerClient) Close() {
	if ds.client == nil {
		return
	}
	if err := ds.stopService(shutdownTimeout); err != nil {
		log.Errorf("%v", err)
	}
	ds.client = nil
}

// ProcessInputs implements DepsScanner.ProcessInputs by sending a ProcessInputs gRPC to the
// connected server and returns the result.
// Returns list of dependencies, boolean indicating whether deps cache was used, and
// error if there was one.
func (ds *DepsScannerClient) ProcessInputs(ctx context.Context, execID string, compileCommand []string, filename, directory string, cmdEnv []string) (dependencies []string, usedCache bool, err error) {
	log.V(3).Infof("%v: Started remote input processing for %v", execID, filename)
	resCh := make(chan *pb.CPPProcessInputsResponse)
	errCh := make(chan error)
	go func() {
		resp, err := ds.client.ProcessInputs(
			ctx,
			&pb.CPPProcessInputsRequest{
				ExecId:    execID,
				Command:   compileCommand,
				Directory: directory,
				CmdEnv:    cmdEnv,
				Filename:  filename,
			})
		log.V(3).Infof("ProcessInputs complete: %v", resp)
		if err != nil {
			errCh <- err
		} else {
			resCh <- resp
		}
	}()

	select {
	case <-ctx.Done():
		// Timeout processing inputs.  Could be the service is offline.  Restart it if possible.
		err := ds.restartService(ds.ctx, ds.executable)
		if err == nil {
			// Successfully restarted service; bubble up DeadlineExceeded to trigger a retry
			return nil, false, fmt.Errorf("failed to get response from scandeps: %w",
				ctx.Err())
		}
		// else unable to restart the service, or reproxy is not responsible for the service
		return nil, false, fmt.Errorf("failed to get response from scandeps; additionally failed to restart service: %w", err)

	case err := <-errCh:
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
			// Unavailable means a disconnect has occurred.
			if restartErr := ds.restartService(ds.ctx, ds.executable); restartErr != nil {
				return nil, false, fmt.Errorf("communication with service lost; failed to restart service: %w", restartErr)
			}
			return nil, false, errors.New("communication with service lost; service restarted")
		}
		// else something unexpected has gone wrong.
		return nil, false, fmt.Errorf("An unexpected error occurred communicating with the service: %w", err)
	case resp := <-resCh:
		if resp.Error != "" {
			return nil, false, fmt.Errorf("input processing failed: %v", resp.Error)
		}
		absDeps := make([]string, 0, len(resp.Dependencies))

		for _, p := range resp.Dependencies {
			if p == "" {
				continue
			}
			if !filepath.IsAbs(p) {
				p = filepath.Join(directory, p)
			}
			absDeps = append(absDeps, p)
		}
		return absDeps, resp.UsedCache, nil
	}
}

// ShouldIgnorePlugin implements DepsScanner.ShouldIgnorePlugin.
func (ds *DepsScannerClient) ShouldIgnorePlugin(plugin string) bool {
	_, present := ds.ignoredPluginsMap[plugin]
	return present
}

func (ds *DepsScannerClient) verifyService(ctx context.Context) error {
	retries := 10
	timeout := 10 * time.Second
	for i := 0; i < retries; i++ {
		sctx, cancel := context.WithTimeout(ds.ctx, timeout)
		defer cancel()

		_, err := ds.client.Status(sctx, &emptypb.Empty{})

		select {
		case <-ctx.Done():
			// timeout, retry
			continue
		default:
			// success?
			if err == nil {
				return nil
			} // else
			// Status call may return an error before the 10 seconds timeout expires if it isn't
			// ready to accept connections yet.  Typically it will error instantly with no delay.
			// In that case we want a delay (up to timeout) to give it more time to be ready.
			time.Sleep(timeout)
		}
	}
	// Still haven't connected; give up
	return fmt.Errorf("Unable to connect to server after %v seconds", retries*(int)(timeout.Seconds()))
}

func (ds *DepsScannerClient) restartService(ctx context.Context, executable string) error {
	if executable == "" {
		// Not responsible for service
		return fmt.Errorf("Service is not managed by reproxy")
	}

	t := time.Now()
	ds.m.Lock()
	defer ds.m.Unlock()
	if t.Before(ds.serviceRestarted) {
		// service has been restarted since this thread was paused
		return nil
	}
	if err := ds.stopService(shutdownTimeout); err != nil {
		log.Errorf("%v", err)
		return fmt.Errorf("Unable to shutdown service: %v", err)
	}
	if err := ds.startService(ctx, executable); err != nil {
		log.Errorf("Failed to start dependency scanner: %v", err)
		return fmt.Errorf("Unable to start service: %v", err)
	}

	err := ds.verifyService(ctx)
	if err == nil {
		ds.serviceRestarted = time.Now()
	}
	return err
}

func (ds *DepsScannerClient) startService(ctx context.Context, executable string) error {
	ctx, ds.terminate = context.WithCancel(ctx)
	ds.executable = executable

	cmdArgs := []string{ds.executable, "--server_address", ds.address}
	cmdArgs = append(cmdArgs, "--cache_dir", ds.cacheDir)
	cmdArgs = append(cmdArgs, "--deps_cache_max_mb", strconv.FormatInt(int64(ds.cacheFileMaxMb), 10))
	if ds.useDepsCache {
		cmdArgs = append(cmdArgs, "--enable_deps_cache")
	} else {
		cmdArgs = append(cmdArgs, "--noenable_deps_cache")
	}

	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		// Debugging parameters `experimental_segfault` and `experimental_deadlock` can be used
		// by setting their corresponding `FLAGS_*` environment variables.
		key, val, err := findKeyVal(e, os.LookupEnv)
		if err != nil {
			return err
		}
		switch key {
		case "FLAGS_experimental_segfault":
			cmdArgs = append(cmdArgs, "--experimental_segfault", val)
		case "FLAGS_experimental_deadlock":
			cmdArgs = append(cmdArgs, "--experimental_deadlock", val)
		default:
			envVars[key] = val
		}
	}
	if ds.logDir != "" {
		log.Infof("Setting GLOG_log_dir=\"%v\"", ds.logDir)
		envVars["GLOG_log_dir"] = ds.logDir
	}

	log.Infof("Starting service: %v", cmdArgs)
	cmd := &command.Command{Args: cmdArgs}
	cmd.InputSpec = &command.InputSpec{
		EnvironmentVariables: envVars,
	}

	ds.oe = &depsScannerLogger{}
	ds.ch = make(chan *command.Result)

	return ds.executor.ExecuteInBackground(ctx, cmd, ds.oe, ds.ch)
}

func findKeyVal(envStr string, lookupEnvFunc func(string) (string, bool)) (string, string, error) {
	envParts := strings.Split(envStr, "=")
	if len(envParts) == 2 {
		return envParts[0], envParts[1], nil
	}
	for i := 1; i < len(envParts)+1; i++ {
		pkey := strings.Join(envParts[:i], "=")
		pval := strings.Join(envParts[i:], "=")
		if val, ok := lookupEnvFunc(pkey); ok && val == pval {
			return pkey, pval, nil
		}
	}
	return "", "", fmt.Errorf("Got %s in env vars list but could not find any matching env var", strings.Join(envParts, "="))
}

// stopService attempts to stop the dependency scanner service started by reproxy.
// If process cannot be stopped within `timeout` it will be killed forcefully.
// If process is still running after a second `timeout` passed, an error will be returned.
// Note the total time before returning may be up to `2*timeout`.
func (ds *DepsScannerClient) stopService(timeout time.Duration) error {
	if ds.executable == "" {
		// Dependency scanner service wasn't started by reclient; nothing to do
		return nil
	}

	// Dependency scanner service was started by reproxy; we must stop it
	ctx, cancel := context.WithTimeout(ds.ctx, 10*time.Second)
	defer cancel()
	if ds.client != nil {
		statusResponse, err := ds.client.Shutdown(ctx, &emptypb.Empty{})
		if err != nil {
			// This could mean the service has crashed, it could be frozen, or may have taken
			// too long to respond. Log the error and wait to see if it shuts itself down.
			log.Errorf("Error sending shutdown command to %v: %v", ds.executable, err)
		} else {
			log.Infof("Shutdown response from %s (v%s)", statusResponse.GetName(), statusResponse.GetVersion())
			if statusResponse.GetUptime() != nil {
				log.Infof("> Uptime: %d seconds", statusResponse.GetUptime().GetSeconds())
			}
			log.Infof("> Completed: %d", statusResponse.GetCompletedActions())
			log.Infof("> Still in Progress: %d", statusResponse.GetRunningActions())
		}
	} else {
		// We started but were never able to connect to the service; must force kill the service.
		// This is a noop if the service crashed unexpectedly.
		ds.terminate()
	}

	for i := 0; i < 2; i++ {
		select {
		case done := <-ds.ch:
			if done.IsOk() {
				log.Infof("%v successfully stopped", ds.executable)
			} else {
				log.Warningf("%v stopped with error code %v", ds.executable, done.ExitCode)
			}
			// service has been successfully stopped
			return nil
		case <-time.After(timeout):
			log.Warningf("Still waiting for shutdown after %.0f seconds", timeout.Seconds())
			if i == 0 {
				log.Info("Sending termination signal to %v", ds.executable)
				// We've waited hardkill seconds for the shutdown to complete.
				// Assume that it's stuck and perform a hard kill.
				// This cancel is tied to the command context and will kill the subprocess and all
				// active ProcessInputs requests
				ds.terminate()
			}
		}
	}
	return fmt.Errorf("could not shutdown %v after %.0f seconds. Giving up", ds.executable, 2*(timeout.Seconds()))
}
