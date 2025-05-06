package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "common/pb/common/proto/users"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
)

// AuthService defines the interface for authentication-related operations
type AuthService interface {
	// GoogleLogin generates a Google OAuth URL with state token
	GoogleLogin(ctx context.Context) (string, error)

	// MicrosoftLogin generates a Microsoft OAuth URL with state token
	MicrosoftLogin(ctx context.Context) (string, error)

	// GoogleCallback handles the callback from Google OAuth
	GoogleCallback(ctx context.Context, state, code string) (*models.AuthResponse, error)

	// MicrosoftCallback handles the callback from Microsoft OAuth
	MicrosoftCallback(ctx context.Context, state, code string) (*models.AuthResponse, error)

	// ValidateStateToken validates the state token to prevent CSRF attacks
	ValidateStateToken(state string) bool

	// Signout signs out the user
	Signout(ctx context.Context, token string) (bool, error)
}

// authService implements the AuthService interface
type authService struct {
	cfg             *config.Config
	logger          *logger.Logger
	userService     UserService
	googleConfig    *oauth2.Config
	microsoftConfig *oauth2.Config
	stateStore      map[string]time.Time // Store state tokens for CSRF protection
	client          pb.UserServiceClient // gRPC client to the users-api
}

// createAuthContext creates a new context with the JWT token in the metadata
func (s *authService) createAuthContext(ctx context.Context) context.Context {
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

// NewAuthService creates a new auth service
func NewAuthService(cfg *config.Config, logger *logger.Logger, userService UserService) AuthService {
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

	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.UsersServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to users service", err)
	}

	// Create a client
	client := pb.NewUserServiceClient(conn)

	return &authService{
		cfg:             cfg,
		logger:          logger,
		userService:     userService,
		googleConfig:    googleConfig,
		microsoftConfig: microsoftConfig,
		stateStore:      make(map[string]time.Time),
		client:          client,
	}
}

// generateAuthStateToken generates a random state token for CSRF protection
func generateAuthStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GoogleLogin generates a Google OAuth URL with state token
func (s *authService) GoogleLogin(ctx context.Context) (string, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Forward the request to the users-api with auth context
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
func (s *authService) MicrosoftLogin(ctx context.Context) (string, error) {
	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)
	// Forward the request to the users-api with auth context
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
func (s *authService) GoogleCallback(ctx context.Context, state, code string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return nil, err
	}

	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the users-api to login the user with the auth context
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
func (s *authService) MicrosoftCallback(ctx context.Context, state, code string) (*models.AuthResponse, error) {
	// Exchange authorization code for token
	token, err := s.microsoftConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return nil, err
	}

	// Create auth request for user service
	authRequest := models.AuthRequest{
		Provider:    "microsoft",
		AccessToken: token.AccessToken,
	}

	// Create context with authorization metadata
	authCtx := s.createAuthContext(ctx)

	// Call the user service to login the user with the auth context
	return s.userService.Login(authCtx, authRequest)
}

// ValidateStateToken validates the state token to prevent CSRF attacks
func (s *authService) ValidateStateToken(state string) bool {
	timestamp, exists := s.stateStore[state]
	if !exists {
		return false
	}

	// Check if the token is expired (10 minutes)
	if time.Since(timestamp) > 10*time.Minute {
		delete(s.stateStore, state)
		return false
	}

	// Remove the token after use
	delete(s.stateStore, state)
	return true
}

// Signout signs out the user
func (s *authService) Signout(ctx context.Context, token string) (bool, error) {
	// Parse the token to get the user ID
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		s.logger.Error("Failed to parse token", err)
		return false, err
	}

	// In a real application, you would add the token to a blacklist here

	return true, nil
}
