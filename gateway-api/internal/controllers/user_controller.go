package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/services"
	"gateway-api/internal/utils/logger"
)

// UserController handles user-related requests
type UserController struct {
	cfg         *config.Config
	logger      *logger.Logger
	userService services.UserService
}

// NewUserController creates a new user controller
func NewUserController(cfg *config.Config, logger *logger.Logger) *UserController {
	userService := services.NewUserService(cfg, logger)

	return &UserController{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with OAuth provider
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.AuthRequest true "Registration request"
// @Success 201 {object} models.AuthResponse "User registered successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var request models.AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Call the user service
	resp, err := c.userService.Register(ctx, request)

	if err != nil {
		c.logger.Error("Failed to register user", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to register user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with OAuth provider
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.AuthRequest true "Login request"
// @Success 200 {object} models.AuthResponse "User logged in successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var request models.AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Call the user service
	resp, err := c.userService.Login(ctx, request)

	if err != nil {
		c.logger.Error("Failed to login user", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to login user",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetProfile gets the user's profile
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserProfile "User profile"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/me [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	fmt.Println("userid", userID)

	// Get JWT token from context
	token := ctx.GetString("jwt_token")

	// Create a new context with the JWT token
	reqCtx := context.WithValue(ctx.Request.Context(), "jwt_token", token)

	// Call the user service with the new context
	resp, err := c.userService.GetProfile(reqCtx, userID)

	if err != nil {
		c.logger.Error("Failed to get user profile", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get user profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateProfile updates the user's profile
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ProfileUpdateRequest true "Update profile request"
// @Success 200 {object} models.UserProfile "User profile updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /users/me [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var request models.ProfileUpdateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Get JWT token from context
	token := ctx.GetString("jwt_token")

	// Create a new context with the JWT token
	reqCtx := context.WithValue(ctx.Request.Context(), "jwt_token", token)

	// Call the user service with the new context
	resp, err := c.userService.UpdateProfile(reqCtx, userID, request)

	if err != nil {
		c.logger.Error("Failed to update user profile", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update user profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
