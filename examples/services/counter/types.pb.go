// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types.proto

package counter

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Message
type Increment struct {
	Value                uint32   `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Increment) Reset()         { *m = Increment{} }
func (m *Increment) String() string { return proto.CompactTextString(m) }
func (*Increment) ProtoMessage()    {}
func (*Increment) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

func (m *Increment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Increment.Unmarshal(m, b)
}
func (m *Increment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Increment.Marshal(b, m, deterministic)
}
func (m *Increment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Increment.Merge(m, src)
}
func (m *Increment) XXX_Size() int {
	return xxx_messageInfo_Increment.Size(m)
}
func (m *Increment) XXX_DiscardUnknown() {
	xxx_messageInfo_Increment.DiscardUnknown(m)
}

var xxx_messageInfo_Increment proto.InternalMessageInfo

func (m *Increment) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Storage
type CountValue struct {
	Current              uint32   `protobuf:"varint,1,opt,name=current,proto3" json:"current,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CountValue) Reset()         { *m = CountValue{} }
func (m *CountValue) String() string { return proto.CompactTextString(m) }
func (*CountValue) ProtoMessage()    {}
func (*CountValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

func (m *CountValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CountValue.Unmarshal(m, b)
}
func (m *CountValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CountValue.Marshal(b, m, deterministic)
}
func (m *CountValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CountValue.Merge(m, src)
}
func (m *CountValue) XXX_Size() int {
	return xxx_messageInfo_CountValue.Size(m)
}
func (m *CountValue) XXX_DiscardUnknown() {
	xxx_messageInfo_CountValue.DiscardUnknown(m)
}

var xxx_messageInfo_CountValue proto.InternalMessageInfo

func (m *CountValue) GetCurrent() uint32 {
	if m != nil {
		return m.Current
	}
	return 0
}

func init() {
	proto.RegisterType((*Increment)(nil), "counter.Increment")
	proto.RegisterType((*CountValue)(nil), "counter.CountValue")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 106 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0xa9, 0x2c, 0x48,
	0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0xce, 0x2f, 0xcd, 0x2b, 0x49, 0x2d,
	0x52, 0x52, 0xe4, 0xe2, 0xf4, 0xcc, 0x4b, 0x2e, 0x4a, 0xcd, 0x4d, 0xcd, 0x2b, 0x11, 0x12, 0xe1,
	0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x82, 0x70, 0x94,
	0xd4, 0xb8, 0xb8, 0x9c, 0x41, 0xaa, 0xc3, 0x40, 0x3c, 0x21, 0x09, 0x2e, 0xf6, 0xe4, 0xd2, 0xa2,
	0xa2, 0xd4, 0xbc, 0x12, 0xa8, 0x2a, 0x18, 0x37, 0x89, 0x0d, 0x6c, 0xb4, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x2c, 0xc9, 0xcc, 0xa0, 0x69, 0x00, 0x00, 0x00,
}