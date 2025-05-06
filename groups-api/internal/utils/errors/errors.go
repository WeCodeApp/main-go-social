package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common errors
var (
	ErrInvalidArgument   = status.Error(codes.InvalidArgument, "invalid argument")
	ErrNotFound          = status.Error(codes.NotFound, "not found")
	ErrAlreadyExists     = status.Error(codes.AlreadyExists, "already exists")
	ErrPermissionDenied  = status.Error(codes.PermissionDenied, "permission denied")
	ErrUnauthenticated   = status.Error(codes.Unauthenticated, "unauthenticated")
	ErrInternal          = status.Error(codes.Internal, "internal error")
	ErrResourceExhausted = status.Error(codes.ResourceExhausted, "resource exhausted")
)