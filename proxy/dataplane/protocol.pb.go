// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/protocol.proto

/*
Package envoy_api_v2 is a generated protocol buffer package.

It is generated from these files:
	api/protocol.proto

It has these top-level messages:
	TcpProtocolOptions
	Http1ProtocolOptions
	Http2ProtocolOptions
	GrpcProtocolOptions
*/
package envoy_api_v2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/wrappers"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TcpProtocolOptions struct {
}

func (m *TcpProtocolOptions) Reset()                    { *m = TcpProtocolOptions{} }
func (m *TcpProtocolOptions) String() string            { return proto.CompactTextString(m) }
func (*TcpProtocolOptions) ProtoMessage()               {}
func (*TcpProtocolOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Http1ProtocolOptions struct {
	AllowAbsoluteUrl *google_protobuf.BoolValue `protobuf:"bytes,1,opt,name=allow_absolute_url,json=allowAbsoluteUrl" json:"allow_absolute_url,omitempty"`
}

func (m *Http1ProtocolOptions) Reset()                    { *m = Http1ProtocolOptions{} }
func (m *Http1ProtocolOptions) String() string            { return proto.CompactTextString(m) }
func (*Http1ProtocolOptions) ProtoMessage()               {}
func (*Http1ProtocolOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Http1ProtocolOptions) GetAllowAbsoluteUrl() *google_protobuf.BoolValue {
	if m != nil {
		return m.AllowAbsoluteUrl
	}
	return nil
}

type Http2ProtocolOptions struct {
	HpackTableSize              *google_protobuf.UInt32Value `protobuf:"bytes,1,opt,name=hpack_table_size,json=hpackTableSize" json:"hpack_table_size,omitempty"`
	MaxConcurrentStreams        *google_protobuf.UInt32Value `protobuf:"bytes,2,opt,name=max_concurrent_streams,json=maxConcurrentStreams" json:"max_concurrent_streams,omitempty"`
	InitialStreamWindowSize     *google_protobuf.UInt32Value `protobuf:"bytes,3,opt,name=initial_stream_window_size,json=initialStreamWindowSize" json:"initial_stream_window_size,omitempty"`
	InitialConnectionWindowSize *google_protobuf.UInt32Value `protobuf:"bytes,4,opt,name=initial_connection_window_size,json=initialConnectionWindowSize" json:"initial_connection_window_size,omitempty"`
}

func (m *Http2ProtocolOptions) Reset()                    { *m = Http2ProtocolOptions{} }
func (m *Http2ProtocolOptions) String() string            { return proto.CompactTextString(m) }
func (*Http2ProtocolOptions) ProtoMessage()               {}
func (*Http2ProtocolOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Http2ProtocolOptions) GetHpackTableSize() *google_protobuf.UInt32Value {
	if m != nil {
		return m.HpackTableSize
	}
	return nil
}

func (m *Http2ProtocolOptions) GetMaxConcurrentStreams() *google_protobuf.UInt32Value {
	if m != nil {
		return m.MaxConcurrentStreams
	}
	return nil
}

func (m *Http2ProtocolOptions) GetInitialStreamWindowSize() *google_protobuf.UInt32Value {
	if m != nil {
		return m.InitialStreamWindowSize
	}
	return nil
}

func (m *Http2ProtocolOptions) GetInitialConnectionWindowSize() *google_protobuf.UInt32Value {
	if m != nil {
		return m.InitialConnectionWindowSize
	}
	return nil
}

type GrpcProtocolOptions struct {
	Http2ProtocolOptions *Http2ProtocolOptions `protobuf:"bytes,1,opt,name=http2_protocol_options,json=http2ProtocolOptions" json:"http2_protocol_options,omitempty"`
}

func (m *GrpcProtocolOptions) Reset()                    { *m = GrpcProtocolOptions{} }
func (m *GrpcProtocolOptions) String() string            { return proto.CompactTextString(m) }
func (*GrpcProtocolOptions) ProtoMessage()               {}
func (*GrpcProtocolOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GrpcProtocolOptions) GetHttp2ProtocolOptions() *Http2ProtocolOptions {
	if m != nil {
		return m.Http2ProtocolOptions
	}
	return nil
}

func init() {
	proto.RegisterType((*TcpProtocolOptions)(nil), "envoy.api.v2.TcpProtocolOptions")
	proto.RegisterType((*Http1ProtocolOptions)(nil), "envoy.api.v2.Http1ProtocolOptions")
	proto.RegisterType((*Http2ProtocolOptions)(nil), "envoy.api.v2.Http2ProtocolOptions")
	proto.RegisterType((*GrpcProtocolOptions)(nil), "envoy.api.v2.GrpcProtocolOptions")
}

func init() { proto.RegisterFile("api/protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcb, 0x4e, 0x32, 0x31,
	0x18, 0x40, 0x03, 0xff, 0x1f, 0x17, 0xd5, 0x18, 0x52, 0x27, 0x48, 0xd0, 0x10, 0x33, 0x2b, 0x57,
	0x25, 0x0e, 0x4f, 0xa0, 0x24, 0x8a, 0x2b, 0x0d, 0x17, 0x2f, 0xab, 0xda, 0xa9, 0x15, 0x1a, 0x4b,
	0xbf, 0xa6, 0xd3, 0x61, 0x90, 0xa7, 0xf6, 0x11, 0x0c, 0xd3, 0x82, 0x8a, 0x2c, 0xd8, 0x4d, 0x3a,
	0x3d, 0xe7, 0x34, 0x5f, 0x8b, 0x30, 0x33, 0xb2, 0x6d, 0x2c, 0x38, 0xe0, 0xa0, 0x48, 0xf9, 0x81,
	0x0f, 0x84, 0x9e, 0xc1, 0x07, 0x61, 0x46, 0x92, 0x59, 0xd2, 0x6c, 0x8d, 0x01, 0xc6, 0x4a, 0xf8,
	0x4d, 0x69, 0xfe, 0xd6, 0x2e, 0x2c, 0x33, 0x46, 0xd8, 0xcc, 0xef, 0x8e, 0x23, 0x84, 0x87, 0xdc,
	0xdc, 0x07, 0xc5, 0x9d, 0x71, 0x12, 0x74, 0x16, 0xbf, 0xa0, 0xa8, 0xe7, 0x9c, 0xb9, 0xd8, 0x58,
	0xc7, 0x3d, 0x84, 0x99, 0x52, 0x50, 0x50, 0x96, 0x66, 0xa0, 0x72, 0x27, 0x68, 0x6e, 0x55, 0xa3,
	0x72, 0x56, 0x39, 0xdf, 0x4f, 0x9a, 0xc4, 0xa7, 0xc8, 0x2a, 0x45, 0xae, 0x00, 0xd4, 0x03, 0x53,
	0xb9, 0xe8, 0xd7, 0x4a, 0xea, 0x32, 0x40, 0x23, 0xab, 0xe2, 0xcf, 0xaa, 0x4f, 0x24, 0x9b, 0x89,
	0x6b, 0x54, 0x9b, 0x18, 0xc6, 0xdf, 0xa9, 0x63, 0xa9, 0x12, 0x34, 0x93, 0x0b, 0x11, 0x02, 0xa7,
	0x7f, 0x02, 0xa3, 0x5b, 0xed, 0x3a, 0x89, 0x4f, 0x1c, 0x96, 0xd4, 0x70, 0x09, 0x0d, 0xe4, 0x42,
	0xe0, 0x3e, 0xaa, 0x4f, 0xd9, 0x9c, 0x72, 0xd0, 0x3c, 0xb7, 0x56, 0x68, 0x47, 0x33, 0x67, 0x05,
	0x9b, 0x66, 0x8d, 0xea, 0x0e, 0xb6, 0x68, 0xca, 0xe6, 0xdd, 0x35, 0x3a, 0xf0, 0x24, 0x7e, 0x46,
	0x4d, 0xa9, 0xa5, 0x93, 0x4c, 0x05, 0x19, 0x2d, 0xa4, 0x7e, 0x85, 0xc2, 0x9f, 0xf2, 0xdf, 0x0e,
	0xde, 0xe3, 0xc0, 0x7b, 0xe3, 0x63, 0x49, 0x97, 0xc7, 0x65, 0xa8, 0xb5, 0x52, 0x73, 0xd0, 0x5a,
	0xf0, 0xe5, 0x34, 0x7e, 0xe9, 0xff, 0xef, 0xa0, 0x3f, 0x09, 0x8e, 0xee, 0x5a, 0xf1, 0x9d, 0x88,
	0x01, 0x1d, 0xdd, 0x58, 0xc3, 0x37, 0x07, 0xfe, 0x84, 0xea, 0x93, 0xe5, 0x45, 0xd0, 0xd5, 0x3b,
	0xa2, 0xe0, 0xff, 0x84, 0xb1, 0xc7, 0xe4, 0xe7, 0x83, 0x22, 0xdb, 0x2e, 0xad, 0x1f, 0x4d, 0xb6,
	0xac, 0xa6, 0x7b, 0xa5, 0xb1, 0xf3, 0x15, 0x00, 0x00, 0xff, 0xff, 0x26, 0x5d, 0xa0, 0x7a, 0xa6,
	0x02, 0x00, 0x00,
}
