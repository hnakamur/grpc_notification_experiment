// Code generated by protoc-gen-go.
// source: sites.proto
// DO NOT EDIT!

/*
Package sites is a generated protocol buffer package.

It is generated from these files:
	sites.proto

It has these top-level messages:
	Empty
	Site
	Sites
	SiteModification
*/
package sites

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type SiteModificationOp int32

const (
	SiteModificationOp_UNKNOWN SiteModificationOp = 0
	SiteModificationOp_ADDED   SiteModificationOp = 1
	SiteModificationOp_EDITED  SiteModificationOp = 2
	SiteModificationOp_REMOVED SiteModificationOp = 3
)

var SiteModificationOp_name = map[int32]string{
	0: "UNKNOWN",
	1: "ADDED",
	2: "EDITED",
	3: "REMOVED",
}
var SiteModificationOp_value = map[string]int32{
	"UNKNOWN": 0,
	"ADDED":   1,
	"EDITED":  2,
	"REMOVED": 3,
}

func (x SiteModificationOp) String() string {
	return proto.EnumName(SiteModificationOp_name, int32(x))
}
func (SiteModificationOp) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// See http://stackoverflow.com/questions/31768665/can-i-define-a-grpc-call-with-a-null-request-or-response/31772973#31772973 for Empty
type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Site struct {
	Domain string `protobuf:"bytes,1,opt,name=domain" json:"domain,omitempty"`
	Origin string `protobuf:"bytes,2,opt,name=origin" json:"origin,omitempty"`
}

func (m *Site) Reset()                    { *m = Site{} }
func (m *Site) String() string            { return proto.CompactTextString(m) }
func (*Site) ProtoMessage()               {}
func (*Site) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Sites struct {
	Sites []*Site `protobuf:"bytes,1,rep,name=sites" json:"sites,omitempty"`
}

func (m *Sites) Reset()                    { *m = Sites{} }
func (m *Sites) String() string            { return proto.CompactTextString(m) }
func (*Sites) ProtoMessage()               {}
func (*Sites) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Sites) GetSites() []*Site {
	if m != nil {
		return m.Sites
	}
	return nil
}

type SiteModification struct {
	Op   SiteModificationOp `protobuf:"varint,1,opt,name=op,enum=sites.SiteModificationOp" json:"op,omitempty"`
	Site *Site              `protobuf:"bytes,2,opt,name=site" json:"site,omitempty"`
}

func (m *SiteModification) Reset()                    { *m = SiteModification{} }
func (m *SiteModification) String() string            { return proto.CompactTextString(m) }
func (*SiteModification) ProtoMessage()               {}
func (*SiteModification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SiteModification) GetSite() *Site {
	if m != nil {
		return m.Site
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "sites.Empty")
	proto.RegisterType((*Site)(nil), "sites.Site")
	proto.RegisterType((*Sites)(nil), "sites.Sites")
	proto.RegisterType((*SiteModification)(nil), "sites.SiteModification")
	proto.RegisterEnum("sites.SiteModificationOp", SiteModificationOp_name, SiteModificationOp_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for SitesService service

type SitesServiceClient interface {
	ListSites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Sites, error)
	NotifySiteModification(ctx context.Context, in *SiteModification, opts ...grpc.CallOption) (*Empty, error)
	WatchSites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SitesService_WatchSitesClient, error)
}

type sitesServiceClient struct {
	cc *grpc.ClientConn
}

func NewSitesServiceClient(cc *grpc.ClientConn) SitesServiceClient {
	return &sitesServiceClient{cc}
}

func (c *sitesServiceClient) ListSites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Sites, error) {
	out := new(Sites)
	err := grpc.Invoke(ctx, "/sites.SitesService/ListSites", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sitesServiceClient) NotifySiteModification(ctx context.Context, in *SiteModification, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/sites.SitesService/NotifySiteModification", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sitesServiceClient) WatchSites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SitesService_WatchSitesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SitesService_serviceDesc.Streams[0], c.cc, "/sites.SitesService/WatchSites", opts...)
	if err != nil {
		return nil, err
	}
	x := &sitesServiceWatchSitesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SitesService_WatchSitesClient interface {
	Recv() (*SiteModification, error)
	grpc.ClientStream
}

type sitesServiceWatchSitesClient struct {
	grpc.ClientStream
}

func (x *sitesServiceWatchSitesClient) Recv() (*SiteModification, error) {
	m := new(SiteModification)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SitesService service

type SitesServiceServer interface {
	ListSites(context.Context, *Empty) (*Sites, error)
	NotifySiteModification(context.Context, *SiteModification) (*Empty, error)
	WatchSites(*Empty, SitesService_WatchSitesServer) error
}

func RegisterSitesServiceServer(s *grpc.Server, srv SitesServiceServer) {
	s.RegisterService(&_SitesService_serviceDesc, srv)
}

func _SitesService_ListSites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SitesServiceServer).ListSites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sites.SitesService/ListSites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SitesServiceServer).ListSites(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SitesService_NotifySiteModification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SiteModification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SitesServiceServer).NotifySiteModification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sites.SitesService/NotifySiteModification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SitesServiceServer).NotifySiteModification(ctx, req.(*SiteModification))
	}
	return interceptor(ctx, in, info, handler)
}

func _SitesService_WatchSites_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SitesServiceServer).WatchSites(m, &sitesServiceWatchSitesServer{stream})
}

type SitesService_WatchSitesServer interface {
	Send(*SiteModification) error
	grpc.ServerStream
}

type sitesServiceWatchSitesServer struct {
	grpc.ServerStream
}

func (x *sitesServiceWatchSitesServer) Send(m *SiteModification) error {
	return x.ServerStream.SendMsg(m)
}

var _SitesService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sites.SitesService",
	HandlerType: (*SitesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSites",
			Handler:    _SitesService_ListSites_Handler,
		},
		{
			MethodName: "NotifySiteModification",
			Handler:    _SitesService_NotifySiteModification_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchSites",
			Handler:       _SitesService_WatchSites_Handler,
			ServerStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x52, 0xcd, 0x4e, 0xf2, 0x40,
	0x14, 0x65, 0x80, 0x42, 0xb8, 0x25, 0x5f, 0xc8, 0x2c, 0xf8, 0x94, 0xc4, 0xa8, 0x5d, 0x21, 0x89,
	0x8d, 0xa9, 0x89, 0xae, 0x21, 0xed, 0xc2, 0x28, 0x2d, 0x29, 0x2a, 0x3b, 0x93, 0x5a, 0x4a, 0x99,
	0xc4, 0x76, 0x26, 0xd3, 0xd1, 0xc8, 0x83, 0xf9, 0x7e, 0xce, 0x4c, 0xbb, 0xa8, 0x22, 0xbb, 0x9e,
	0x7b, 0xcf, 0xcf, 0x3d, 0x6d, 0xc1, 0x2c, 0x88, 0x48, 0x0a, 0x9b, 0x71, 0x2a, 0x28, 0x36, 0x34,
	0xb0, 0xba, 0x60, 0x78, 0x19, 0x13, 0x3b, 0xeb, 0x06, 0xda, 0x4b, 0x39, 0xc1, 0x43, 0xe8, 0xac,
	0x69, 0x16, 0x91, 0xfc, 0x08, 0x9d, 0xa1, 0x71, 0x2f, 0xac, 0x90, 0x9a, 0x53, 0x4e, 0x52, 0x39,
	0x6f, 0x96, 0xf3, 0x12, 0x59, 0x13, 0x30, 0x94, 0xae, 0xc0, 0xe7, 0x50, 0x5a, 0x4a, 0x5d, 0x6b,
	0x6c, 0x3a, 0xa6, 0x5d, 0xa6, 0xa9, 0x65, 0x58, 0x85, 0xbd, 0xc0, 0x40, 0xc1, 0x39, 0x5d, 0x93,
	0x0d, 0x89, 0x23, 0x41, 0x68, 0x8e, 0x2f, 0xa0, 0x49, 0x99, 0xce, 0xfa, 0xe7, 0x1c, 0xd7, 0x34,
	0x75, 0x52, 0xc0, 0x42, 0x49, 0xc2, 0xa7, 0xd0, 0x56, 0x7b, 0x7d, 0xc0, 0xaf, 0x00, 0xbd, 0x98,
	0x78, 0x80, 0xf7, 0xa5, 0xd8, 0x84, 0xee, 0x93, 0x7f, 0xef, 0x07, 0x2b, 0x7f, 0xd0, 0xc0, 0x3d,
	0x30, 0xa6, 0xae, 0xeb, 0xb9, 0x03, 0x84, 0x01, 0x3a, 0x9e, 0x7b, 0xf7, 0x28, 0x9f, 0x9b, 0x8a,
	0x13, 0x7a, 0xf3, 0xe0, 0x59, 0x82, 0x96, 0xf3, 0x85, 0xa0, 0xaf, 0x3b, 0x2d, 0x13, 0xfe, 0x41,
	0xe2, 0x44, 0xde, 0xd8, 0x7b, 0x20, 0x85, 0x28, 0x7b, 0xf6, 0xab, 0x5c, 0xfd, 0xda, 0x46, 0xfd,
	0xda, 0x15, 0x85, 0xd5, 0xc0, 0x53, 0x18, 0xfa, 0x54, 0x90, 0xcd, 0x6e, 0xaf, 0xe8, 0xff, 0x03,
	0xe5, 0x46, 0x3f, 0x0c, 0xa5, 0xc5, 0x2d, 0xc0, 0x2a, 0x12, 0xf1, 0xf6, 0xaf, 0xb8, 0x43, 0x26,
	0x56, 0xe3, 0x0a, 0xcd, 0x2e, 0xe1, 0x84, 0xf2, 0xd4, 0xce, 0x23, 0xfe, 0xbe, 0xb5, 0x53, 0xce,
	0x62, 0x3b, 0xf9, 0x8c, 0x32, 0xf6, 0x26, 0xf9, 0x5a, 0x35, 0x03, 0x6d, 0xb9, 0x50, 0xdf, 0x7f,
	0x81, 0x5e, 0x3b, 0xfa, 0x47, 0xb8, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x82, 0xa7, 0x56, 0x5b,
	0x17, 0x02, 0x00, 0x00,
}
