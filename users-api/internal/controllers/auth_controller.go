package controllers

import (
	"context"
	"users-api/internal/services"
	"users-api/internal/utils/logger"

	pb "common/pb/common/proto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthController handles gRPC requests for the auth service
type AuthController struct {
	pb.UnimplementedUserServiceServer
	authService    services.AuthService
	userController *UserController // Embed the user controller to delegate existing RPC methods
	logger         *logger.Logger
}

// NewAuthController creates a new auth controller
func NewAuthController(authService services.AuthService, userController *UserController, logger *logger.Logger) *AuthController {
	return &AuthController{
		authService:    authService,
		userController: userController,
		logger:         logger,
	}
}

// Register delegates to the user controller
func (c *AuthController) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return c.userController.Register(ctx, req)
}

// Login delegates to the user controller
func (c *AuthController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return c.userController.Login(ctx, req)
}

// GetProfile delegates to the user controller
func (c *AuthController) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.ProfileResponse, error) {
	return c.userController.GetProfile(ctx, req)
}

// UpdateProfile delegates to the user controller
func (c *AuthController) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.ProfileResponse, error) {
	return c.userController.UpdateProfile(ctx, req)
}

// GoogleLogin generates a Google OAuth URL with state token
func (c *AuthController) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.OAuthURLResponse, error) {
	c.logger.Info("GoogleLogin request received")

	// Call service to generate Google OAuth URL
	url, state, err := c.authService.GoogleLogin(ctx, req.RedirectUrl)
	if err != nil {
		c.logger.Error("Failed to generate Google OAuth URL", err)
		return nil, status.Errorf(codes.Internal, "failed to generate Google OAuth URL: %v", err)
	}

	return &pb.OAuthURLResponse{
		Url:   url,
		State: state,
	}, nil
}

// MicrosoftLogin generates a Microsoft OAuth URL with state token
func (c *AuthController) MicrosoftLogin(ctx context.Context, req *pb.MicrosoftLoginRequest) (*pb.OAuthURLResponse, error) {
	c.logger.Info("MicrosoftLogin request received")

	// Call service to generate Microsoft OAuth URL
	url, state, err := c.authService.MicrosoftLogin(ctx, req.RedirectUrl)
	if err != nil {
		c.logger.Error("Failed to generate Microsoft OAuth URL", err)
		return nil, status.Errorf(codes.Internal, "failed to generate Microsoft OAuth URL: %v", err)
	}

	return &pb.OAuthURLResponse{
		Url:   url,
		State: state,
	}, nil
}

// GoogleCallback handles the callback from Google OAuth
func (c *AuthController) GoogleCallback(ctx context.Context, req *pb.OAuthCallbackRequest) (*pb.LoginResponse, error) {
	c.logger.Info("GoogleCallback request received")

	// Validate request
	if req.State == "" || req.Code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "state and code are required")
	}

	// Call service to handle Google callback
	userID, accessToken, err := c.authService.GoogleCallback(ctx, req.State, req.Code)
	if err != nil {
		c.logger.Error("Failed to handle Google callback", err)
		return nil, status.Errorf(codes.Internal, "failed to handle Google callback: %v", err)
	}

	return &pb.LoginResponse{
		UserId:      userID,
		AccessToken: accessToken,
	}, nil
}

// MicrosoftCallback handles the callback from Microsoft OAuth
func (c *AuthController) MicrosoftCallback(ctx context.Context, req *pb.OAuthCallbackRequest) (*pb.LoginResponse, error) {
	c.logger.Info("MicrosoftCallback request received")

	// Validate request
	if req.State == "" || req.Code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "state and code are required")
	}

	// Call service to handle Microsoft callback
	userID, accessToken, err := c.authService.MicrosoftCallback(ctx, req.State, req.Code)
	if err != nil {
		c.logger.Error("Failed to handle Microsoft callback", err)
		return nil, status.Errorf(codes.Internal, "failed to handle Microsoft callback: %v", err)
	}

	return &pb.LoginResponse{
		UserId:      userID,
		AccessToken: accessToken,
	}, nil
}

// ValidateStateToken validates the state token to prevent CSRF attacks
func (c *AuthController) ValidateStateToken(ctx context.Context, req *pb.ValidateStateTokenRequest) (*pb.ValidateStateTokenResponse, error) {
	c.logger.Info("ValidateStateToken request received")

	// Validate request
	if req.State == "" {
		return nil, status.Errorf(codes.InvalidArgument, "state is required")
	}

	// Call service to validate state token
	valid := c.authService.ValidateStateToken(req.State)

	return &pb.ValidateStateTokenResponse{
		Valid: valid,
	}, nil
}

// Signout signs out the user
func (c *AuthController) Signout(ctx context.Context, req *pb.SignoutRequest) (*pb.SignoutResponse, error) {
	c.logger.Info("Signout request received")

	// Validate request
	if req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "token is required")
	}

	// Call service to sign out user
	success, err := c.authService.Signout(ctx, req.Token)
	if err != nil {
		c.logger.Error("Failed to sign out user", err)
		return nil, status.Errorf(codes.Internal, "failed to sign out user: %v", err)
	}

	return &pb.SignoutResponse{
		Success: success,
	}, nil
}
