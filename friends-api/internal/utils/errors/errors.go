package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Common errors
var (
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "unauthenticated")
	ErrNotFound         = status.Error(codes.NotFound, "not found")
	ErrInvalidArgument  = status.Error(codes.InvalidArgument, "invalid argument")
	ErrInternal         = status.Error(codes.Internal, "internal error")
	ErrAlreadyExists    = status.Error(codes.AlreadyExists, "already exists")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
)
