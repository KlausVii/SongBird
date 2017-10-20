// Code generated by protoc-gen-go.
// source: github.com/KlausVii/SongBird/proto/api.proto
// DO NOT EDIT!

/*
Package songbird is a generated protocol buffer package.

It is generated from these files:
	github.com/KlausVii/SongBird/proto/api.proto

It has these top-level messages:
	Request
	Response
*/
package songbird

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	Topic string `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	Req   string `protobuf:"bytes,2,opt,name=req" json:"req,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *Request) GetReq() string {
	if m != nil {
		return m.Req
	}
	return ""
}

type Response struct {
	Rsp string `protobuf:"bytes,1,opt,name=rsp" json:"rsp,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetRsp() string {
	if m != nil {
		return m.Rsp
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "songbird.Request")
	proto.RegisterType((*Response)(nil), "songbird.Response")
}

func init() { proto.RegisterFile("github.com/KlausVii/SongBird/proto/api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x49, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xf7, 0xce, 0x49, 0x2c, 0x2d, 0x0e, 0xcb, 0xcc, 0xd4,
	0x0f, 0xce, 0xcf, 0x4b, 0x77, 0xca, 0x2c, 0x4a, 0xd1, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f,
	0x2c, 0xc8, 0xd4, 0x03, 0xb3, 0x84, 0x38, 0x8a, 0xf3, 0xf3, 0xd2, 0x93, 0x32, 0x8b, 0x52, 0x94,
	0x0c, 0xb9, 0xd8, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x44, 0xb8, 0x58, 0x4b, 0xf2,
	0x0b, 0x32, 0x93, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x21, 0x01, 0x2e, 0xe6,
	0xa2, 0xd4, 0x42, 0x09, 0x26, 0xb0, 0x18, 0x88, 0xa9, 0x24, 0xc3, 0xc5, 0x11, 0x94, 0x5a, 0x5c,
	0x90, 0x9f, 0x57, 0x9c, 0x0a, 0x96, 0x2d, 0x2e, 0x80, 0xea, 0x00, 0x31, 0x93, 0xd8, 0xc0, 0x36,
	0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x05, 0x72, 0xec, 0x91, 0x00, 0x00, 0x00,
}
