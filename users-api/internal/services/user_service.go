package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"users-api/internal/models"
	"users-api/internal/repository"
	"users-api/internal/utils/logger"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// UserService defines the interface for user service operations
type UserService interface {
	Register(ctx context.Context, provider, token string) (string, string, error)
	Login(ctx context.Context, provider, token string) (string, string, error)
	GetProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, userID, name, avatar string) (*models.User, error)
}

// userService implements the UserService interface
type userService struct {
	userRepo        repository.UserRepository
	logger          *logger.Logger
	jwtSecret       string
	jwtExpiration   time.Duration
	googleConfig    *oauth2.Config
	microsoftConfig *oauth2.Config
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	logger *logger.Logger,
	jwtSecret string,
	jwtExpiration time.Duration,
	googleClientID string,
	googleClientSecret string,
	microsoftClientID string,
	microsoftClientSecret string,
) UserService {
	// Configure Google OAuth2
	googleConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"User.Read", "profile", "email", "openid"},
	}

	// Configure Microsoft OAuth2
	microsoftConfig := &oauth2.Config{
		ClientID:     microsoftClientID,
		ClientSecret: microsoftClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/token",
		},
		Scopes: []string{"User.Read", "profile", "email", "openid"},
	}

	return &userService{
		userRepo:        userRepo,
		logger:          logger,
		jwtSecret:       jwtSecret,
		jwtExpiration:   jwtExpiration,
		googleConfig:    googleConfig,
		microsoftConfig: microsoftConfig,
	}
}

// Register registers a new user with OAuth provider
func (s *userService) Register(ctx context.Context, provider, token string) (string, string, error) {
	// Validate provider
	if provider != "google" && provider != "microsoft" {
		return "", "", errors.New("invalid provider")
	}

	// Get user info from OAuth provider
	userInfo, err := s.getUserInfoFromOAuth(ctx, provider, token)
	if err != nil {
		s.logger.Error("Failed to get user info from OAuth", err)
		return "", "", err
	}

	// Check if user already exists
	_, err = s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err == nil {
		// User already exists, return error
		return "", "", errors.New("user already exists")
	}

	// Create new user
	user := &models.User{
		Name:     userInfo.Name,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		Provider: provider,
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", err)
		return "", "", err
	}

	// Generate JWT token
	accessToken, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT", err)
		return "", "", err
	}

	return user.ID, accessToken, nil
}

// Login authenticates a user with OAuth provider
func (s *userService) Login(ctx context.Context, provider, token string) (string, string, error) {
	// Validate provider
	if provider != "google" && provider != "microsoft" {
		return "", "", errors.New("invalid provider")
	}

	// Get user info from OAuth provider
	userInfo, err := s.getUserInfoFromOAuth(ctx, provider, token)
	if err != nil {
		s.logger.Error("Failed to get user info from OAuth", err)
		return "", "", err
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		// User not found, register new user
		return s.Register(ctx, provider, token)
	}

	// Generate JWT token
	accessToken, err := s.generateJWT(user.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT", err)
		return "", "", err
	}

	return user.ID, accessToken, nil
}

// GetProfile retrieves a user's profile
func (s *userService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// UpdateProfile updates a user's profile
func (s *userService) UpdateProfile(ctx context.Context, userID, name, avatar string) (*models.User, error) {
	// Find user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update user fields
	if name != "" {
		user.Name = name
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	// Save user to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("Failed to update user", err)
		return nil, err
	}

	return user, nil
}

// UserInfo represents user information from OAuth provider
type UserInfo struct {
	ID     string
	Name   string
	Email  string
	Avatar string
}

// getUserInfoFromOAuth gets user information from OAuth provider
func (s *userService) getUserInfoFromOAuth(ctx context.Context, provider, token string) (*UserInfo, error) {
	if provider == "google" {
		return s.getGoogleUserInfo(ctx, token)
	} else if provider == "microsoft" {
		return s.getMicrosoftUserInfo(ctx, token)
	}

	return nil, errors.New("invalid provider")
}

// getGoogleUserInfo retrieves user information from Google OAuth
func (s *userService) getGoogleUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// Google's userinfo endpoint
	userInfoURL := "https://www.googleapis.com/oauth2/v3/userinfo"

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Google user info: %w", err)
	}

	// Extract user info
	userInfo := &UserInfo{}

	if sub, ok := result["sub"].(string); ok {
		userInfo.ID = sub
	}

	if name, ok := result["name"].(string); ok {
		userInfo.Name = name
	}

	if email, ok := result["email"].(string); ok {
		userInfo.Email = email
	}

	if picture, ok := result["picture"].(string); ok {
		userInfo.Avatar = picture
	}

	// Validate required fields
	if userInfo.ID == "" || userInfo.Email == "" {
		return nil, errors.New("incomplete user info from Google")
	}

	return userInfo, nil
}

// getMicrosoftUserInfo retrieves user information from Microsoft OAuth
func (s *userService) getMicrosoftUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// Microsoft Graph API endpoint for user info
	userInfoURL := "https://graph.microsoft.com/v1.0/me"

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Microsoft: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Microsoft API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Microsoft user info: %w", err)
	}

	// Extract user info
	userInfo := &UserInfo{}

	if id, ok := result["id"].(string); ok {
		userInfo.ID = id
	}

	if displayName, ok := result["displayName"].(string); ok {
		userInfo.Name = displayName
	}

	if mail, ok := result["mail"].(string); ok {
		userInfo.Email = mail
	} else if userPrincipalName, ok := result["userPrincipalName"].(string); ok {
		// Fallback to userPrincipalName if mail is not available
		userInfo.Email = userPrincipalName
	}

	// Get photo (requires a separate API call)
	photoURL := "https://graph.microsoft.com/v1.0/me/photo/$value"
	photoReq, err := http.NewRequestWithContext(ctx, "GET", photoURL, nil)
	if err == nil {
		photoReq.Header.Add("Authorization", "Bearer "+accessToken)
		photoResp, err := client.Do(photoReq)
		if err == nil && photoResp.StatusCode == http.StatusOK {
			// If photo is available, construct a data URL
			// In a real app, you might want to save this to a CDN or file storage
			userInfo.Avatar = fmt.Sprintf("https://graph.microsoft.com/v1.0/me/photo/$value")
			photoResp.Body.Close()
		}
	}

	// Validate required fields
	if userInfo.ID == "" || userInfo.Email == "" {
		return nil, errors.New("incomplete user info from Microsoft")
	}

	return userInfo, nil
}

// generateJWT generates a JWT token for the user
func (s *userService) generateJWT(userID string) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                                 // Subject (user ID)
		"iat": time.Now().Unix(),                      // Issued at
		"exp": time.Now().Add(s.jwtExpiration).Unix(), // Expiration
	})

	// Sign token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
