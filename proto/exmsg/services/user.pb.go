// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.5.1
// source: services/user.proto

package services

import (
	models "github.com/nhdms/base-go/proto/exmsg/models"
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

type UserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *UserRequest) Reset() {
	*x = UserRequest{}
	mi := &file_services_user_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRequest) ProtoMessage() {}

func (x *UserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_user_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRequest.ProtoReflect.Descriptor instead.
func (*UserRequest) Descriptor() ([]byte, []int) {
	return file_services_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId                int64         `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Query                 *models.Query `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	Status                int64         `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	DepartmentIds         []int64       `protobuf:"varint,5,rep,packed,name=department_ids,json=departmentIds,proto3" json:"department_ids,omitempty"`
	AncestorDepartmentIds []int64       `protobuf:"varint,6,rep,packed,name=ancestor_department_ids,json=ancestorDepartmentIds,proto3" json:"ancestor_department_ids,omitempty"`
}

func (x *ProfileRequest) Reset() {
	*x = ProfileRequest{}
	mi := &file_services_user_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileRequest) ProtoMessage() {}

func (x *ProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_user_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileRequest.ProtoReflect.Descriptor instead.
func (*ProfileRequest) Descriptor() ([]byte, []int) {
	return file_services_user_proto_rawDescGZIP(), []int{1}
}

func (x *ProfileRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ProfileRequest) GetQuery() *models.Query {
	if x != nil {
		return x.Query
	}
	return nil
}

func (x *ProfileRequest) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *ProfileRequest) GetDepartmentIds() []int64 {
	if x != nil {
		return x.DepartmentIds
	}
	return nil
}

func (x *ProfileRequest) GetAncestorDepartmentIds() []int64 {
	if x != nil {
		return x.AncestorDepartmentIds
	}
	return nil
}

type UserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User     *models.User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	CacheHit bool         `protobuf:"varint,2,opt,name=cache_hit,json=cacheHit,proto3" json:"cache_hit,omitempty"`
}

func (x *UserResponse) Reset() {
	*x = UserResponse{}
	mi := &file_services_user_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserResponse) ProtoMessage() {}

func (x *UserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_user_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserResponse.ProtoReflect.Descriptor instead.
func (*UserResponse) Descriptor() ([]byte, []int) {
	return file_services_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserResponse) GetUser() *models.User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UserResponse) GetCacheHit() bool {
	if x != nil {
		return x.CacheHit
	}
	return false
}

var File_services_user_proto protoreflect.FileDescriptor

var file_services_user_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x11, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a,
	0x0b, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0xcb, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x29, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0d, 0x64, 0x65,
	0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x12, 0x36, 0x0a, 0x17, 0x61,
	0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x5f, 0x64, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x03, 0x52, 0x15, 0x61, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x44, 0x65, 0x70, 0x61, 0x72, 0x74, 0x6d, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x73, 0x22, 0x53, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x63,
	0x61, 0x63, 0x68, 0x65, 0x5f, 0x68, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x63, 0x61, 0x63, 0x68, 0x65, 0x48, 0x69, 0x74, 0x32, 0x57, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x42, 0x79, 0x49, 0x44, 0x12, 0x1b, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6e, 0x68, 0x64, 0x6d, 0x73, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x3b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_services_user_proto_rawDescOnce sync.Once
	file_services_user_proto_rawDescData = file_services_user_proto_rawDesc
)

func file_services_user_proto_rawDescGZIP() []byte {
	file_services_user_proto_rawDescOnce.Do(func() {
		file_services_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_user_proto_rawDescData)
	})
	return file_services_user_proto_rawDescData
}

var file_services_user_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_services_user_proto_goTypes = []any{
	(*UserRequest)(nil),    // 0: exmsg.services.UserRequest
	(*ProfileRequest)(nil), // 1: exmsg.services.ProfileRequest
	(*UserResponse)(nil),   // 2: exmsg.services.UserResponse
	(*models.Query)(nil),   // 3: exmsg.models.Query
	(*models.User)(nil),    // 4: exmsg.models.User
}
var file_services_user_proto_depIdxs = []int32{
	3, // 0: exmsg.services.ProfileRequest.query:type_name -> exmsg.models.Query
	4, // 1: exmsg.services.UserResponse.user:type_name -> exmsg.models.User
	0, // 2: exmsg.services.UserService.GetUserByID:input_type -> exmsg.services.UserRequest
	2, // 3: exmsg.services.UserService.GetUserByID:output_type -> exmsg.services.UserResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_services_user_proto_init() }
func file_services_user_proto_init() {
	if File_services_user_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_user_proto_goTypes,
		DependencyIndexes: file_services_user_proto_depIdxs,
		MessageInfos:      file_services_user_proto_msgTypes,
	}.Build()
	File_services_user_proto = out.File
	file_services_user_proto_rawDesc = nil
	file_services_user_proto_goTypes = nil
	file_services_user_proto_depIdxs = nil
}
