package services

import (
	"context"
	"errors"

	pb "common/pb/common/proto/friends"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// FriendService defines the interface for friend-related operations
type FriendService interface {
	// GetFriends retrieves friends for a user
	GetFriends(ctx context.Context, userID string, page, limit int) (*models.FriendsResponse, error)

	// SendFriendRequest sends a friend request
	SendFriendRequest(ctx context.Context, userID, friendID string) (*models.FriendRequestDetails, error)

	// GetFriendRequests retrieves friend requests for a user
	GetFriendRequests(ctx context.Context, userID, status string, page, limit int) (*models.FriendRequestsResponse, error)

	// AcceptFriendRequest accepts a friend request
	AcceptFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequestDetails, error)

	// RejectFriendRequest rejects a friend request
	RejectFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequestDetails, error)

	// RemoveFriend removes a friend
	RemoveFriend(ctx context.Context, userID, friendID string) (bool, error)

	// BlockUser blocks a user
	BlockUser(ctx context.Context, userID, blockedUserID string) (bool, error)

	// UnblockUser unblocks a user
	UnblockUser(ctx context.Context, userID, blockedUserID string) (bool, error)
}

// friendService implements the FriendService interface
type friendService struct {
	cfg    *config.Config
	logger *logger.Logger
	client pb.FriendServiceClient
}

// createAuthContext creates a new context with the JWT token in the metadata
func (s *friendService) createAuthContext(ctx context.Context) (context.Context, error) {
	// Extract token from context
	token, ok := ctx.Value("jwt_token").(string)
	if !ok || token == "" {
		return ctx, errors.New("no token found in context")
	}

	// Create metadata with authorization header
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create a new context with the metadata
	return metadata.NewOutgoingContext(ctx, md), nil
}

// NewFriendService creates a new friend service
func NewFriendService(cfg *config.Config, logger *logger.Logger) FriendService {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.FriendsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to friends service", err)
	}

	// Create a client
	client := pb.NewFriendServiceClient(conn)

	return &friendService{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// GetFriends retrieves friends for a user
func (s *friendService) GetFriends(ctx context.Context, userID string, page, limit int) (*models.FriendsResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service
	resp, err := s.client.GetFriends(authCtx, &pb.GetFriendsRequest{
		UserId: userID,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get friends", err)
		return nil, err
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

	return &models.FriendsResponse{
		Friends:    friends,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// SendFriendRequest sends a friend request
func (s *friendService) SendFriendRequest(ctx context.Context, userID, friendID string) (*models.FriendRequestDetails, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service
	resp, err := s.client.SendFriendRequest(authCtx, &pb.SendFriendRequestRequest{
		UserId:   userID,
		FriendId: friendID,
	})

	if err != nil {
		s.logger.Error("Failed to send friend request", err)
		return nil, err
	}

	return &models.FriendRequestDetails{
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
	}, nil
}

// GetFriendRequests retrieves friend requests for a user
func (s *friendService) GetFriendRequests(ctx context.Context, userID, status string, page, limit int) (*models.FriendRequestsResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service
	resp, err := s.client.GetFriendRequests(authCtx, &pb.GetFriendRequestsRequest{
		UserId: userID,
		Status: status,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get friend requests", err)
		return nil, err
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

	return &models.FriendRequestsResponse{
		Requests:   requests,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// AcceptFriendRequest accepts a friend request
func (s *friendService) AcceptFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequestDetails, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service
	resp, err := s.client.AcceptFriendRequest(authCtx, &pb.AcceptFriendRequestRequest{
		RequestId: requestID,
		UserId:    userID,
	})

	if err != nil {
		s.logger.Error("Failed to accept friend request", err)
		return nil, err
	}

	return &models.FriendRequestDetails{
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
	}, nil
}

// RejectFriendRequest rejects a friend request
func (s *friendService) RejectFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequestDetails, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service
	resp, err := s.client.RejectFriendRequest(authCtx, &pb.RejectFriendRequestRequest{
		RequestId: requestID,
		UserId:    userID,
	})

	if err != nil {
		s.logger.Error("Failed to reject friend request", err)
		return nil, err
	}

	return &models.FriendRequestDetails{
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
	}, nil
}

// RemoveFriend removes a friend
func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID string) (bool, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return false, err
	}

	// Call the gRPC service
	resp, err := s.client.RemoveFriend(authCtx, &pb.RemoveFriendRequest{
		UserId:   userID,
		FriendId: friendID,
	})

	if err != nil {
		s.logger.Error("Failed to remove friend", err)
		return false, err
	}

	return resp.Success, nil
}

// BlockUser blocks a user
func (s *friendService) BlockUser(ctx context.Context, userID, blockedUserID string) (bool, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return false, err
	}

	// Call the gRPC service
	resp, err := s.client.BlockUser(authCtx, &pb.BlockUserRequest{
		UserId:        userID,
		BlockedUserId: blockedUserID,
	})

	if err != nil {
		s.logger.Error("Failed to block user", err)
		return false, err
	}

	return resp.Success, nil
}

// UnblockUser unblocks a user
func (s *friendService) UnblockUser(ctx context.Context, userID, blockedUserID string) (bool, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return false, err
	}

	// Call the gRPC service
	resp, err := s.client.UnblockUser(authCtx, &pb.UnblockUserRequest{
		UserId:        userID,
		BlockedUserId: blockedUserID,
	})

	if err != nil {
		s.logger.Error("Failed to unblock user", err)
		return false, err
	}

	return resp.Success, nil
}
