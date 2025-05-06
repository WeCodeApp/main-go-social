package controllers

import (
	pb "common/pb/common/proto/friends"
	"context"
	"friends-api/internal/services"
	"friends-api/internal/utils/errors"
	"friends-api/internal/utils/logger"
)

// FriendController handles gRPC requests for friend-related operations
type FriendController struct {
	pb.UnimplementedFriendServiceServer
	service services.FriendService
	logger  *logger.Logger
}

// NewFriendController creates a new friend controller
func NewFriendController(service services.FriendService, logger *logger.Logger) *FriendController {
	return &FriendController{
		service: service,
		logger:  logger,
	}
}

// SendFriendRequest sends a friend request
func (c *FriendController) SendFriendRequest(ctx context.Context, req *pb.SendFriendRequestRequest) (*pb.FriendRequestResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Send friend request
	request, err := c.service.SendFriendRequest(ctx, userID, req.FriendId)
	if err != nil {
		c.logger.Error("Failed to send friend request", err)
		return nil, err
	}

	// Create response
	return &pb.FriendRequestResponse{
		RequestId:      request.ID,
		SenderId:       request.SenderID,
		SenderName:     "", // Would need to fetch from users service
		SenderAvatar:   "", // Would need to fetch from users service
		ReceiverId:     request.ReceiverID,
		ReceiverName:   "", // Would need to fetch from users service
		ReceiverAvatar: "", // Would need to fetch from users service
		Status:         request.Status,
		CreatedAt:      request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      request.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetFriendRequests retrieves friend requests for a user
func (c *FriendController) GetFriendRequests(ctx context.Context, req *pb.GetFriendRequestsRequest) (*pb.GetFriendRequestsResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Get friend requests
	requests, totalCount, totalPages, err := c.service.GetFriendRequests(ctx, userID, req.Status, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get friend requests", err)
		return nil, err
	}

	// Create response
	response := &pb.GetFriendRequestsResponse{
		Requests:   make([]*pb.FriendRequestResponse, 0, len(requests)),
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}

	// Add requests to response
	for _, request := range requests {
		response.Requests = append(response.Requests, &pb.FriendRequestResponse{
			RequestId:      request.ID,
			SenderId:       request.SenderID,
			SenderName:     "", // Would need to fetch from users service
			SenderAvatar:   "", // Would need to fetch from users service
			ReceiverId:     request.ReceiverID,
			ReceiverName:   "", // Would need to fetch from users service
			ReceiverAvatar: "", // Would need to fetch from users service
			Status:         request.Status,
			CreatedAt:      request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      request.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// AcceptFriendRequest accepts a friend request
func (c *FriendController) AcceptFriendRequest(ctx context.Context, req *pb.AcceptFriendRequestRequest) (*pb.FriendRequestResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Accept friend request
	request, err := c.service.AcceptFriendRequest(ctx, req.RequestId, userID)
	if err != nil {
		c.logger.Error("Failed to accept friend request", err)
		return nil, err
	}

	// Create response
	return &pb.FriendRequestResponse{
		RequestId:      request.ID,
		SenderId:       request.SenderID,
		SenderName:     "", // Would need to fetch from users service
		SenderAvatar:   "", // Would need to fetch from users service
		ReceiverId:     request.ReceiverID,
		ReceiverName:   "", // Would need to fetch from users service
		ReceiverAvatar: "", // Would need to fetch from users service
		Status:         request.Status,
		CreatedAt:      request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      request.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// RejectFriendRequest rejects a friend request
func (c *FriendController) RejectFriendRequest(ctx context.Context, req *pb.RejectFriendRequestRequest) (*pb.FriendRequestResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Reject friend request
	request, err := c.service.RejectFriendRequest(ctx, req.RequestId, userID)
	if err != nil {
		c.logger.Error("Failed to reject friend request", err)
		return nil, err
	}

	// Create response
	return &pb.FriendRequestResponse{
		RequestId:      request.ID,
		SenderId:       request.SenderID,
		SenderName:     "", // Would need to fetch from users service
		SenderAvatar:   "", // Would need to fetch from users service
		ReceiverId:     request.ReceiverID,
		ReceiverName:   "", // Would need to fetch from users service
		ReceiverAvatar: "", // Would need to fetch from users service
		Status:         request.Status,
		CreatedAt:      request.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      request.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetFriends retrieves friends for a user
func (c *FriendController) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	// Get user ID from context or request
	userID := req.UserId
	if userID == "" {
		var ok bool
		userID, ok = ctx.Value("userID").(string)
		if !ok {
			c.logger.Error("Failed to get user ID from context", nil)
			return nil, errors.ErrUnauthenticated
		}
	}

	// Get friends
	friendships, totalCount, totalPages, err := c.service.GetFriends(ctx, userID, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get friends", err)
		return nil, err
	}

	// Create response
	response := &pb.GetFriendsResponse{
		Friends:    make([]*pb.FriendResponse, 0, len(friendships)),
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}

	// Add friends to response
	for _, friendship := range friendships {
		response.Friends = append(response.Friends, &pb.FriendResponse{
			UserId:       friendship.FriendID,
			Name:         "", // Would need to fetch from users service
			Avatar:       "", // Would need to fetch from users service
			Email:        "", // Would need to fetch from users service
			FriendsSince: friendship.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// RemoveFriend removes a friend
func (c *FriendController) RemoveFriend(ctx context.Context, req *pb.RemoveFriendRequest) (*pb.RemoveFriendResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Remove friend
	err := c.service.RemoveFriend(ctx, userID, req.FriendId)
	if err != nil {
		c.logger.Error("Failed to remove friend", err)
		return nil, err
	}

	// Create response
	return &pb.RemoveFriendResponse{
		Success: true,
	}, nil
}

// BlockUser blocks a user
func (c *FriendController) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Block user
	err := c.service.BlockUser(ctx, userID, req.BlockedUserId)
	if err != nil {
		c.logger.Error("Failed to block user", err)
		return nil, err
	}

	// Create response
	return &pb.BlockUserResponse{
		Success: true,
	}, nil
}

// UnblockUser unblocks a user
func (c *FriendController) UnblockUser(ctx context.Context, req *pb.UnblockUserRequest) (*pb.UnblockUserResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Unblock user
	err := c.service.UnblockUser(ctx, userID, req.BlockedUserId)
	if err != nil {
		c.logger.Error("Failed to unblock user", err)
		return nil, err
	}

	// Create response
	return &pb.UnblockUserResponse{
		Success: true,
	}, nil
}

// GetBlockedUsers retrieves blocked users for a user
func (c *FriendController) GetBlockedUsers(ctx context.Context, req *pb.GetBlockedUsersRequest) (*pb.GetBlockedUsersResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Get blocked users
	blockedUsers, totalCount, totalPages, err := c.service.GetBlockedUsers(ctx, userID, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get blocked users", err)
		return nil, err
	}

	// Create response
	response := &pb.GetBlockedUsersResponse{
		BlockedUsers: make([]*pb.BlockedUserResponse, 0, len(blockedUsers)),
		TotalCount:   int32(totalCount),
		Page:         req.Page,
		TotalPages:   totalPages,
	}

	// Add blocked users to response
	for _, blockedUser := range blockedUsers {
		response.BlockedUsers = append(response.BlockedUsers, &pb.BlockedUserResponse{
			UserId:    blockedUser.BlockedUserID,
			Name:      "", // Would need to fetch from users service
			Avatar:    "", // Would need to fetch from users service
			BlockedAt: blockedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// CheckFriendship checks if two users are friends
func (c *FriendController) CheckFriendship(ctx context.Context, req *pb.CheckFriendshipRequest) (*pb.CheckFriendshipResponse, error) {
	// Check friendship
	status, requestID, err := c.service.CheckFriendship(ctx, req.UserId, req.FriendId)
	if err != nil {
		c.logger.Error("Failed to check friendship", err)
		return nil, err
	}

	// Create response
	return &pb.CheckFriendshipResponse{
		AreFriends: status == "friends",
		Status:     status,
		RequestId:  requestID,
	}, nil
}
