package controllers

import (
	"context"
	"users-api/internal/services"
	"users-api/internal/utils/logger"

	pb "common/pb/common/proto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserController handles gRPC requests for the user service
type UserController struct {
	pb.UnimplementedUserServiceServer
	userService services.UserService
	logger      *logger.Logger
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserService, logger *logger.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// Register registers a new user with OAuth provider
func (c *UserController) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	c.logger.Info("Register request received", logger.Field("provider", req.Provider))

	// Validate request
	if req.Provider == "" || req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "provider and token are required")
	}

	// Call service to register user
	userID, accessToken, err := c.userService.Register(ctx, req.Provider, req.Token)
	if err != nil {
		c.logger.Error("Failed to register user", err)
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:      userID,
		AccessToken: accessToken,
	}, nil
}

// Login authenticates a user with OAuth provider
func (c *UserController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	c.logger.Info("Login request received", logger.Field("provider", req.Provider))

	// Validate request
	if req.Provider == "" || req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "provider and token are required")
	}

	// Call service to login user
	userID, accessToken, err := c.userService.Login(ctx, req.Provider, req.Token)
	if err != nil {
		c.logger.Error("Failed to login user", err)
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &pb.LoginResponse{
		UserId:      userID,
		AccessToken: accessToken,
	}, nil
}

// GetProfile retrieves a user's profile
func (c *UserController) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	c.logger.Info("GetProfile request received", logger.Field("user_id", req.UserId))

	// Validate request
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	// Get user ID from context if not provided in request
	if req.UserId == "" {
		userID, ok := ctx.Value("userID").(string)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
		}
		req.UserId = userID
	}

	// Call service to get user profile
	user, err := c.userService.GetProfile(ctx, req.UserId)
	if err != nil {
		c.logger.Error("Failed to get user profile", err)
		return nil, status.Errorf(codes.Internal, "failed to get user profile: %v", err)
	}

	return &pb.ProfileResponse{
		UserId:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// UpdateProfile updates a user's profile
func (c *UserController) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.ProfileResponse, error) {
	c.logger.Info("UpdateProfile request received", logger.Field("user_id", req.UserId))

	// Validate request
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID is required")
	}

	// Get user ID from context if not provided in request
	if req.UserId == "" {
		userID, ok := ctx.Value("userID").(string)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "user ID not found in context")
		}
		req.UserId = userID
	}

	// Call service to update user profile
	user, err := c.userService.UpdateProfile(ctx, req.UserId, req.Name, req.Avatar)
	if err != nil {
		c.logger.Error("Failed to update user profile", err)
		return nil, status.Errorf(codes.Internal, "failed to update user profile: %v", err)
	}

	return &pb.ProfileResponse{
		UserId:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
