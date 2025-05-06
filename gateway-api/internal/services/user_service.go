package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	pb "common/pb/common/proto/users"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// UserService defines the interface for user-related operations
type UserService interface {
	// Register registers a new user with OAuth provider
	Register(ctx context.Context, request models.AuthRequest) (*models.AuthResponse, error)

	// Login authenticates a user with OAuth provider
	Login(ctx context.Context, request models.AuthRequest) (*models.AuthResponse, error)

	// GetProfile gets the user's profile
	GetProfile(ctx context.Context, userID string) (*models.UserProfile, error)

	// UpdateProfile updates the user's profile
	UpdateProfile(ctx context.Context, userID string, request models.ProfileUpdateRequest) (*models.UserProfile, error)

	// GoogleLogin generates a Google OAuth URL with state token
	GoogleLogin() (string, error)

	// MicrosoftLogin generates a Microsoft OAuth URL with state token
	MicrosoftLogin() (string, error)

	// GoogleCallback handles the callback from Google OAuth
	GoogleCallback(ctx context.Context, state, code string) (*models.AuthResponse, error)

	// MicrosoftCallback handles the callback from Microsoft OAuth
	MicrosoftCallback(ctx context.Context, state, code string) (*models.AuthResponse, error)

	// ValidateStateToken validates the state token to prevent CSRF attacks
	ValidateStateToken(state string) bool

	// Signout signs out the user
	Signout(ctx context.Context, token string) (bool, error)
}

// userService implements the UserService interface
type userService struct {
	cfg             *config.Config
	logger          *logger.Logger
	client          pb.UserServiceClient
	googleConfig    *oauth2.Config
	microsoftConfig *oauth2.Config
	stateStore      map[string]time.Time // Store state tokens for CSRF protection
}

// createAuthContext creates a new context with the JWT token in the metadata
func (s *userService) createAuthContext(ctx context.Context) context.Context {
	// Get JWT token from context if available
	token := ""
	if t, exists := ctx.Value("jwt_token").(string); exists {
		token = t
	}

	if token == "" {
		// If no token is provided, use the original context
		return ctx
	}

	// Create metadata with authorization header
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create a new context with the metadata
	return metadata.NewOutgoingContext(ctx, md)
}

// NewUserService creates a new user service
func NewUserService(cfg *config.Config, logger *logger.Logger) UserService {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.UsersServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to users service", err)
	}

	// Create a client
	client := pb.NewUserServiceClient(conn)

	// Configure Google OAuth2
	googleConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes:       []string{"User.Read", "profile", "email", "openid"},
		Endpoint:     google.Endpoint,
	}

	// Configure Microsoft OAuth2
	microsoftConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.Microsoft.ClientID,
		ClientSecret: cfg.OAuth.Microsoft.ClientSecret,
		RedirectURL:  cfg.OAuth.Microsoft.RedirectURL,
		Scopes:       []string{"User.Read", "profile", "email", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		},
	}

	return &userService{
		cfg:             cfg,
		logger:          logger,
		client:          client,
		googleConfig:    googleConfig,
		microsoftConfig: microsoftConfig,
		stateStore:      make(map[string]time.Time),
	}
}

// generateStateToken generates a random state token for CSRF protection
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Register registers a new user with OAuth provider
func (s *userService) Register(ctx context.Context, request models.AuthRequest) (*models.AuthResponse, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service with the auth context
	resp, err := s.client.Register(authCtx, &pb.RegisterRequest{
		Provider: request.Provider,
		Token:    request.AccessToken,
	})

	if err != nil {
		s.logger.Error("Failed to register user", err)
		return nil, err
	}

	return &models.AuthResponse{
		UserID:      resp.UserId,
		AccessToken: resp.AccessToken,
	}, nil
}

// Login authenticates a user with OAuth provider
func (s *userService) Login(ctx context.Context, request models.AuthRequest) (*models.AuthResponse, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service with the auth context
	resp, err := s.client.Login(authCtx, &pb.LoginRequest{
		Provider: request.Provider,
		Token:    request.AccessToken,
	})

	if err != nil {
		s.logger.Error("Failed to login user", err)
		return nil, err
	}

	return &models.AuthResponse{
		UserID:      resp.UserId,
		AccessToken: resp.AccessToken,
	}, nil
}

// GetProfile gets the user's profile
func (s *userService) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service with the auth context
	resp, err := s.client.GetProfile(authCtx, &pb.GetProfileRequest{
		UserId: userID,
	})

	if err != nil {
		s.logger.Error("Failed to get user profile", err)
		return nil, err
	}

	return &models.UserProfile{
		UserID:    resp.UserId,
		Name:      resp.Name,
		Email:     resp.Email,
		Avatar:    resp.Avatar,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.CreatedAt, // Using CreatedAt as UpdatedAt since it's not provided by the gRPC service
	}, nil
}

// UpdateProfile updates the user's profile
func (s *userService) UpdateProfile(ctx context.Context, userID string, request models.ProfileUpdateRequest) (*models.UserProfile, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service with the auth context
	resp, err := s.client.UpdateProfile(authCtx, &pb.UpdateProfileRequest{
		UserId: userID,
		Name:   request.Name,
		Avatar: request.Avatar,
	})

	if err != nil {
		s.logger.Error("Failed to update user profile", err)
		return nil, err
	}

	return &models.UserProfile{
		UserID:    resp.UserId,
		Name:      resp.Name,
		Email:     resp.Email,
		Avatar:    resp.Avatar,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.CreatedAt, // Using CreatedAt as UpdatedAt since it's not provided by the gRPC service
	}, nil
}

// GoogleLogin generates a Google OAuth URL with state token
func (s *userService) GoogleLogin() (string, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(context.Background())

	// Call the gRPC service with the auth context
	resp, err := s.client.GoogleLogin(authCtx, &pb.GoogleLoginRequest{
		RedirectUrl: s.googleConfig.RedirectURL,
	})
	if err != nil {
		s.logger.Error("Failed to generate Google OAuth URL", err)
		return "", err
	}

	// Store the state token with timestamp
	s.stateStore[resp.State] = time.Now()

	return resp.Url, nil
}

// MicrosoftLogin generates a Microsoft OAuth URL with state token
func (s *userService) MicrosoftLogin() (string, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(context.Background())

	// Call the gRPC service with the auth context
	resp, err := s.client.MicrosoftLogin(authCtx, &pb.MicrosoftLoginRequest{
		RedirectUrl: s.microsoftConfig.RedirectURL,
	})
	if err != nil {
		s.logger.Error("Failed to generate Microsoft OAuth URL", err)
		return "", err
	}

	// Store the state token with timestamp
	s.stateStore[resp.State] = time.Now()

	return resp.Url, nil
}

// GoogleCallback handles the callback from Google OAuth
func (s *userService) GoogleCallback(ctx context.Context, state, code string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return nil, err
	}

	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service to login the user with the auth context
	resp, err := s.client.Login(authCtx, &pb.LoginRequest{
		Provider: "google",
		Token:    token.AccessToken,
	})

	if err != nil {
		s.logger.Error("Failed to login user", err)
		return nil, err
	}

	return &models.AuthResponse{
		UserID:      resp.UserId,
		AccessToken: resp.AccessToken,
	}, nil
}

// MicrosoftCallback handles the callback from Microsoft OAuth
func (s *userService) MicrosoftCallback(ctx context.Context, state, code string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.microsoftConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return nil, err
	}

	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service to login the user with the auth context
	resp, err := s.client.Login(authCtx, &pb.LoginRequest{
		Provider: "microsoft",
		Token:    token.AccessToken,
	})

	if err != nil {
		s.logger.Error("Failed to login user", err)
		return nil, err
	}

	return &models.AuthResponse{
		UserID:      resp.UserId,
		AccessToken: resp.AccessToken,
	}, nil
}

// ValidateStateToken validates the state token to prevent CSRF attacks
func (s *userService) ValidateStateToken(state string) bool {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(context.Background())

	// Call the gRPC service with the auth context
	resp, err := s.client.ValidateStateToken(authCtx, &pb.ValidateStateTokenRequest{
		State: state,
	})

	if err != nil {
		s.logger.Error("Failed to validate state token", err)
		return false
	}

	return resp.Valid
}

// Signout signs out the user
func (s *userService) Signout(ctx context.Context, token string) (bool, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the gRPC service with the auth context
	resp, err := s.client.Signout(authCtx, &pb.SignoutRequest{
		Token: token,
	})

	if err != nil {
		s.logger.Error("Failed to sign out user", err)
		return false, err
	}

	return resp.Success, nil
}
