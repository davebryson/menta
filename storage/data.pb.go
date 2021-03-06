// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data.proto

package storage

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

// State commit information
type CommitData struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Version              int64    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommitData) Reset()         { *m = CommitData{} }
func (m *CommitData) String() string { return proto.CompactTextString(m) }
func (*CommitData) ProtoMessage()    {}
func (*CommitData) Descriptor() ([]byte, []int) {
	return fileDescriptor_871986018790d2fd, []int{0}
}

func (m *CommitData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitData.Unmarshal(m, b)
}
func (m *CommitData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitData.Marshal(b, m, deterministic)
}
func (m *CommitData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitData.Merge(m, src)
}
func (m *CommitData) XXX_Size() int {
	return xxx_messageInfo_CommitData.Size(m)
}
func (m *CommitData) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitData.DiscardUnknown(m)
}

var xxx_messageInfo_CommitData proto.InternalMessageInfo

func (m *CommitData) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *CommitData) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func init() {
	proto.RegisterType((*CommitData)(nil), "storage.CommitData")
}

func init() { proto.RegisterFile("data.proto", fileDescriptor_871986018790d2fd) }

var fileDescriptor_871986018790d2fd = []byte{
	// 101 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x49, 0x2c, 0x49,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2f, 0x2e, 0xc9, 0x2f, 0x4a, 0x4c, 0x4f, 0x55,
	0xb2, 0xe2, 0xe2, 0x72, 0xce, 0xcf, 0xcd, 0xcd, 0x2c, 0x71, 0x49, 0x2c, 0x49, 0x14, 0x12, 0xe2,
	0x62, 0xc9, 0x48, 0x2c, 0xce, 0x90, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0xb3, 0x85, 0x24,
	0xb8, 0xd8, 0xcb, 0x52, 0x8b, 0x8a, 0x33, 0xf3, 0xf3, 0x24, 0x98, 0x14, 0x18, 0x35, 0x98, 0x83,
	0x60, 0xdc, 0x24, 0x36, 0xb0, 0x59, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x59, 0x19, 0x96,
	0x24, 0x59, 0x00, 0x00, 0x00,
}
