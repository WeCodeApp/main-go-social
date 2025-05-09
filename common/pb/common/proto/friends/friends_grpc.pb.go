// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: friends/friends.proto

package friends

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FriendService_SendFriendRequest_FullMethodName   = "/friends.FriendService/SendFriendRequest"
	FriendService_GetFriendRequests_FullMethodName   = "/friends.FriendService/GetFriendRequests"
	FriendService_AcceptFriendRequest_FullMethodName = "/friends.FriendService/AcceptFriendRequest"
	FriendService_RejectFriendRequest_FullMethodName = "/friends.FriendService/RejectFriendRequest"
	FriendService_GetFriends_FullMethodName          = "/friends.FriendService/GetFriends"
	FriendService_RemoveFriend_FullMethodName        = "/friends.FriendService/RemoveFriend"
	FriendService_BlockUser_FullMethodName           = "/friends.FriendService/BlockUser"
	FriendService_UnblockUser_FullMethodName         = "/friends.FriendService/UnblockUser"
	FriendService_GetBlockedUsers_FullMethodName     = "/friends.FriendService/GetBlockedUsers"
	FriendService_CheckFriendship_FullMethodName     = "/friends.FriendService/CheckFriendship"
)

// FriendServiceClient is the client API for FriendService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// FriendService provides friend-related operations
type FriendServiceClient interface {
	// SendFriendRequest sends a friend request to another user
	SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error)
	// GetFriendRequests retrieves friend requests for a user
	GetFriendRequests(ctx context.Context, in *GetFriendRequestsRequest, opts ...grpc.CallOption) (*GetFriendRequestsResponse, error)
	// AcceptFriendRequest accepts a friend request
	AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error)
	// RejectFriendRequest rejects a friend request
	RejectFriendRequest(ctx context.Context, in *RejectFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error)
	// GetFriends retrieves friends for a user
	GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error)
	// RemoveFriend removes a friend
	RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error)
	// BlockUser blocks a user
	BlockUser(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*BlockUserResponse, error)
	// UnblockUser unblocks a user
	UnblockUser(ctx context.Context, in *UnblockUserRequest, opts ...grpc.CallOption) (*UnblockUserResponse, error)
	// GetBlockedUsers retrieves blocked users for a user
	GetBlockedUsers(ctx context.Context, in *GetBlockedUsersRequest, opts ...grpc.CallOption) (*GetBlockedUsersResponse, error)
	// CheckFriendship checks if two users are friends
	CheckFriendship(ctx context.Context, in *CheckFriendshipRequest, opts ...grpc.CallOption) (*CheckFriendshipResponse, error)
}

type friendServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFriendServiceClient(cc grpc.ClientConnInterface) FriendServiceClient {
	return &friendServiceClient{cc}
}

func (c *friendServiceClient) SendFriendRequest(ctx context.Context, in *SendFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendService_SendFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) GetFriendRequests(ctx context.Context, in *GetFriendRequestsRequest, opts ...grpc.CallOption) (*GetFriendRequestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendRequestsResponse)
	err := c.cc.Invoke(ctx, FriendService_GetFriendRequests_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) AcceptFriendRequest(ctx context.Context, in *AcceptFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendService_AcceptFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) RejectFriendRequest(ctx context.Context, in *RejectFriendRequestRequest, opts ...grpc.CallOption) (*FriendRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FriendRequestResponse)
	err := c.cc.Invoke(ctx, FriendService_RejectFriendRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) GetFriends(ctx context.Context, in *GetFriendsRequest, opts ...grpc.CallOption) (*GetFriendsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFriendsResponse)
	err := c.cc.Invoke(ctx, FriendService_GetFriends_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) RemoveFriend(ctx context.Context, in *RemoveFriendRequest, opts ...grpc.CallOption) (*RemoveFriendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveFriendResponse)
	err := c.cc.Invoke(ctx, FriendService_RemoveFriend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) BlockUser(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*BlockUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BlockUserResponse)
	err := c.cc.Invoke(ctx, FriendService_BlockUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) UnblockUser(ctx context.Context, in *UnblockUserRequest, opts ...grpc.CallOption) (*UnblockUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UnblockUserResponse)
	err := c.cc.Invoke(ctx, FriendService_UnblockUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) GetBlockedUsers(ctx context.Context, in *GetBlockedUsersRequest, opts ...grpc.CallOption) (*GetBlockedUsersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBlockedUsersResponse)
	err := c.cc.Invoke(ctx, FriendService_GetBlockedUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *friendServiceClient) CheckFriendship(ctx context.Context, in *CheckFriendshipRequest, opts ...grpc.CallOption) (*CheckFriendshipResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckFriendshipResponse)
	err := c.cc.Invoke(ctx, FriendService_CheckFriendship_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FriendServiceServer is the server API for FriendService service.
// All implementations must embed UnimplementedFriendServiceServer
// for forward compatibility.
//
// FriendService provides friend-related operations
type FriendServiceServer interface {
	// SendFriendRequest sends a friend request to another user
	SendFriendRequest(context.Context, *SendFriendRequestRequest) (*FriendRequestResponse, error)
	// GetFriendRequests retrieves friend requests for a user
	GetFriendRequests(context.Context, *GetFriendRequestsRequest) (*GetFriendRequestsResponse, error)
	// AcceptFriendRequest accepts a friend request
	AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*FriendRequestResponse, error)
	// RejectFriendRequest rejects a friend request
	RejectFriendRequest(context.Context, *RejectFriendRequestRequest) (*FriendRequestResponse, error)
	// GetFriends retrieves friends for a user
	GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error)
	// RemoveFriend removes a friend
	RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error)
	// BlockUser blocks a user
	BlockUser(context.Context, *BlockUserRequest) (*BlockUserResponse, error)
	// UnblockUser unblocks a user
	UnblockUser(context.Context, *UnblockUserRequest) (*UnblockUserResponse, error)
	// GetBlockedUsers retrieves blocked users for a user
	GetBlockedUsers(context.Context, *GetBlockedUsersRequest) (*GetBlockedUsersResponse, error)
	// CheckFriendship checks if two users are friends
	CheckFriendship(context.Context, *CheckFriendshipRequest) (*CheckFriendshipResponse, error)
	mustEmbedUnimplementedFriendServiceServer()
}

// UnimplementedFriendServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFriendServiceServer struct{}

func (UnimplementedFriendServiceServer) SendFriendRequest(context.Context, *SendFriendRequestRequest) (*FriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequest not implemented")
}
func (UnimplementedFriendServiceServer) GetFriendRequests(context.Context, *GetFriendRequestsRequest) (*GetFriendRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendRequests not implemented")
}
func (UnimplementedFriendServiceServer) AcceptFriendRequest(context.Context, *AcceptFriendRequestRequest) (*FriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptFriendRequest not implemented")
}
func (UnimplementedFriendServiceServer) RejectFriendRequest(context.Context, *RejectFriendRequestRequest) (*FriendRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectFriendRequest not implemented")
}
func (UnimplementedFriendServiceServer) GetFriends(context.Context, *GetFriendsRequest) (*GetFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriends not implemented")
}
func (UnimplementedFriendServiceServer) RemoveFriend(context.Context, *RemoveFriendRequest) (*RemoveFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFriend not implemented")
}
func (UnimplementedFriendServiceServer) BlockUser(context.Context, *BlockUserRequest) (*BlockUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlockUser not implemented")
}
func (UnimplementedFriendServiceServer) UnblockUser(context.Context, *UnblockUserRequest) (*UnblockUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnblockUser not implemented")
}
func (UnimplementedFriendServiceServer) GetBlockedUsers(context.Context, *GetBlockedUsersRequest) (*GetBlockedUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockedUsers not implemented")
}
func (UnimplementedFriendServiceServer) CheckFriendship(context.Context, *CheckFriendshipRequest) (*CheckFriendshipResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckFriendship not implemented")
}
func (UnimplementedFriendServiceServer) mustEmbedUnimplementedFriendServiceServer() {}
func (UnimplementedFriendServiceServer) testEmbeddedByValue()                       {}

// UnsafeFriendServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FriendServiceServer will
// result in compilation errors.
type UnsafeFriendServiceServer interface {
	mustEmbedUnimplementedFriendServiceServer()
}

func RegisterFriendServiceServer(s grpc.ServiceRegistrar, srv FriendServiceServer) {
	// If the following call pancis, it indicates UnimplementedFriendServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FriendService_ServiceDesc, srv)
}

func _FriendService_SendFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).SendFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_SendFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).SendFriendRequest(ctx, req.(*SendFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_GetFriendRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).GetFriendRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_GetFriendRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).GetFriendRequests(ctx, req.(*GetFriendRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_AcceptFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).AcceptFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_AcceptFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).AcceptFriendRequest(ctx, req.(*AcceptFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_RejectFriendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RejectFriendRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).RejectFriendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_RejectFriendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).RejectFriendRequest(ctx, req.(*RejectFriendRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_GetFriends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).GetFriends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_GetFriends_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).GetFriends(ctx, req.(*GetFriendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_RemoveFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).RemoveFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_RemoveFriend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).RemoveFriend(ctx, req.(*RemoveFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_BlockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).BlockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_BlockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).BlockUser(ctx, req.(*BlockUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_UnblockUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnblockUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).UnblockUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_UnblockUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).UnblockUser(ctx, req.(*UnblockUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_GetBlockedUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockedUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).GetBlockedUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_GetBlockedUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).GetBlockedUsers(ctx, req.(*GetBlockedUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FriendService_CheckFriendship_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckFriendshipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FriendServiceServer).CheckFriendship(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FriendService_CheckFriendship_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FriendServiceServer).CheckFriendship(ctx, req.(*CheckFriendshipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FriendService_ServiceDesc is the grpc.ServiceDesc for FriendService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FriendService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "friends.FriendService",
	HandlerType: (*FriendServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendFriendRequest",
			Handler:    _FriendService_SendFriendRequest_Handler,
		},
		{
			MethodName: "GetFriendRequests",
			Handler:    _FriendService_GetFriendRequests_Handler,
		},
		{
			MethodName: "AcceptFriendRequest",
			Handler:    _FriendService_AcceptFriendRequest_Handler,
		},
		{
			MethodName: "RejectFriendRequest",
			Handler:    _FriendService_RejectFriendRequest_Handler,
		},
		{
			MethodName: "GetFriends",
			Handler:    _FriendService_GetFriends_Handler,
		},
		{
			MethodName: "RemoveFriend",
			Handler:    _FriendService_RemoveFriend_Handler,
		},
		{
			MethodName: "BlockUser",
			Handler:    _FriendService_BlockUser_Handler,
		},
		{
			MethodName: "UnblockUser",
			Handler:    _FriendService_UnblockUser_Handler,
		},
		{
			MethodName: "GetBlockedUsers",
			Handler:    _FriendService_GetBlockedUsers_Handler,
		},
		{
			MethodName: "CheckFriendship",
			Handler:    _FriendService_CheckFriendship_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "friends/friends.proto",
}
