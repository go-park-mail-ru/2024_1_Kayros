// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: session/session.proto

package session

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SetSessionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value    string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Database int32  `protobuf:"varint,3,opt,name=database,proto3" json:"database,omitempty"`
}

func (x *SetSessionData) Reset() {
	*x = SetSessionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_session_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetSessionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetSessionData) ProtoMessage() {}

func (x *SetSessionData) ProtoReflect() protoreflect.Message {
	mi := &file_session_session_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetSessionData.ProtoReflect.Descriptor instead.
func (*SetSessionData) Descriptor() ([]byte, []int) {
	return file_session_session_proto_rawDescGZIP(), []int{0}
}

func (x *SetSessionData) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetSessionData) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SetSessionData) GetDatabase() int32 {
	if x != nil {
		return x.Database
	}
	return 0
}

type GetSessionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Database int32  `protobuf:"varint,2,opt,name=database,proto3" json:"database,omitempty"`
}

func (x *GetSessionData) Reset() {
	*x = GetSessionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_session_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSessionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSessionData) ProtoMessage() {}

func (x *GetSessionData) ProtoReflect() protoreflect.Message {
	mi := &file_session_session_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSessionData.ProtoReflect.Descriptor instead.
func (*GetSessionData) Descriptor() ([]byte, []int) {
	return file_session_session_proto_rawDescGZIP(), []int{1}
}

func (x *GetSessionData) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GetSessionData) GetDatabase() int32 {
	if x != nil {
		return x.Database
	}
	return 0
}

type DeleteSessionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Database int32  `protobuf:"varint,2,opt,name=database,proto3" json:"database,omitempty"`
}

func (x *DeleteSessionData) Reset() {
	*x = DeleteSessionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_session_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSessionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSessionData) ProtoMessage() {}

func (x *DeleteSessionData) ProtoReflect() protoreflect.Message {
	mi := &file_session_session_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSessionData.ProtoReflect.Descriptor instead.
func (*DeleteSessionData) Descriptor() ([]byte, []int) {
	return file_session_session_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteSessionData) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DeleteSessionData) GetDatabase() int32 {
	if x != nil {
		return x.Database
	}
	return 0
}

type SessionValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *SessionValue) Reset() {
	*x = SessionValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_session_session_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionValue) ProtoMessage() {}

func (x *SessionValue) ProtoReflect() protoreflect.Message {
	mi := &file_session_session_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionValue.ProtoReflect.Descriptor instead.
func (*SessionValue) Descriptor() ([]byte, []int) {
	return file_session_session_proto_rawDescGZIP(), []int{3}
}

func (x *SessionValue) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_session_session_proto protoreflect.FileDescriptor

var file_session_session_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a,
	0x0e, 0x53, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x22, 0x3e, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x22, 0x41, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x22, 0x22, 0x0a, 0x0c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0xd2, 0x01, 0x0a, 0x0e, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x3d, 0x0a,
	0x0a, 0x53, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x2e, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3c, 0x0a, 0x0a,
	0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x2e, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44,
	0x61, 0x74, 0x61, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x43, 0x0a, 0x0d, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_session_session_proto_rawDescOnce sync.Once
	file_session_session_proto_rawDescData = file_session_session_proto_rawDesc
)

func file_session_session_proto_rawDescGZIP() []byte {
	file_session_session_proto_rawDescOnce.Do(func() {
		file_session_session_proto_rawDescData = protoimpl.X.CompressGZIP(file_session_session_proto_rawDescData)
	})
	return file_session_session_proto_rawDescData
}

var file_session_session_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_session_session_proto_goTypes = []interface{}{
	(*SetSessionData)(nil),    // 0: session.SetSessionData
	(*GetSessionData)(nil),    // 1: session.GetSessionData
	(*DeleteSessionData)(nil), // 2: session.DeleteSessionData
	(*SessionValue)(nil),      // 3: session.SessionValue
	(*emptypb.Empty)(nil),     // 4: google.protobuf.Empty
}
var file_session_session_proto_depIdxs = []int32{
	0, // 0: session.SessionManager.SetSession:input_type -> session.SetSessionData
	1, // 1: session.SessionManager.GetSession:input_type -> session.GetSessionData
	2, // 2: session.SessionManager.DeleteSession:input_type -> session.DeleteSessionData
	4, // 3: session.SessionManager.SetSession:output_type -> google.protobuf.Empty
	3, // 4: session.SessionManager.GetSession:output_type -> session.SessionValue
	4, // 5: session.SessionManager.DeleteSession:output_type -> google.protobuf.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_session_session_proto_init() }
func file_session_session_proto_init() {
	if File_session_session_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_session_session_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetSessionData); i {
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
		file_session_session_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSessionData); i {
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
		file_session_session_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteSessionData); i {
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
		file_session_session_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionValue); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_session_session_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_session_session_proto_goTypes,
		DependencyIndexes: file_session_session_proto_depIdxs,
		MessageInfos:      file_session_session_proto_msgTypes,
	}.Build()
	File_session_session_proto = out.File
	file_session_session_proto_rawDesc = nil
	file_session_session_proto_goTypes = nil
	file_session_session_proto_depIdxs = nil
}