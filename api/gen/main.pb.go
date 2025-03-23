// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: api/proto/main.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AskRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Text          string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AskRequest) Reset() {
	*x = AskRequest{}
	mi := &file_api_proto_main_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AskRequest) ProtoMessage() {}

func (x *AskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_main_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AskRequest.ProtoReflect.Descriptor instead.
func (*AskRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_main_proto_rawDescGZIP(), []int{0}
}

func (x *AskRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AskRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type AskResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 因为是流式的 需要判断是否是终止的
	Final bool `protobuf:"varint,1,opt,name=final,proto3" json:"final,omitempty"`
	// 返回的 结果
	Text          string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AskResponse) Reset() {
	*x = AskResponse{}
	mi := &file_api_proto_main_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AskResponse) ProtoMessage() {}

func (x *AskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_main_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AskResponse.ProtoReflect.Descriptor instead.
func (*AskResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_main_proto_rawDescGZIP(), []int{1}
}

func (x *AskResponse) GetFinal() bool {
	if x != nil {
		return x.Final
	}
	return false
}

func (x *AskResponse) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_api_proto_main_proto protoreflect.FileDescriptor

var file_api_proto_main_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x30, 0x0a,
	0x0a, 0x41, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22,
	0x37, 0x0a, 0x0b, 0x41, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x66, 0x69, 0x6e, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x66,
	0x69, 0x6e, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x32, 0x3b, 0x0a, 0x09, 0x41, 0x49, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x03, 0x41, 0x73, 0x6b, 0x12, 0x11, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x12, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x41, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x5e, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x42, 0x09, 0x4d, 0x61, 0x69, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x12, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x64, 0x65, 0x65, 0x70, 0x73,
	0x65, 0x65, 0x6b, 0xa2, 0x02, 0x03, 0x4d, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0xca, 0x02, 0x05, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0xe2, 0x02, 0x11, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x05,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_proto_main_proto_rawDescOnce sync.Once
	file_api_proto_main_proto_rawDescData []byte
)

func file_api_proto_main_proto_rawDescGZIP() []byte {
	file_api_proto_main_proto_rawDescOnce.Do(func() {
		file_api_proto_main_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_main_proto_rawDesc), len(file_api_proto_main_proto_rawDesc)))
	})
	return file_api_proto_main_proto_rawDescData
}

var file_api_proto_main_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_proto_main_proto_goTypes = []any{
	(*AskRequest)(nil),  // 0: model.AskRequest
	(*AskResponse)(nil), // 1: model.AskResponse
}
var file_api_proto_main_proto_depIdxs = []int32{
	0, // 0: model.AIService.Ask:input_type -> model.AskRequest
	1, // 1: model.AIService.Ask:output_type -> model.AskResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_main_proto_init() }
func file_api_proto_main_proto_init() {
	if File_api_proto_main_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_main_proto_rawDesc), len(file_api_proto_main_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_main_proto_goTypes,
		DependencyIndexes: file_api_proto_main_proto_depIdxs,
		MessageInfos:      file_api_proto_main_proto_msgTypes,
	}.Build()
	File_api_proto_main_proto = out.File
	file_api_proto_main_proto_goTypes = nil
	file_api_proto_main_proto_depIdxs = nil
}
