// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: internal/state/state.proto

package state

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	player "mafia-grpc/internal/player"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PartOfTheDay int32

const (
	PartOfTheDay_UNKNOWN PartOfTheDay = 0
	PartOfTheDay_DAY     PartOfTheDay = 1
	PartOfTheDay_NIGHT   PartOfTheDay = 2
)

// Enum value maps for PartOfTheDay.
var (
	PartOfTheDay_name = map[int32]string{
		0: "UNKNOWN",
		1: "DAY",
		2: "NIGHT",
	}
	PartOfTheDay_value = map[string]int32{
		"UNKNOWN": 0,
		"DAY":     1,
		"NIGHT":   2,
	}
)

func (x PartOfTheDay) Enum() *PartOfTheDay {
	p := new(PartOfTheDay)
	*p = x
	return p
}

func (x PartOfTheDay) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PartOfTheDay) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_state_state_proto_enumTypes[0].Descriptor()
}

func (PartOfTheDay) Type() protoreflect.EnumType {
	return &file_internal_state_state_proto_enumTypes[0]
}

func (x PartOfTheDay) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PartOfTheDay.Descriptor instead.
func (PartOfTheDay) EnumDescriptor() ([]byte, []int) {
	return file_internal_state_state_proto_rawDescGZIP(), []int{0}
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Players      []*player.Player  `protobuf:"bytes,1,rep,name=players,proto3" json:"players,omitempty"`
	PartOfTheDay PartOfTheDay      `protobuf:"varint,2,opt,name=partOfTheDay,proto3,enum=state.PartOfTheDay" json:"partOfTheDay,omitempty"`
	Voices       map[string]string `protobuf:"bytes,3,rep,name=voices,proto3" json:"voices,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_state_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_internal_state_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_internal_state_state_proto_rawDescGZIP(), []int{0}
}

func (x *State) GetPlayers() []*player.Player {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *State) GetPartOfTheDay() PartOfTheDay {
	if x != nil {
		return x.PartOfTheDay
	}
	return PartOfTheDay_UNKNOWN
}

func (x *State) GetVoices() map[string]string {
	if x != nil {
		return x.Voices
	}
	return nil
}

var File_internal_state_state_proto protoreflect.FileDescriptor

var file_internal_state_state_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x1a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xd7, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x73, 0x12, 0x37, 0x0a, 0x0c, 0x70, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x54,
	0x68, 0x65, 0x44, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x54, 0x68, 0x65, 0x44, 0x61, 0x79,
	0x52, 0x0c, 0x70, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x54, 0x68, 0x65, 0x44, 0x61, 0x79, 0x12, 0x30,
	0x0a, 0x06, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x56, 0x6f, 0x69,
	0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73,
	0x1a, 0x39, 0x0a, 0x0b, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x2f, 0x0a, 0x0c, 0x50,
	0x61, 0x72, 0x74, 0x4f, 0x66, 0x54, 0x68, 0x65, 0x44, 0x61, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x41, 0x59, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x4e, 0x49, 0x47, 0x48, 0x54, 0x10, 0x02, 0x42, 0x1b, 0x5a, 0x19,
	0x6d, 0x61, 0x66, 0x69, 0x61, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_internal_state_state_proto_rawDescOnce sync.Once
	file_internal_state_state_proto_rawDescData = file_internal_state_state_proto_rawDesc
)

func file_internal_state_state_proto_rawDescGZIP() []byte {
	file_internal_state_state_proto_rawDescOnce.Do(func() {
		file_internal_state_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_state_state_proto_rawDescData)
	})
	return file_internal_state_state_proto_rawDescData
}

var file_internal_state_state_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_state_state_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_state_state_proto_goTypes = []interface{}{
	(PartOfTheDay)(0),     // 0: state.PartOfTheDay
	(*State)(nil),         // 1: state.State
	nil,                   // 2: state.State.VoicesEntry
	(*player.Player)(nil), // 3: player.Player
}
var file_internal_state_state_proto_depIdxs = []int32{
	3, // 0: state.State.players:type_name -> player.Player
	0, // 1: state.State.partOfTheDay:type_name -> state.PartOfTheDay
	2, // 2: state.State.voices:type_name -> state.State.VoicesEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_state_state_proto_init() }
func file_internal_state_state_proto_init() {
	if File_internal_state_state_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_state_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
			RawDescriptor: file_internal_state_state_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_state_state_proto_goTypes,
		DependencyIndexes: file_internal_state_state_proto_depIdxs,
		EnumInfos:         file_internal_state_state_proto_enumTypes,
		MessageInfos:      file_internal_state_state_proto_msgTypes,
	}.Build()
	File_internal_state_state_proto = out.File
	file_internal_state_state_proto_rawDesc = nil
	file_internal_state_state_proto_goTypes = nil
	file_internal_state_state_proto_depIdxs = nil
}
