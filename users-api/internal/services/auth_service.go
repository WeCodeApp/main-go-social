package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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

// AuthService defines the interface for authentication-related operations
type AuthService interface {
	// GoogleLogin generates a Google OAuth URL with state token
	GoogleLogin(ctx context.Context, redirectURL string) (string, string, error)

	// MicrosoftLogin generates a Microsoft OAuth URL with state token
	MicrosoftLogin(ctx context.Context, redirectURL string) (string, string, error)

	// GoogleCallback handles the callback from Google OAuth
	GoogleCallback(ctx context.Context, state, code string) (string, string, error)

	// MicrosoftCallback handles the callback from Microsoft OAuth
	MicrosoftCallback(ctx context.Context, state, code string) (string, string, error)

	// ValidateStateToken validates the state token to prevent CSRF attacks
	ValidateStateToken(state string) bool

	// Signout signs out the user
	Signout(ctx context.Context, token string) (bool, error)
}

// authService implements the AuthService interface
type authService struct {
	userRepo        repository.UserRepository
	logger          *logger.Logger
	jwtSecret       string
	jwtExpiration   time.Duration
	googleConfig    *oauth2.Config
	microsoftConfig *oauth2.Config
	stateStore      map[string]time.Time // Store state tokens for CSRF protection
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	logger *logger.Logger,
	jwtSecret string,
	jwtExpiration time.Duration,
	googleClientID string,
	googleClientSecret string,
	microsoftClientID string,
	microsoftClientSecret string,
) AuthService {
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

	return &authService{
		userRepo:        userRepo,
		logger:          logger,
		jwtSecret:       jwtSecret,
		jwtExpiration:   jwtExpiration,
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

// GoogleLogin generates a Google OAuth URL with state token
func (s *authService) GoogleLogin(ctx context.Context, redirectURL string) (string, string, error) {
	// Set redirect URL if provided
	if redirectURL != "" {
		s.googleConfig.RedirectURL = redirectURL
	}

	// Generate a state token for CSRF protection
	state, err := generateStateToken()
	if err != nil {
		s.logger.Error("Failed to generate state token", err)
		return "", "", err
	}

	// Store the state token with timestamp
	s.stateStore[state] = time.Now()

	// Generate Google's OAuth login URL
	url := s.googleConfig.AuthCodeURL(state)
	return url, state, nil
}

// MicrosoftLogin generates a Microsoft OAuth URL with state token
func (s *authService) MicrosoftLogin(ctx context.Context, redirectURL string) (string, string, error) {
	// Set redirect URL if provided
	if redirectURL != "" {
		s.microsoftConfig.RedirectURL = redirectURL
	}

	// Generate a state token for CSRF protection
	state, err := generateStateToken()
	if err != nil {
		s.logger.Error("Failed to generate state token", err)
		return "", "", err
	}

	// Store the state token with timestamp
	s.stateStore[state] = time.Now()

	// Generate Microsoft's OAuth login URL
	url := s.microsoftConfig.AuthCodeURL(state)
	return url, state, nil
}

// GoogleCallback handles the callback from Google OAuth
func (s *authService) GoogleCallback(ctx context.Context, state, code string) (string, string, error) {
	// Validate state token
	if !s.ValidateStateToken(state) {
		return "", "", errors.New("invalid state token")
	}

	// Exchange authorization code for token
	token, err := s.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return "", "", err
	}

	// Get user info from Google
	userInfo, err := s.getUserInfoFromOAuth(ctx, "google", token.AccessToken)
	if err != nil {
		s.logger.Error("Failed to get user info from Google", err)
		return "", "", err
	}

	// Find user by email
	existingUser, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		// User not found, create new user
		if err := s.userRepo.Create(ctx, userInfo); err != nil {
			s.logger.Error("Failed to create user", err)
			return "", "", err
		}
		existingUser = userInfo
	}

	// Generate JWT token
	accessToken, err := s.generateJWT(existingUser.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT", err)
		return "", "", err
	}

	return existingUser.ID, accessToken, nil
}

// MicrosoftCallback handles the callback from Microsoft OAuth
func (s *authService) MicrosoftCallback(ctx context.Context, state, code string) (string, string, error) {
	s.logger.Info("Microsoft callback received",
		logger.Field("state_length", len(state)),
		logger.Field("code_length", len(code)))

	// Validate state token
	if !s.ValidateStateToken(state) {
		s.logger.Error("Invalid state token", nil)
		return "", "", errors.New("invalid state token")
	}

	// Exchange authorization code for token
	s.logger.Info("Exchanging code for token",
		logger.Field("code_length", len(code)),
		logger.Field("redirect_url", s.microsoftConfig.RedirectURL),
		logger.Field("client_id", s.microsoftConfig.ClientID),
		logger.Field("scopes", s.microsoftConfig.Scopes),
		logger.Field("auth_url", s.microsoftConfig.Endpoint.AuthURL),
		logger.Field("token_url", s.microsoftConfig.Endpoint.TokenURL))

	token, err := s.microsoftConfig.Exchange(context.Background(), code)
	if err != nil {
		s.logger.Error("Failed to exchange code for token", err)
		return "", "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	s.logger.Info("Token exchange successful",
		logger.Field("token_type", token.TokenType),
		logger.Field("expiry", token.Expiry.String()),
		logger.Field("access_token_length", len(token.AccessToken)))

	// Get user info from Microsoft
	s.logger.Info("Getting user info from Microsoft")
	userInfo, err := s.getUserInfoFromOAuth(ctx, "microsoft", token.AccessToken)
	if err != nil {
		s.logger.Error("Failed to get user info from Microsoft", err)
		return "", "", fmt.Errorf("failed to get user info from Microsoft: %w", err)
	}

	s.logger.Info("User info retrieved successfully",
		logger.Field("email", userInfo.Email),
		logger.Field("name", userInfo.Name))

	// Find user by email
	s.logger.Info("Finding user by email", logger.Field("email", userInfo.Email))
	existingUser, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		s.logger.Info("User not found, creating new user", logger.Field("email", userInfo.Email))
		// User not found, create new user
		if err := s.userRepo.Create(ctx, userInfo); err != nil {
			s.logger.Error("Failed to create user", err)
			return "", "", fmt.Errorf("failed to create user: %w", err)
		}
		existingUser = userInfo
		s.logger.Info("New user created", logger.Field("id", existingUser.ID))
	} else {
		s.logger.Info("User found", logger.Field("id", existingUser.ID))
	}

	// Generate JWT token
	s.logger.Info("Generating JWT token")
	accessToken, err := s.generateJWT(existingUser.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT", err)
		return "", "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	s.logger.Info("JWT token generated successfully",
		logger.Field("token_length", len(accessToken)))

	return existingUser.ID, accessToken, nil
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
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Error("Failed to parse token", err)
		return false, err
	}

	// In a real application, you would add the token to a blacklist here

	return true, nil
}

// getUserInfoFromOAuth gets user information from OAuth provider
func (s *authService) getUserInfoFromOAuth(ctx context.Context, provider, token string) (*models.User, error) {
	if provider == "google" {
		return s.getGoogleUserInfo(ctx, token)
	} else if provider == "microsoft" {
		return s.getMicrosoftUserInfo(ctx, token)
	}
	return nil, errors.New("invalid provider")
}

// getGoogleUserInfo retrieves user information from Google OAuth
func (s *authService) getGoogleUserInfo(ctx context.Context, accessToken string) (*models.User, error) {
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
	user := &models.User{
		Provider: "google",
	}

	if sub, ok := result["sub"].(string); ok {
		user.ID = sub
	}

	if name, ok := result["name"].(string); ok {
		user.Name = name
	}

	if email, ok := result["email"].(string); ok {
		user.Email = email
	}

	if picture, ok := result["picture"].(string); ok {
		user.Avatar = picture
	}

	// Validate required fields
	if user.Email == "" {
		return nil, errors.New("incomplete user info from Google")
	}

	return user, nil
}

// getMicrosoftUserInfo retrieves user information from Microsoft OAuth
func (s *authService) getMicrosoftUserInfo(ctx context.Context, accessToken string) (*models.User, error) {
	// Microsoft Graph API endpoint for user info
	userInfoURL := "https://graph.microsoft.com/beta/me"

	s.logger.Info("Getting Microsoft user info", logger.Field("url", userInfoURL))
	s.logger.Info("Access token", logger.Field("token_length", len(accessToken)))

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	// Make the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("Failed to make request to Microsoft Graph API", err,
			logger.Field("url", userInfoURL))
		return nil, fmt.Errorf("failed to get user info from Microsoft: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("Microsoft API error", nil,
			logger.Field("status", resp.StatusCode),
			logger.Field("body", string(body)))
		return nil, fmt.Errorf("Microsoft API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Microsoft user info: %w", err)
	}

	// Extract user info
	user := &models.User{
		Provider: "microsoft",
	}

	if id, ok := result["id"].(string); ok {
		user.ID = id
	}

	if displayName, ok := result["displayName"].(string); ok {
		user.Name = displayName
	}

	if mail, ok := result["mail"].(string); ok {
		user.Email = mail
	} else if userPrincipalName, ok := result["userPrincipalName"].(string); ok {
		// Fallback to userPrincipalName if mail is not available
		user.Email = userPrincipalName
	}

	// Get photo (requires a separate API call)
	photoURL := "https://graph.microsoft.com/beta/me/photo/$value"
	s.logger.Info("Getting Microsoft user photo", logger.Field("url", photoURL))

	photoReq, err := http.NewRequestWithContext(ctx, "GET", photoURL, nil)
	if err != nil {
		s.logger.Error("Failed to create photo request", err)
		// Continue without photo, not a critical error
	} else {
		photoReq.Header.Add("Authorization", "Bearer "+accessToken)
		photoReq.Header.Add("Accept", "application/json")

		photoResp, err := client.Do(photoReq)
		if err != nil {
			s.logger.Error("Failed to get user photo", err)
			// Continue without photo, not a critical error
		} else {
			defer photoResp.Body.Close()

			if photoResp.StatusCode == http.StatusOK {
				s.logger.Info("Photo request successful")
				// If photo is available, construct a data URL
				// In a real app, you might want to save this to a CDN or file storage
				user.Avatar = fmt.Sprintf("https://graph.microsoft.com/beta/me/photo/$value")
			} else {
				body, _ := io.ReadAll(photoResp.Body)
				s.logger.Error("Failed to get user photo", nil,
					logger.Field("status", photoResp.StatusCode),
					logger.Field("body", string(body)))
			}
		}
	}

	// Validate required fields
	if user.Email == "" {
		return nil, errors.New("incomplete user info from Microsoft")
	}

	return user, nil
}

// generateJWT generates a JWT token for the user
func (s *authService) generateJWT(userID string) (string, error) {
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
