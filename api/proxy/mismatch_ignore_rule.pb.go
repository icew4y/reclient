// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.26.0
// source: api/proxy/mismatch_ignore_rule.proto

package proxy

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MismatchIgnoreConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []*Rule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *MismatchIgnoreConfig) Reset() {
	*x = MismatchIgnoreConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MismatchIgnoreConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MismatchIgnoreConfig) ProtoMessage() {}

func (x *MismatchIgnoreConfig) ProtoReflect() protoreflect.Message {
	mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MismatchIgnoreConfig.ProtoReflect.Descriptor instead.
func (*MismatchIgnoreConfig) Descriptor() ([]byte, []int) {
	return file_api_proxy_mismatch_ignore_rule_proto_rawDescGZIP(), []int{0}
}

func (x *MismatchIgnoreConfig) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to RuleSpec:
	//
	//	*Rule_OutputFilePathRuleSpec
	RuleSpec isRule_RuleSpec `protobuf_oneof:"rule_spec"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_api_proxy_mismatch_ignore_rule_proto_rawDescGZIP(), []int{1}
}

func (m *Rule) GetRuleSpec() isRule_RuleSpec {
	if m != nil {
		return m.RuleSpec
	}
	return nil
}

func (x *Rule) GetOutputFilePathRuleSpec() *OutputFilePathRuleSpec {
	if x, ok := x.GetRuleSpec().(*Rule_OutputFilePathRuleSpec); ok {
		return x.OutputFilePathRuleSpec
	}
	return nil
}

type isRule_RuleSpec interface {
	isRule_RuleSpec()
}

type Rule_OutputFilePathRuleSpec struct {
	OutputFilePathRuleSpec *OutputFilePathRuleSpec `protobuf:"bytes,1,opt,name=output_file_path_rule_spec,json=outputFilePathRuleSpec,proto3,oneof"`
}

func (*Rule_OutputFilePathRuleSpec) isRule_RuleSpec() {}

type OutputFilePathRuleSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PathPattern *RegexPattern `protobuf:"bytes,1,opt,name=path_pattern,json=pathPattern,proto3" json:"path_pattern,omitempty"`
}

func (x *OutputFilePathRuleSpec) Reset() {
	*x = OutputFilePathRuleSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutputFilePathRuleSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutputFilePathRuleSpec) ProtoMessage() {}

func (x *OutputFilePathRuleSpec) ProtoReflect() protoreflect.Message {
	mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutputFilePathRuleSpec.ProtoReflect.Descriptor instead.
func (*OutputFilePathRuleSpec) Descriptor() ([]byte, []int) {
	return file_api_proxy_mismatch_ignore_rule_proto_rawDescGZIP(), []int{2}
}

func (x *OutputFilePathRuleSpec) GetPathPattern() *RegexPattern {
	if x != nil {
		return x.PathPattern
	}
	return nil
}

type RegexPattern struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	Inverted   bool   `protobuf:"varint,2,opt,name=inverted,proto3" json:"inverted,omitempty"`
}

func (x *RegexPattern) Reset() {
	*x = RegexPattern{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegexPattern) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegexPattern) ProtoMessage() {}

func (x *RegexPattern) ProtoReflect() protoreflect.Message {
	mi := &file_api_proxy_mismatch_ignore_rule_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegexPattern.ProtoReflect.Descriptor instead.
func (*RegexPattern) Descriptor() ([]byte, []int) {
	return file_api_proxy_mismatch_ignore_rule_proto_rawDescGZIP(), []int{3}
}

func (x *RegexPattern) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *RegexPattern) GetInverted() bool {
	if x != nil {
		return x.Inverted
	}
	return false
}

var File_api_proxy_mismatch_ignore_rule_proto protoreflect.FileDescriptor

var file_api_proxy_mismatch_ignore_rule_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x6d, 0x69, 0x73, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x5f, 0x72, 0x75, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x22, 0x39, 0x0a,
	0x14, 0x4d, 0x69, 0x73, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x49, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x21, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x52, 0x75, 0x6c,
	0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x22, 0x70, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65,
	0x12, 0x5b, 0x0a, 0x1a, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x70, 0x61, 0x74, 0x68, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65, 0x53,
	0x70, 0x65, 0x63, 0x48, 0x00, 0x52, 0x16, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x50, 0x61, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65, 0x53, 0x70, 0x65, 0x63, 0x42, 0x0b, 0x0a,
	0x09, 0x72, 0x75, 0x6c, 0x65, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x22, 0x50, 0x0a, 0x16, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x52, 0x75, 0x6c, 0x65,
	0x53, 0x70, 0x65, 0x63, 0x12, 0x36, 0x0a, 0x0c, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x70, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x2e, 0x52, 0x65, 0x67, 0x65, 0x78, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x52,
	0x0b, 0x70, 0x61, 0x74, 0x68, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x22, 0x4a, 0x0a, 0x0c,
	0x52, 0x65, 0x67, 0x65, 0x78, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x12, 0x1e, 0x0a, 0x0a,
	0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x69, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x69, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x65, 0x64, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x7a, 0x65, 0x6c, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x2f, 0x72, 0x65, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x6f, 0x78, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proxy_mismatch_ignore_rule_proto_rawDescOnce sync.Once
	file_api_proxy_mismatch_ignore_rule_proto_rawDescData = file_api_proxy_mismatch_ignore_rule_proto_rawDesc
)

func file_api_proxy_mismatch_ignore_rule_proto_rawDescGZIP() []byte {
	file_api_proxy_mismatch_ignore_rule_proto_rawDescOnce.Do(func() {
		file_api_proxy_mismatch_ignore_rule_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proxy_mismatch_ignore_rule_proto_rawDescData)
	})
	return file_api_proxy_mismatch_ignore_rule_proto_rawDescData
}

var file_api_proxy_mismatch_ignore_rule_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_proxy_mismatch_ignore_rule_proto_goTypes = []any{
	(*MismatchIgnoreConfig)(nil),   // 0: proxy.MismatchIgnoreConfig
	(*Rule)(nil),                   // 1: proxy.Rule
	(*OutputFilePathRuleSpec)(nil), // 2: proxy.OutputFilePathRuleSpec
	(*RegexPattern)(nil),           // 3: proxy.RegexPattern
}
var file_api_proxy_mismatch_ignore_rule_proto_depIdxs = []int32{
	1, // 0: proxy.MismatchIgnoreConfig.rules:type_name -> proxy.Rule
	2, // 1: proxy.Rule.output_file_path_rule_spec:type_name -> proxy.OutputFilePathRuleSpec
	3, // 2: proxy.OutputFilePathRuleSpec.path_pattern:type_name -> proxy.RegexPattern
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_proxy_mismatch_ignore_rule_proto_init() }
func file_api_proxy_mismatch_ignore_rule_proto_init() {
	if File_api_proxy_mismatch_ignore_rule_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_proxy_mismatch_ignore_rule_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*MismatchIgnoreConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_proxy_mismatch_ignore_rule_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Rule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_proxy_mismatch_ignore_rule_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*OutputFilePathRuleSpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_proxy_mismatch_ignore_rule_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RegexPattern); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_api_proxy_mismatch_ignore_rule_proto_msgTypes[1].OneofWrappers = []any{
		(*Rule_OutputFilePathRuleSpec)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proxy_mismatch_ignore_rule_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_proxy_mismatch_ignore_rule_proto_goTypes,
		DependencyIndexes: file_api_proxy_mismatch_ignore_rule_proto_depIdxs,
		MessageInfos:      file_api_proxy_mismatch_ignore_rule_proto_msgTypes,
	}.Build()
	File_api_proxy_mismatch_ignore_rule_proto = out.File
	file_api_proxy_mismatch_ignore_rule_proto_rawDesc = nil
	file_api_proxy_mismatch_ignore_rule_proto_goTypes = nil
	file_api_proxy_mismatch_ignore_rule_proto_depIdxs = nil
}
