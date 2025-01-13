// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.5.1
// source: services/webhook.proto

package services

import (
	models "github.com/nhdms/base-go/proto/exmsg/models"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_services_webhook_proto protoreflect.FileDescriptor

var file_services_webhook_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x14, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0x54, 0x0a, 0x0e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x42, 0x0a, 0x0a, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x4c,
	0x6f, 0x67, 0x73, 0x12, 0x1b, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x1a, 0x17, 0x2e, 0x65, 0x78, 0x6d, 0x73, 0x67, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e,
	0x53, 0x51, 0x4c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74,
	0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x37, 0x39, 0x32, 0x33, 0x2f, 0x61, 0x74,
	0x68, 0x65, 0x6e, 0x61, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78,
	0x6d, 0x73, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x3b, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_services_webhook_proto_goTypes = []any{
	(*models.WebhookEvents)(nil), // 0: exmsg.models.WebhookEvents
	(*models.SQLResult)(nil),     // 1: exmsg.models.SQLResult
}
var file_services_webhook_proto_depIdxs = []int32{
	0, // 0: exmsg.services.WebhookService.InsertLogs:input_type -> exmsg.models.WebhookEvents
	1, // 1: exmsg.services.WebhookService.InsertLogs:output_type -> exmsg.models.SQLResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_webhook_proto_init() }
func file_services_webhook_proto_init() {
	if File_services_webhook_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_services_webhook_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_webhook_proto_goTypes,
		DependencyIndexes: file_services_webhook_proto_depIdxs,
	}.Build()
	File_services_webhook_proto = out.File
	file_services_webhook_proto_rawDesc = nil
	file_services_webhook_proto_goTypes = nil
	file_services_webhook_proto_depIdxs = nil
}
