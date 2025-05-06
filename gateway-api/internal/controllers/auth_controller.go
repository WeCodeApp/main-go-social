package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/services"
	"gateway-api/internal/utils/logger"
)

// AuthController handles authentication-related requests
type AuthController struct {
	cfg         *config.Config
	logger      *logger.Logger
	authService services.AuthService
	userService services.UserService
}

// NewAuthController creates a new auth controller
func NewAuthController(cfg *config.Config, logger *logger.Logger, authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		cfg:         cfg,
		logger:      logger,
		authService: authService,
		userService: userService,
	}
}

// GoogleLogin initiates Google OAuth login
// @Summary Initiate Google OAuth login
// @Description Redirects the user to Google's OAuth login page
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/google [get]
func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	// Call the auth service
	url, err := c.authService.GoogleLogin(ctx)

	if err != nil {
		c.logger.Error("Failed to generate Google login URL", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to initiate Google login",
		})
		return
	}

	//ctx.Redirect(http.StatusFound, url)

	// Return the login URL as JSON
	ctx.JSON(http.StatusOK, gin.H{"login_url": url})
}

// MicrosoftLogin initiates Microsoft OAuth login
// @Summary Initiate Microsoft OAuth login
// @Description Redirects the user to Microsoft's OAuth login page
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to Microsoft"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/microsoft [get]
func (c *AuthController) MicrosoftLogin(ctx *gin.Context) {
	// Call the auth service
	loginUrl, err := c.authService.MicrosoftLogin(ctx)
	if err != nil {
		c.logger.Error("Failed to generate Microsoft login URL", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to initiate Microsoft login",
		})
		return
	}

	//ctx.Redirect(http.StatusFound, url)
	// Return the login URL as JSON
	ctx.JSON(http.StatusOK, gin.H{"login_url": loginUrl})
}

// MicrosoftCallback handles the callback from Microsoft OAuth
// @Summary Handle Microsoft OAuth callback
// @Description Handles the callback from Microsoft's OAuth login
// @Tags auth
// @Produce json
// @Param state query string true "State token for CSRF protection"
// @Param code query string true "Authorization code"
// @Param redirect_url query string false "URL to redirect to after authentication"
// @Success 302 {string} string "Redirect to frontend with token and user data"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/microsoft/callback [get]
func (c *AuthController) MicrosoftCallback(ctx *gin.Context) {
	// Get state, code, and redirect_url from query parameters
	state := ctx.Query("state")
	code := ctx.Query("code")
	redirectURLStr := ctx.Query("redirect_url")

	// Validate state token to prevent CSRF attacks
	if !c.authService.ValidateStateToken(state) {
		c.logger.Error("Invalid state token", nil)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid state token",
		})
		return
	}

	// Call the auth service
	resp, err := c.authService.MicrosoftCallback(ctx, state, code)
	if err != nil {
		c.logger.Error("Failed to handle Microsoft callback", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate with Microsoft",
		})
		return
	}

	// Create a new context with the JWT token
	ctxWithToken := context.WithValue(context.Background(), "jwt_token", resp.AccessToken)

	// Get user profile from user service
	userProfile, err := c.userService.GetProfile(ctxWithToken, resp.UserID)
	if err != nil {
		c.logger.Error("Failed to get user profile", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get user profile",
		})
		return
	}

	// Convert token and user to JSON
	tokenJSON, err := json.Marshal(resp)
	if err != nil {
		c.logger.Error("Failed to marshal token", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	userJSON, err := json.Marshal(userProfile)
	if err != nil {
		c.logger.Error("Failed to marshal user", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	// Determine the redirect URL
	var redirectURL *url.URL
	var parseErr error

	if redirectURLStr != "" {
		// Validate the redirect URL
		redirectURL, parseErr = url.Parse(redirectURLStr)
		if parseErr != nil || redirectURL.Scheme == "" || redirectURL.Host == "" {
			c.logger.Error("Invalid redirect URL", parseErr)
			// Use a hardcoded default URL
			redirectURL, parseErr = url.Parse(c.cfg.AppURL + "/login")
			if parseErr != nil {
				c.logger.Error("Failed to parse default redirect URL", parseErr)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
				return
			}
		}
	} else {
		// Use a hardcoded default URL
		redirectURL, parseErr = url.Parse(c.cfg.AppURL + "/login")
		if parseErr != nil {
			c.logger.Error("Failed to parse default redirect URL", parseErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
			return
		}
	}

	// Add token and user data to query parameters
	query := redirectURL.Query()
	query.Set("token", string(tokenJSON))
	query.Set("user", string(userJSON))
	redirectURL.RawQuery = query.Encode()

	// Redirect to frontend application with token and user data
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}

// GoogleCallback handles the callback from Google OAuth
// @Summary Handle Google OAuth callback
// @Description Handles the callback from Google's OAuth login
// @Tags auth
// @Produce json
// @Param state query string true "State token for CSRF protection"
// @Param code query string true "Authorization code"
// @Param redirect_url query string false "URL to redirect to after authentication"
// @Success 302 {string} string "Redirect to frontend with token and user data"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/google/callback [get]
func (c *AuthController) GoogleCallback(ctx *gin.Context) {
	// Get state, code, and redirect_url from query parameters
	state := ctx.Query("state")
	code := ctx.Query("code")
	redirectURLStr := ctx.Query("redirect_url")

	// Validate state token to prevent CSRF attacks
	if !c.authService.ValidateStateToken(state) {
		c.logger.Error("Invalid state token", nil)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid state token",
		})
		return
	}

	// Call the auth service
	resp, err := c.authService.GoogleCallback(ctx, state, code)
	if err != nil {
		c.logger.Error("Failed to handle Google callback", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate with Google",
		})
		return
	}

	// Create a new context with the JWT token
	ctxWithToken := context.WithValue(context.Background(), "jwt_token", resp.AccessToken)

	// Get user profile from user service
	userProfile, err := c.userService.GetProfile(ctxWithToken, resp.UserID)
	if err != nil {
		c.logger.Error("Failed to get user profile", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get user profile",
		})
		return
	}

	// Convert token and user to JSON
	tokenJSON, err := json.Marshal(resp)
	if err != nil {
		c.logger.Error("Failed to marshal token", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	userJSON, err := json.Marshal(userProfile)
	if err != nil {
		c.logger.Error("Failed to marshal user", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
		return
	}

	// Determine the redirect URL
	var redirectURL *url.URL
	var parseErr error

	if redirectURLStr != "" {
		// Validate the redirect URL
		redirectURL, parseErr = url.Parse(redirectURLStr)
		if parseErr != nil || redirectURL.Scheme == "" || redirectURL.Host == "" {
			c.logger.Error("Invalid redirect URL", parseErr)
			// Use a hardcoded default URL
			redirectURL, parseErr = url.Parse(c.cfg.AppURL + "/login")
			if parseErr != nil {
				c.logger.Error("Failed to parse default redirect URL", parseErr)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
				return
			}
		}
	} else {
		// Use a hardcoded default URL
		redirectURL, parseErr = url.Parse(c.cfg.AppURL + "/login")
		if parseErr != nil {
			c.logger.Error("Failed to parse default redirect URL", parseErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication"})
			return
		}
	}

	// Add token and user data to query parameters
	query := redirectURL.Query()
	query.Set("token", string(tokenJSON))
	query.Set("user", string(userJSON))
	redirectURL.RawQuery = query.Encode()

	// Redirect to frontend application with token and user data
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}

// Signout signs out the user
// @Summary Sign out the user
// @Description Signs out the user by invalidating the token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.SuccessResponse "User signed out successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/signout [post]
func (c *AuthController) Signout(ctx *gin.Context) {
	// Get the token from the Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Call the auth service
	success, err := c.authService.Signout(ctx, token)

	if err != nil {
		c.logger.Error("Failed to sign out user", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to sign out user",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: success,
	})
}
