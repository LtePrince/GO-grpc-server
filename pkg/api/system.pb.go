// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: pkg/api/system.proto

package api

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

type SendFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FilePath      string                 `protobuf:"bytes,1,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"` // 要传输的文件路径
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendFileRequest) Reset() {
	*x = SendFileRequest{}
	mi := &file_pkg_api_system_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFileRequest) ProtoMessage() {}

func (x *SendFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_system_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFileRequest.ProtoReflect.Descriptor instead.
func (*SendFileRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_system_proto_rawDescGZIP(), []int{0}
}

func (x *SendFileRequest) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

type SendFileResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*SendFileResponse_Metadata
	//	*SendFileResponse_Chunk
	Data          isSendFileResponse_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendFileResponse) Reset() {
	*x = SendFileResponse{}
	mi := &file_pkg_api_system_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFileResponse) ProtoMessage() {}

func (x *SendFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_system_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFileResponse.ProtoReflect.Descriptor instead.
func (*SendFileResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_system_proto_rawDescGZIP(), []int{1}
}

func (x *SendFileResponse) GetData() isSendFileResponse_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *SendFileResponse) GetMetadata() *FileMetadata {
	if x != nil {
		if x, ok := x.Data.(*SendFileResponse_Metadata); ok {
			return x.Metadata
		}
	}
	return nil
}

func (x *SendFileResponse) GetChunk() []byte {
	if x != nil {
		if x, ok := x.Data.(*SendFileResponse_Chunk); ok {
			return x.Chunk
		}
	}
	return nil
}

type isSendFileResponse_Data interface {
	isSendFileResponse_Data()
}

type SendFileResponse_Metadata struct {
	Metadata *FileMetadata `protobuf:"bytes,1,opt,name=metadata,proto3,oneof"` // 首次响应携带元数据
}

type SendFileResponse_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"` // 后续响应携带数据块
}

func (*SendFileResponse_Metadata) isSendFileResponse_Data() {}

func (*SendFileResponse_Chunk) isSendFileResponse_Data() {}

type FileMetadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileName      string                 `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	MimeType      string                 `protobuf:"bytes,2,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
	FileSize      uint64                 `protobuf:"varint,3,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileMetadata) Reset() {
	*x = FileMetadata{}
	mi := &file_pkg_api_system_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMetadata) ProtoMessage() {}

func (x *FileMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_system_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMetadata.ProtoReflect.Descriptor instead.
func (*FileMetadata) Descriptor() ([]byte, []int) {
	return file_pkg_api_system_proto_rawDescGZIP(), []int{2}
}

func (x *FileMetadata) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileMetadata) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

func (x *FileMetadata) GetFileSize() uint64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

var File_pkg_api_system_proto protoreflect.FileDescriptor

const file_pkg_api_system_proto_rawDesc = "" +
	"\n" +
	"\x14pkg/api/system.proto\x12\x06system\".\n" +
	"\x0fSendFileRequest\x12\x1b\n" +
	"\tfile_path\x18\x01 \x01(\tR\bfilePath\"f\n" +
	"\x10SendFileResponse\x122\n" +
	"\bmetadata\x18\x01 \x01(\v2\x14.system.FileMetadataH\x00R\bmetadata\x12\x16\n" +
	"\x05chunk\x18\x02 \x01(\fH\x00R\x05chunkB\x06\n" +
	"\x04data\"e\n" +
	"\fFileMetadata\x12\x1b\n" +
	"\tfile_name\x18\x01 \x01(\tR\bfileName\x12\x1b\n" +
	"\tmime_type\x18\x02 \x01(\tR\bmimeType\x12\x1b\n" +
	"\tfile_size\x18\x03 \x01(\x04R\bfileSize2P\n" +
	"\rSystemService\x12?\n" +
	"\bSendFile\x12\x17.system.SendFileRequest\x1a\x18.system.SendFileResponse0\x01B\x0fZ\r/pkg/api/;apib\x06proto3"

var (
	file_pkg_api_system_proto_rawDescOnce sync.Once
	file_pkg_api_system_proto_rawDescData []byte
)

func file_pkg_api_system_proto_rawDescGZIP() []byte {
	file_pkg_api_system_proto_rawDescOnce.Do(func() {
		file_pkg_api_system_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_pkg_api_system_proto_rawDesc), len(file_pkg_api_system_proto_rawDesc)))
	})
	return file_pkg_api_system_proto_rawDescData
}

var file_pkg_api_system_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_api_system_proto_goTypes = []any{
	(*SendFileRequest)(nil),  // 0: system.SendFileRequest
	(*SendFileResponse)(nil), // 1: system.SendFileResponse
	(*FileMetadata)(nil),     // 2: system.FileMetadata
}
var file_pkg_api_system_proto_depIdxs = []int32{
	2, // 0: system.SendFileResponse.metadata:type_name -> system.FileMetadata
	0, // 1: system.SystemService.SendFile:input_type -> system.SendFileRequest
	1, // 2: system.SystemService.SendFile:output_type -> system.SendFileResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_api_system_proto_init() }
func file_pkg_api_system_proto_init() {
	if File_pkg_api_system_proto != nil {
		return
	}
	file_pkg_api_system_proto_msgTypes[1].OneofWrappers = []any{
		(*SendFileResponse_Metadata)(nil),
		(*SendFileResponse_Chunk)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_pkg_api_system_proto_rawDesc), len(file_pkg_api_system_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_api_system_proto_goTypes,
		DependencyIndexes: file_pkg_api_system_proto_depIdxs,
		MessageInfos:      file_pkg_api_system_proto_msgTypes,
	}.Build()
	File_pkg_api_system_proto = out.File
	file_pkg_api_system_proto_goTypes = nil
	file_pkg_api_system_proto_depIdxs = nil
}
