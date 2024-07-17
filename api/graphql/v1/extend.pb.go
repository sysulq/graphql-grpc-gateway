// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: graphql/v1/extend.proto

package graphqlv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Type of the graphql operation
type Type int32

const (
	// DEFAULT invalid type
	Type_DEFAULT Type = 0
	// QUERY query type
	Type_QUERY Type = 1
	// MUTATION mutation type
	Type_MUTATION Type = 2
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "DEFAULT",
		1: "QUERY",
		2: "MUTATION",
	}
	Type_value = map[string]int32{
		"DEFAULT":  0,
		"QUERY":    1,
		"MUTATION": 2,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_graphql_v1_extend_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_graphql_v1_extend_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_graphql_v1_extend_proto_rawDescGZIP(), []int{0}
}

// The options for the oneof
type Oneof struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ignore the oneof
	Ignore bool `protobuf:"varint,1,opt,name=ignore,proto3" json:"ignore,omitempty"`
	// name of the oneof
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Oneof) Reset() {
	*x = Oneof{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graphql_v1_extend_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Oneof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Oneof) ProtoMessage() {}

func (x *Oneof) ProtoReflect() protoreflect.Message {
	mi := &file_graphql_v1_extend_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Oneof.ProtoReflect.Descriptor instead.
func (*Oneof) Descriptor() ([]byte, []int) {
	return file_graphql_v1_extend_proto_rawDescGZIP(), []int{0}
}

func (x *Oneof) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

func (x *Oneof) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The options for the field
type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// required the field is required
	Required bool `protobuf:"varint,1,opt,name=required,proto3" json:"required,omitempty"`
	// name of the field
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// ignore the field
	Ignore bool `protobuf:"varint,3,opt,name=ignore,proto3" json:"ignore,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graphql_v1_extend_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_graphql_v1_extend_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_graphql_v1_extend_proto_rawDescGZIP(), []int{1}
}

func (x *Field) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

func (x *Field) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Field) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

// The options for the rpc
type Rpc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// type of the graphql operation
	Type Type `protobuf:"varint,1,opt,name=type,proto3,enum=graphql.v1.Type" json:"type,omitempty"`
	// name of the graphql operation
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// ignore the rpc
	Ignore bool `protobuf:"varint,3,opt,name=ignore,proto3" json:"ignore,omitempty"`
}

func (x *Rpc) Reset() {
	*x = Rpc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_graphql_v1_extend_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rpc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rpc) ProtoMessage() {}

func (x *Rpc) ProtoReflect() protoreflect.Message {
	mi := &file_graphql_v1_extend_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rpc.ProtoReflect.Descriptor instead.
func (*Rpc) Descriptor() ([]byte, []int) {
	return file_graphql_v1_extend_proto_rawDescGZIP(), []int{2}
}

func (x *Rpc) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_DEFAULT
}

func (x *Rpc) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Rpc) GetIgnore() bool {
	if x != nil {
		return x.Ignore
	}
	return false
}

var file_graphql_v1_extend_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*Rpc)(nil),
		Field:         95279528,
		Name:          "graphql.v1.rpc",
		Tag:           "bytes,95279528,opt,name=rpc",
		Filename:      "graphql/v1/extend.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Field)(nil),
		Field:         95279528,
		Name:          "graphql.v1.field",
		Tag:           "bytes,95279528,opt,name=field",
		Filename:      "graphql/v1/extend.proto",
	},
	{
		ExtendedType:  (*descriptorpb.OneofOptions)(nil),
		ExtensionType: (*Oneof)(nil),
		Field:         95279528,
		Name:          "graphql.v1.oneof",
		Tag:           "bytes,95279528,opt,name=oneof",
		Filename:      "graphql/v1/extend.proto",
	},
}

// Extension fields to descriptorpb.MethodOptions.
var (
	// rpc defines the graphql information for the method
	//
	// optional graphql.v1.Rpc rpc = 95279528;
	E_Rpc = &file_graphql_v1_extend_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// field defines the graphql information for the field
	//
	// optional graphql.v1.Field field = 95279528;
	E_Field = &file_graphql_v1_extend_proto_extTypes[1]
)

// Extension fields to descriptorpb.OneofOptions.
var (
	// oneof defines the graphql information for the oneof
	//
	// optional graphql.v1.Oneof oneof = 95279528;
	E_Oneof = &file_graphql_v1_extend_proto_extTypes[2]
)

var File_graphql_v1_extend_proto protoreflect.FileDescriptor

var file_graphql_v1_extend_proto_rawDesc = []byte{
	0x0a, 0x17, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x71, 0x6c, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x05, 0x4f, 0x6e, 0x65, 0x6f, 0x66,
	0x12, 0x16, 0x0a, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4f, 0x0a, 0x05,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x22, 0x57, 0x0a,
	0x03, 0x52, 0x70, 0x63, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
	0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x2a, 0x2c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b,
	0x0a, 0x07, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x51,
	0x55, 0x45, 0x52, 0x59, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4d, 0x55, 0x54, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x10, 0x02, 0x3a, 0x44, 0x0a, 0x03, 0x72, 0x70, 0x63, 0x12, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa8, 0xb3, 0xb7, 0x2d,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x70, 0x63, 0x52, 0x03, 0x72, 0x70, 0x63, 0x3a, 0x49, 0x0a, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xa8, 0xb3, 0xb7, 0x2d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x71, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x3a, 0x49, 0x0a, 0x05, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa8, 0xb3,
	0xb7, 0x2d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x52, 0x05, 0x6f, 0x6e, 0x65, 0x6f, 0x66,
	0x42, 0xa7, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c,
	0x2e, 0x76, 0x31, 0x42, 0x0b, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x79, 0x73, 0x75, 0x6c, 0x71, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2d, 0x67, 0x72,
	0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71,
	0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x47, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x47, 0x72, 0x61, 0x70,
	0x68, 0x71, 0x6c, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x47, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x47, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x47,
	0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_graphql_v1_extend_proto_rawDescOnce sync.Once
	file_graphql_v1_extend_proto_rawDescData = file_graphql_v1_extend_proto_rawDesc
)

func file_graphql_v1_extend_proto_rawDescGZIP() []byte {
	file_graphql_v1_extend_proto_rawDescOnce.Do(func() {
		file_graphql_v1_extend_proto_rawDescData = protoimpl.X.CompressGZIP(file_graphql_v1_extend_proto_rawDescData)
	})
	return file_graphql_v1_extend_proto_rawDescData
}

var file_graphql_v1_extend_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_graphql_v1_extend_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_graphql_v1_extend_proto_goTypes = []interface{}{
	(Type)(0),                          // 0: graphql.v1.Type
	(*Oneof)(nil),                      // 1: graphql.v1.Oneof
	(*Field)(nil),                      // 2: graphql.v1.Field
	(*Rpc)(nil),                        // 3: graphql.v1.Rpc
	(*descriptorpb.MethodOptions)(nil), // 4: google.protobuf.MethodOptions
	(*descriptorpb.FieldOptions)(nil),  // 5: google.protobuf.FieldOptions
	(*descriptorpb.OneofOptions)(nil),  // 6: google.protobuf.OneofOptions
}
var file_graphql_v1_extend_proto_depIdxs = []int32{
	0, // 0: graphql.v1.Rpc.type:type_name -> graphql.v1.Type
	4, // 1: graphql.v1.rpc:extendee -> google.protobuf.MethodOptions
	5, // 2: graphql.v1.field:extendee -> google.protobuf.FieldOptions
	6, // 3: graphql.v1.oneof:extendee -> google.protobuf.OneofOptions
	3, // 4: graphql.v1.rpc:type_name -> graphql.v1.Rpc
	2, // 5: graphql.v1.field:type_name -> graphql.v1.Field
	1, // 6: graphql.v1.oneof:type_name -> graphql.v1.Oneof
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	4, // [4:7] is the sub-list for extension type_name
	1, // [1:4] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_graphql_v1_extend_proto_init() }
func file_graphql_v1_extend_proto_init() {
	if File_graphql_v1_extend_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_graphql_v1_extend_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Oneof); i {
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
		file_graphql_v1_extend_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
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
		file_graphql_v1_extend_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rpc); i {
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
			RawDescriptor: file_graphql_v1_extend_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 3,
			NumServices:   0,
		},
		GoTypes:           file_graphql_v1_extend_proto_goTypes,
		DependencyIndexes: file_graphql_v1_extend_proto_depIdxs,
		EnumInfos:         file_graphql_v1_extend_proto_enumTypes,
		MessageInfos:      file_graphql_v1_extend_proto_msgTypes,
		ExtensionInfos:    file_graphql_v1_extend_proto_extTypes,
	}.Build()
	File_graphql_v1_extend_proto = out.File
	file_graphql_v1_extend_proto_rawDesc = nil
	file_graphql_v1_extend_proto_goTypes = nil
	file_graphql_v1_extend_proto_depIdxs = nil
}
