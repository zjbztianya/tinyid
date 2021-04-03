// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: id.proto

package errors

import (
	_ "github.com/go-kratos/kratos/v2/api/kratos/api"
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

type Idgen int32

const (
	Idgen_NotGenSegmentID Idgen = 0
	Idgen_ClockBackwards  Idgen = 1
)

// Enum value maps for Idgen.
var (
	Idgen_name = map[int32]string{
		0: "NotGenSegmentID",
		1: "ClockBackwards",
	}
	Idgen_value = map[string]int32{
		"NotGenSegmentID": 0,
		"ClockBackwards":  1,
	}
)

func (x Idgen) Enum() *Idgen {
	p := new(Idgen)
	*p = x
	return p
}

func (x Idgen) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Idgen) Descriptor() protoreflect.EnumDescriptor {
	return file_id_proto_enumTypes[0].Descriptor()
}

func (Idgen) Type() protoreflect.EnumType {
	return &file_id_proto_enumTypes[0]
}

func (x Idgen) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Idgen.Descriptor instead.
func (Idgen) EnumDescriptor() ([]byte, []int) {
	return file_id_proto_rawDescGZIP(), []int{0}
}

var File_id_proto protoreflect.FileDescriptor

var file_id_proto_rawDesc = []byte{
	0x0a, 0x08, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x61, 0x70, 0x69, 0x2e,
	0x69, 0x64, 0x67, 0x65, 0x6e, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x1a, 0x1c, 0x6b, 0x72,
	0x61, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x35, 0x0a, 0x05, 0x49, 0x64,
	0x67, 0x65, 0x6e, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x6f, 0x74, 0x47, 0x65, 0x6e, 0x53, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x6c, 0x6f, 0x63,
	0x6b, 0x42, 0x61, 0x63, 0x6b, 0x77, 0x61, 0x72, 0x64, 0x73, 0x10, 0x01, 0x1a, 0x03, 0xa0, 0x45,
	0x01, 0x42, 0x2a, 0x5a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x64, 0x67, 0x65, 0x6e, 0x2f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x3b, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0xa2, 0x02, 0x0e, 0x41,
	0x50, 0x49, 0x49, 0x64, 0x67, 0x65, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_id_proto_rawDescOnce sync.Once
	file_id_proto_rawDescData = file_id_proto_rawDesc
)

func file_id_proto_rawDescGZIP() []byte {
	file_id_proto_rawDescOnce.Do(func() {
		file_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_id_proto_rawDescData)
	})
	return file_id_proto_rawDescData
}

var file_id_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_id_proto_goTypes = []interface{}{
	(Idgen)(0), // 0: api.idgen.errors.Idgen
}
var file_id_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_id_proto_init() }
func file_id_proto_init() {
	if File_id_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_id_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_id_proto_goTypes,
		DependencyIndexes: file_id_proto_depIdxs,
		EnumInfos:         file_id_proto_enumTypes,
	}.Build()
	File_id_proto = out.File
	file_id_proto_rawDesc = nil
	file_id_proto_goTypes = nil
	file_id_proto_depIdxs = nil
}
