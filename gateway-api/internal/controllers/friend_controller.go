package controllers

import (
	friends2 "common/pb/common/proto/friends"
	"context"
	"errors"
	"gateway-api/internal/models"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"gateway-api/internal/config"
	"gateway-api/internal/utils/logger"
)

// FriendController handles friend-related requests
type FriendController struct {
	cfg    *config.Config
	logger *logger.Logger
	client friends2.FriendServiceClient
}

// createAuthContext creates a new context with the JWT token in the metadata
func (c *FriendController) createAuthContext(ctx *gin.Context) (context.Context, error) {
	// Extract token from context
	token, exists := ctx.Get("jwt_token")
	if !exists {
		return nil, errors.New("no token found in context")
	}

	tokenStr, ok := token.(string)
	if !ok || tokenStr == "" {
		return nil, errors.New("invalid token format")
	}

	// Create metadata with authorization header
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + tokenStr,
	})

	// Create a new context with the metadata
	return metadata.NewOutgoingContext(context.Background(), md), nil
}

// NewFriendController creates a new friend controller
func NewFriendController(cfg *config.Config, logger *logger.Logger) *FriendController {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.FriendsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to friends service", err)
	}

	// Create a client
	client := friends2.NewFriendServiceClient(conn)

	return &FriendController{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// GetFriends handles retrieving friends for a user
// @Summary Get friends
// @Description Get friends with pagination
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of friends per page" default(10)
// @Success 200 {object} models.FriendsResponse "Friends list with pagination"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends [get]
func (c *FriendController) GetFriends(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.GetFriends(authCtx, &friends2.GetFriendsRequest{
		UserId: userID,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		c.logger.Error("Failed to get friends", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get friends",
		})
		return
	}

	// Convert friends to model format
	friends := make([]models.Friend, len(resp.Friends))
	for i, friend := range resp.Friends {
		friends[i] = models.Friend{
			UserID:       friend.UserId,
			Name:         friend.Name,
			Avatar:       friend.Avatar,
			Email:        friend.Email,
			FriendsSince: friend.FriendsSince,
		}
	}

	ctx.JSON(http.StatusOK, models.FriendsResponse{
		Friends:    friends,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	})
}

// SendFriendRequest handles sending a friend request
// @Summary Send a friend request
// @Description Send a friend request to another user
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.FriendRequest true "Friend request"
// @Success 201 {object} models.FriendRequestDetails "Friend request sent successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/requests [post]
func (c *FriendController) SendFriendRequest(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var request models.FriendRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.SendFriendRequest(authCtx, &friends2.SendFriendRequestRequest{
		UserId:   userID,
		FriendId: request.FriendID,
	})

	if err != nil {
		c.logger.Error("Failed to send friend request", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to send friend request",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.FriendRequestDetails{
		RequestID:      resp.RequestId,
		SenderID:       resp.SenderId,
		SenderName:     resp.SenderName,
		SenderAvatar:   resp.SenderAvatar,
		ReceiverID:     resp.ReceiverId,
		ReceiverName:   resp.ReceiverName,
		ReceiverAvatar: resp.ReceiverAvatar,
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	})
}

// GetFriendRequests handles retrieving friend requests for a user
// @Summary Get friend requests
// @Description Get friend requests with pagination
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter requests by status" Enums(pending, accepted, rejected) default(pending)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of requests per page" default(10)
// @Success 200 {object} models.FriendRequestsResponse "Friend requests with pagination"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/requests [get]
func (c *FriendController) GetFriendRequests(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	status := ctx.DefaultQuery("status", "pending")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.GetFriendRequests(authCtx, &friends2.GetFriendRequestsRequest{
		UserId: userID,
		Status: status,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		c.logger.Error("Failed to get friend requests", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get friend requests",
		})
		return
	}

	// Convert requests to model format
	requests := make([]models.FriendRequestDetails, len(resp.Requests))
	for i, request := range resp.Requests {
		requests[i] = models.FriendRequestDetails{
			RequestID:      request.RequestId,
			SenderID:       request.SenderId,
			SenderName:     request.SenderName,
			SenderAvatar:   request.SenderAvatar,
			ReceiverID:     request.ReceiverId,
			ReceiverName:   request.ReceiverName,
			ReceiverAvatar: request.ReceiverAvatar,
			Status:         request.Status,
			CreatedAt:      request.CreatedAt,
			UpdatedAt:      request.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, models.FriendRequestsResponse{
		Requests:   requests,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	})
}

// AcceptFriendRequest handles accepting a friend request
// @Summary Accept a friend request
// @Description Accept a friend request
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} models.FriendRequestDetails "Friend request accepted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Request not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/requests/{id}/accept [put]
func (c *FriendController) AcceptFriendRequest(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	requestID := ctx.Param("id")

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.AcceptFriendRequest(authCtx, &friends2.AcceptFriendRequestRequest{
		RequestId: requestID,
		UserId:    userID,
	})

	if err != nil {
		c.logger.Error("Failed to accept friend request", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to accept friend request",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.FriendRequestDetails{
		RequestID:      resp.RequestId,
		SenderID:       resp.SenderId,
		SenderName:     resp.SenderName,
		SenderAvatar:   resp.SenderAvatar,
		ReceiverID:     resp.ReceiverId,
		ReceiverName:   resp.ReceiverName,
		ReceiverAvatar: resp.ReceiverAvatar,
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	})
}

// RejectFriendRequest handles rejecting a friend request
// @Summary Reject a friend request
// @Description Reject a friend request
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} models.FriendRequestDetails "Friend request rejected successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Request not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/requests/{id}/reject [put]
func (c *FriendController) RejectFriendRequest(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	requestID := ctx.Param("id")

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.RejectFriendRequest(authCtx, &friends2.RejectFriendRequestRequest{
		RequestId: requestID,
		UserId:    userID,
	})

	if err != nil {
		c.logger.Error("Failed to reject friend request", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to reject friend request",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.FriendRequestDetails{
		RequestID:      resp.RequestId,
		SenderID:       resp.SenderId,
		SenderName:     resp.SenderName,
		SenderAvatar:   resp.SenderAvatar,
		ReceiverID:     resp.ReceiverId,
		ReceiverName:   resp.ReceiverName,
		ReceiverAvatar: resp.ReceiverAvatar,
		Status:         resp.Status,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
	})
}

// RemoveFriend handles removing a friend
// @Summary Remove a friend
// @Description Remove a friend
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "Friend ID"
// @Success 200 {object} models.SuccessResponse "Friend removed successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Friend not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/{id} [delete]
func (c *FriendController) RemoveFriend(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	friendID := ctx.Param("id")

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.RemoveFriend(authCtx, &friends2.RemoveFriendRequest{
		UserId:   userID,
		FriendId: friendID,
	})

	if err != nil {
		c.logger.Error("Failed to remove friend", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to remove friend",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: resp.Success,
	})
}

// BlockUser handles blocking a user
// @Summary Block a user
// @Description Block a user
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID to block"
// @Success 200 {object} models.SuccessResponse "User blocked successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/block/{id} [post]
func (c *FriendController) BlockUser(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	blockedUserID := ctx.Param("id")

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.BlockUser(authCtx, &friends2.BlockUserRequest{
		UserId:        userID,
		BlockedUserId: blockedUserID,
	})

	if err != nil {
		c.logger.Error("Failed to block user", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to block user",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: resp.Success,
	})
}

// UnblockUser handles unblocking a user
// @Summary Unblock a user
// @Description Unblock a user
// @Tags friends
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID to unblock"
// @Success 200 {object} models.SuccessResponse "User unblocked successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /friends/block/{id} [delete]
func (c *FriendController) UnblockUser(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	blockedUserID := ctx.Param("id")

	// Create context with authorization metadata
	authCtx, err := c.createAuthContext(ctx)
	if err != nil {
		c.logger.Error("Failed to create auth context", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to authenticate request",
		})
		return
	}

	// Call the gRPC service
	resp, err := c.client.UnblockUser(authCtx, &friends2.UnblockUserRequest{
		UserId:        userID,
		BlockedUserId: blockedUserID,
	})

	if err != nil {
		c.logger.Error("Failed to unblock user", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to unblock user",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: resp.Success,
	})
}
