// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: education_service.proto

// Всё как в go

package protobuf

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
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
	0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d,
	0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xd6, 0x15, 0x0a, 0x09, 0x45, 0x64, 0x75, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x60, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x22, 0x0c, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x51, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12,
	0x09, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x12, 0x7f, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x18, 0x12, 0x16, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x79, 0x5f, 0x73,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x58, 0x0a, 0x08, 0x45,
	0x64, 0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64,
	0x69, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x32, 0x0a, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x3a, 0x01, 0x2a, 0x12, 0x5d, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x2a, 0x0c, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x12, 0x5c, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x10, 0x22, 0x0b, 0x2f, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x3a,
	0x01, 0x2a, 0x12, 0x79, 0x0a, 0x10, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x65, 0x77, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x72, 0x65, 0x6e, 0x65, 0x77, 0x5f, 0x61, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x3a, 0x01, 0x2a, 0x12, 0x61, 0x0a,
	0x0b, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0f, 0x12, 0x0d, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x6c, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x22, 0x0f, 0x2f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x01, 0x2a, 0x12, 0x5d,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12,
	0x0c, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x83, 0x01,
	0x0a, 0x14, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x4f, 0x66, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x4f, 0x66, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x4f, 0x66, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x6f, 0x66, 0x5f, 0x61, 0x72, 0x74, 0x69,
	0x63, 0x6c, 0x65, 0x12, 0x64, 0x0a, 0x0b, 0x45, 0x64, 0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64,
	0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x32, 0x0d, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x01, 0x2a, 0x12, 0x69, 0x0a, 0x0d, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x11, 0x2a, 0x0f, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x6c, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x72,
	0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x22, 0x0f,
	0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x3a,
	0x01, 0x2a, 0x12, 0x5d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65,
	0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69,
	0x63, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x12, 0x64, 0x0a, 0x0b, 0x45, 0x64, 0x69, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65,
	0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74,
	0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x41, 0x72,
	0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x12, 0x32, 0x0d, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x61, 0x72, 0x74,
	0x69, 0x63, 0x6c, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x69, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x11, 0x2a, 0x0f, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x6c, 0x65, 0x12, 0x97, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12,
	0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x72,
	0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65,
	0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c,
	0x2f, 0x67, 0x65, 0x74, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x5f, 0x77, 0x69,
	0x74, 0x68, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x8b, 0x01, 0x0a,
	0x15, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41, 0x72,
	0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e,
	0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x12,
	0x19, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x79, 0x5f, 0x73, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0xc1, 0x01, 0x0a, 0x22, 0x47,
	0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41, 0x72, 0x74, 0x69,
	0x63, 0x6c, 0x65, 0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x12, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74,
	0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x57, 0x69, 0x74, 0x68, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x30, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x2a, 0x12, 0x28, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x6d, 0x61, 0x6e, 0x79, 0x5f,
	0x73, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x5f,
	0x77, 0x69, 0x74, 0x68, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x60,
	0x0a, 0x08, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x47, 0x65, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x67, 0x65, 0x74, 0x5f,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x7d,
	0x12, 0x5c, 0x0a, 0x09, 0x45, 0x64, 0x69, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x32, 0x0b,
	0x2f, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x6c,
	0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1a, 0x2a, 0x18, 0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x5f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x2f, 0x7b, 0x69, 0x64, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x7d, 0x12, 0x5c, 0x0a, 0x09,
	0x4c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x22, 0x0b, 0x2f, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x64, 0x0a, 0x0b, 0x52, 0x65,
	0x6e, 0x61, 0x6d, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x52, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x32, 0x0d,
	0x2f, 0x72, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x3a, 0x01, 0x2a,
	0x42, 0x63, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x5a,
	0x65, 0x6e, 0x53, 0x61, 0x6d, 0x37, 0x2f, 0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x92, 0x41, 0x39, 0x12, 0x37, 0x0a, 0x09,
	0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x25, 0x0a, 0x07, 0x5a, 0x65, 0x6e,
	0x53, 0x61, 0x6d, 0x37, 0x12, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x5a, 0x65, 0x6e, 0x53, 0x61, 0x6d, 0x37,
	0x32, 0x03, 0x30, 0x2e, 0x38, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_education_service_proto_goTypes = []interface{}{
	(*CreateUserRequest)(nil),                          // 0: protobuf.CreateUserRequest
	(*GetUserRequest)(nil),                             // 1: protobuf.GetUserRequest
	(*GetManySortedUsersRequest)(nil),                  // 2: protobuf.GetManySortedUsersRequest
	(*EditUserRequest)(nil),                            // 3: protobuf.EditUserRequest
	(*DeleteUserRequest)(nil),                          // 4: protobuf.DeleteUserRequest
	(*LoginUserRequest)(nil),                           // 5: protobuf.LoginUserRequest
	(*RenewAccessTokenRequest)(nil),                    // 6: protobuf.RenewAccessTokenRequest
	(*VerifyEmailRequest)(nil),                         // 7: protobuf.VerifyEmailRequest
	(*CreateCommentRequest)(nil),                       // 8: protobuf.CreateCommentRequest
	(*GetCommentRequest)(nil),                          // 9: protobuf.GetCommentRequest
	(*GetCommentsOfArticleRequest)(nil),                // 10: protobuf.GetCommentsOfArticleRequest
	(*EditCommentRequest)(nil),                         // 11: protobuf.EditCommentRequest
	(*DeleteCommentRequest)(nil),                       // 12: protobuf.DeleteCommentRequest
	(*CreateArticleRequest)(nil),                       // 13: protobuf.CreateArticleRequest
	(*GetArticleRequest)(nil),                          // 14: protobuf.GetArticleRequest
	(*EditArticleRequest)(nil),                         // 15: protobuf.EditArticleRequest
	(*DeleteArticleRequest)(nil),                       // 16: protobuf.DeleteArticleRequest
	(*GetArticlesWithAttributeRequest)(nil),            // 17: protobuf.GetArticlesWithAttributeRequest
	(*GetManySortedArticlesRequest)(nil),               // 18: protobuf.GetManySortedArticlesRequest
	(*GetManySortedArticlesWithAttributeRequest)(nil),  // 19: protobuf.GetManySortedArticlesWithAttributeRequest
	(*GetImageRequest)(nil),                            // 20: protobuf.GetImageRequest
	(*EditImageRequest)(nil),                           // 21: protobuf.EditImageRequest
	(*DeleteImageRequest)(nil),                         // 22: protobuf.DeleteImageRequest
	(*LoadImageRequest)(nil),                           // 23: protobuf.LoadImageRequest
	(*RenameImageRequest)(nil),                         // 24: protobuf.RenameImageRequest
	(*CreateUserResponse)(nil),                         // 25: protobuf.CreateUserResponse
	(*GetUserResponse)(nil),                            // 26: protobuf.GetUserResponse
	(*GetManySortedUsersResponse)(nil),                 // 27: protobuf.GetManySortedUsersResponse
	(*EditUserResponse)(nil),                           // 28: protobuf.EditUserResponse
	(*DeleteUserResponse)(nil),                         // 29: protobuf.DeleteUserResponse
	(*LoginUserResponse)(nil),                          // 30: protobuf.LoginUserResponse
	(*RenewAccessTokenResponse)(nil),                   // 31: protobuf.RenewAccessTokenResponse
	(*VerifyEmailResponse)(nil),                        // 32: protobuf.VerifyEmailResponse
	(*CreateCommentResponse)(nil),                      // 33: protobuf.CreateCommentResponse
	(*GetCommentResponse)(nil),                         // 34: protobuf.GetCommentResponse
	(*GetCommentsOfArticleResponse)(nil),               // 35: protobuf.GetCommentsOfArticleResponse
	(*EditCommentResponse)(nil),                        // 36: protobuf.EditCommentResponse
	(*DeleteCommentResponse)(nil),                      // 37: protobuf.DeleteCommentResponse
	(*CreateArticleResponse)(nil),                      // 38: protobuf.CreateArticleResponse
	(*GetArticleResponse)(nil),                         // 39: protobuf.GetArticleResponse
	(*EditArticleResponse)(nil),                        // 40: protobuf.EditArticleResponse
	(*DeleteArticleResponse)(nil),                      // 41: protobuf.DeleteArticleResponse
	(*GetArticlesWithAttributeResponse)(nil),           // 42: protobuf.GetArticlesWithAttributeResponse
	(*GetManySortedArticlesResponse)(nil),              // 43: protobuf.GetManySortedArticlesResponse
	(*GetManySortedArticlesWithAttributeResponse)(nil), // 44: protobuf.GetManySortedArticlesWithAttributeResponse
	(*GetImageResponse)(nil),                           // 45: protobuf.GetImageResponse
	(*EditImageResponse)(nil),                          // 46: protobuf.EditImageResponse
	(*DeleteImageResponse)(nil),                        // 47: protobuf.DeleteImageResponse
	(*LoadImageResponse)(nil),                          // 48: protobuf.LoadImageResponse
	(*RenameImageResponse)(nil),                        // 49: protobuf.RenameImageResponse
}
var file_education_service_proto_depIdxs = []int32{
	0,  // 0: protobuf.Education.CreateUser:input_type -> protobuf.CreateUserRequest
	1,  // 1: protobuf.Education.GetUser:input_type -> protobuf.GetUserRequest
	2,  // 2: protobuf.Education.GetManySortedUsers:input_type -> protobuf.GetManySortedUsersRequest
	3,  // 3: protobuf.Education.EditUser:input_type -> protobuf.EditUserRequest
	4,  // 4: protobuf.Education.DeleteUser:input_type -> protobuf.DeleteUserRequest
	5,  // 5: protobuf.Education.LoginUser:input_type -> protobuf.LoginUserRequest
	6,  // 6: protobuf.Education.RenewAccessToken:input_type -> protobuf.RenewAccessTokenRequest
	7,  // 7: protobuf.Education.VerifyEmail:input_type -> protobuf.VerifyEmailRequest
	8,  // 8: protobuf.Education.CreateComment:input_type -> protobuf.CreateCommentRequest
	9,  // 9: protobuf.Education.GetComment:input_type -> protobuf.GetCommentRequest
	10, // 10: protobuf.Education.GetCommentsOfArticle:input_type -> protobuf.GetCommentsOfArticleRequest
	11, // 11: protobuf.Education.EditComment:input_type -> protobuf.EditCommentRequest
	12, // 12: protobuf.Education.DeleteComment:input_type -> protobuf.DeleteCommentRequest
	13, // 13: protobuf.Education.CreateArticle:input_type -> protobuf.CreateArticleRequest
	14, // 14: protobuf.Education.GetArticle:input_type -> protobuf.GetArticleRequest
	15, // 15: protobuf.Education.EditArticle:input_type -> protobuf.EditArticleRequest
	16, // 16: protobuf.Education.DeleteArticle:input_type -> protobuf.DeleteArticleRequest
	17, // 17: protobuf.Education.GetArticlesWithAttribute:input_type -> protobuf.GetArticlesWithAttributeRequest
	18, // 18: protobuf.Education.GetManySortedArticles:input_type -> protobuf.GetManySortedArticlesRequest
	19, // 19: protobuf.Education.GetManySortedArticlesWithAttribute:input_type -> protobuf.GetManySortedArticlesWithAttributeRequest
	20, // 20: protobuf.Education.GetImage:input_type -> protobuf.GetImageRequest
	21, // 21: protobuf.Education.EditImage:input_type -> protobuf.EditImageRequest
	22, // 22: protobuf.Education.DeleteImage:input_type -> protobuf.DeleteImageRequest
	23, // 23: protobuf.Education.LoadImage:input_type -> protobuf.LoadImageRequest
	24, // 24: protobuf.Education.RenameImage:input_type -> protobuf.RenameImageRequest
	25, // 25: protobuf.Education.CreateUser:output_type -> protobuf.CreateUserResponse
	26, // 26: protobuf.Education.GetUser:output_type -> protobuf.GetUserResponse
	27, // 27: protobuf.Education.GetManySortedUsers:output_type -> protobuf.GetManySortedUsersResponse
	28, // 28: protobuf.Education.EditUser:output_type -> protobuf.EditUserResponse
	29, // 29: protobuf.Education.DeleteUser:output_type -> protobuf.DeleteUserResponse
	30, // 30: protobuf.Education.LoginUser:output_type -> protobuf.LoginUserResponse
	31, // 31: protobuf.Education.RenewAccessToken:output_type -> protobuf.RenewAccessTokenResponse
	32, // 32: protobuf.Education.VerifyEmail:output_type -> protobuf.VerifyEmailResponse
	33, // 33: protobuf.Education.CreateComment:output_type -> protobuf.CreateCommentResponse
	34, // 34: protobuf.Education.GetComment:output_type -> protobuf.GetCommentResponse
	35, // 35: protobuf.Education.GetCommentsOfArticle:output_type -> protobuf.GetCommentsOfArticleResponse
	36, // 36: protobuf.Education.EditComment:output_type -> protobuf.EditCommentResponse
	37, // 37: protobuf.Education.DeleteComment:output_type -> protobuf.DeleteCommentResponse
	38, // 38: protobuf.Education.CreateArticle:output_type -> protobuf.CreateArticleResponse
	39, // 39: protobuf.Education.GetArticle:output_type -> protobuf.GetArticleResponse
	40, // 40: protobuf.Education.EditArticle:output_type -> protobuf.EditArticleResponse
	41, // 41: protobuf.Education.DeleteArticle:output_type -> protobuf.DeleteArticleResponse
	42, // 42: protobuf.Education.GetArticlesWithAttribute:output_type -> protobuf.GetArticlesWithAttributeResponse
	43, // 43: protobuf.Education.GetManySortedArticles:output_type -> protobuf.GetManySortedArticlesResponse
	44, // 44: protobuf.Education.GetManySortedArticlesWithAttribute:output_type -> protobuf.GetManySortedArticlesWithAttributeResponse
	45, // 45: protobuf.Education.GetImage:output_type -> protobuf.GetImageResponse
	46, // 46: protobuf.Education.EditImage:output_type -> protobuf.EditImageResponse
	47, // 47: protobuf.Education.DeleteImage:output_type -> protobuf.DeleteImageResponse
	48, // 48: protobuf.Education.LoadImage:output_type -> protobuf.LoadImageResponse
	49, // 49: protobuf.Education.RenameImage:output_type -> protobuf.RenameImageResponse
	25, // [25:50] is the sub-list for method output_type
	0,  // [0:25] is the sub-list for method input_type
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
	file_comment_proto_init()
	file_article_proto_init()
	file_image_proto_init()
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
