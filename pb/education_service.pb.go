// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: education_service.proto

// Всё как в go

package protobuf

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

var File_education_service_proto protoreflect.FileDescriptor

var file_education_service_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x1a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xd3, 0x05,
	0x0a, 0x09, 0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x60, 0x0a, 0x0a, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x22, 0x0c, 0x2f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x51, 0x0a,
	0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12, 0x09, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x12, 0x7f, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f,
	0x72, 0x74, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x67, 0x65, 0x74, 0x5f,
	0x6d, 0x61, 0x6e, 0x79, 0x5f, 0x73, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x12, 0x58, 0x0a, 0x08, 0x45, 0x64, 0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x32, 0x0a, 0x2f, 0x65,
	0x64, 0x69, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x5d, 0x0a, 0x0a, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x2a, 0x0c, 0x2f, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x12, 0x5c, 0x0a, 0x09, 0x4c, 0x6f,
	0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x79, 0x0a, 0x10, 0x52, 0x65, 0x6e, 0x65,
	0x77, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x21, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x65, 0x77,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x72, 0x65,
	0x6e, 0x65, 0x77, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x3a, 0x01, 0x2a, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x5a, 0x65, 0x6e, 0x53, 0x61, 0x6d, 0x37, 0x2f, 0x45, 0x64, 0x75, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_education_service_proto_goTypes = []interface{}{
	(*CreateUserRequest)(nil),          // 0: protobuf.CreateUserRequest
	(*GetUserRequest)(nil),             // 1: protobuf.GetUserRequest
	(*GetManySortedUsersRequest)(nil),  // 2: protobuf.GetManySortedUsersRequest
	(*EditUserRequest)(nil),            // 3: protobuf.EditUserRequest
	(*DeleteUserRequest)(nil),          // 4: protobuf.DeleteUserRequest
	(*LoginUserRequest)(nil),           // 5: protobuf.LoginUserRequest
	(*RenewAccessTokenRequest)(nil),    // 6: protobuf.RenewAccessTokenRequest
	(*CreateUserResponse)(nil),         // 7: protobuf.CreateUserResponse
	(*GetUserResponse)(nil),            // 8: protobuf.GetUserResponse
	(*GetManySortedUsersResponse)(nil), // 9: protobuf.GetManySortedUsersResponse
	(*EditUserResponse)(nil),           // 10: protobuf.EditUserResponse
	(*DeleteUserResponse)(nil),         // 11: protobuf.DeleteUserResponse
	(*LoginUserResponse)(nil),          // 12: protobuf.LoginUserResponse
	(*RenewAccessTokenResponse)(nil),   // 13: protobuf.RenewAccessTokenResponse
}
var file_education_service_proto_depIdxs = []int32{
	0,  // 0: protobuf.Education.CreateUser:input_type -> protobuf.CreateUserRequest
	1,  // 1: protobuf.Education.GetUser:input_type -> protobuf.GetUserRequest
	2,  // 2: protobuf.Education.GetManySortedUsers:input_type -> protobuf.GetManySortedUsersRequest
	3,  // 3: protobuf.Education.EditUser:input_type -> protobuf.EditUserRequest
	4,  // 4: protobuf.Education.DeleteUser:input_type -> protobuf.DeleteUserRequest
	5,  // 5: protobuf.Education.LoginUser:input_type -> protobuf.LoginUserRequest
	6,  // 6: protobuf.Education.RenewAccessToken:input_type -> protobuf.RenewAccessTokenRequest
	7,  // 7: protobuf.Education.CreateUser:output_type -> protobuf.CreateUserResponse
	8,  // 8: protobuf.Education.GetUser:output_type -> protobuf.GetUserResponse
	9,  // 9: protobuf.Education.GetManySortedUsers:output_type -> protobuf.GetManySortedUsersResponse
	10, // 10: protobuf.Education.EditUser:output_type -> protobuf.EditUserResponse
	11, // 11: protobuf.Education.DeleteUser:output_type -> protobuf.DeleteUserResponse
	12, // 12: protobuf.Education.LoginUser:output_type -> protobuf.LoginUserResponse
	13, // 13: protobuf.Education.RenewAccessToken:output_type -> protobuf.RenewAccessTokenResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_education_service_proto_init() }
func file_education_service_proto_init() {
	if File_education_service_proto != nil {
		return
	}
	file_user_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_education_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_education_service_proto_goTypes,
		DependencyIndexes: file_education_service_proto_depIdxs,
	}.Build()
	File_education_service_proto = out.File
	file_education_service_proto_rawDesc = nil
	file_education_service_proto_goTypes = nil
	file_education_service_proto_depIdxs = nil
}
