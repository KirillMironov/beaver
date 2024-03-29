// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: api/storage.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	Upload(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (Storage_UploadClient, error)
	Download(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (Storage_DownloadClient, error)
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResponse, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) Upload(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (Storage_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[0], "/proto.Storage/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_UploadClient interface {
	Recv() (*File, error)
	grpc.ClientStream
}

type storageUploadClient struct {
	grpc.ClientStream
}

func (x *storageUploadClient) Recv() (*File, error) {
	m := new(File)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) Download(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (Storage_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[1], "/proto.Storage/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadClient interface {
	Recv() (*File, error)
	grpc.ClientStream
}

type storageDownloadClient struct {
	grpc.ClientStream
}

func (x *storageDownloadClient) Recv() (*File, error) {
	m := new(File)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/proto.Storage/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
// All implementations should embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	Upload(*FileRequest, Storage_UploadServer) error
	Download(*FileRequest, Storage_DownloadServer) error
	List(context.Context, *emptypb.Empty) (*ListResponse, error)
}

// UnimplementedStorageServer should be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) Upload(*FileRequest, Storage_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedStorageServer) Download(*FileRequest, Storage_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedStorageServer) List(context.Context, *emptypb.Empty) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).Upload(m, &storageUploadServer{stream})
}

type Storage_UploadServer interface {
	Send(*File) error
	grpc.ServerStream
}

type storageUploadServer struct {
	grpc.ServerStream
}

func (x *storageUploadServer) Send(m *File) error {
	return x.ServerStream.SendMsg(m)
}

func _Storage_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).Download(m, &storageDownloadServer{stream})
}

type Storage_DownloadServer interface {
	Send(*File) error
	grpc.ServerStream
}

type storageDownloadServer struct {
	grpc.ServerStream
}

func (x *storageDownloadServer) Send(m *File) error {
	return x.ServerStream.SendMsg(m)
}

func _Storage_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Storage/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).List(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Storage_List_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _Storage_Upload_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _Storage_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/storage.proto",
}
