package services

import (
	"context"
	"errors"
	"friends-api/internal/models"
	"friends-api/internal/repository"
	"friends-api/internal/utils/logger"
	"time"

	"gorm.io/gorm"
)

// FriendService is the interface for friend-related operations
type FriendService interface {
	// Friend requests
	SendFriendRequest(ctx context.Context, senderID, receiverID string) (*models.FriendRequest, error)
	GetFriendRequests(ctx context.Context, userID, status string, page, limit int) ([]*models.FriendRequest, int64, int32, error)
	AcceptFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequest, error)
	RejectFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequest, error)

	// Friendships
	GetFriends(ctx context.Context, userID string, page, limit int) ([]*models.Friendship, int64, int32, error)
	RemoveFriend(ctx context.Context, userID, friendID string) error

	// Blocked users
	BlockUser(ctx context.Context, userID, blockedUserID string) error
	UnblockUser(ctx context.Context, userID, blockedUserID string) error
	GetBlockedUsers(ctx context.Context, userID string, page, limit int) ([]*models.BlockedUser, int64, int32, error)

	// Check friendship status
	CheckFriendship(ctx context.Context, userID, friendID string) (string, string, error)
}

// friendService is the implementation of FriendService
type friendService struct {
	repo   repository.FriendRepository
	logger *logger.Logger
}

// NewFriendService creates a new friend service
func NewFriendService(repo repository.FriendRepository, logger *logger.Logger) FriendService {
	return &friendService{
		repo:   repo,
		logger: logger,
	}
}

// SendFriendRequest sends a friend request
func (s *friendService) SendFriendRequest(ctx context.Context, senderID, receiverID string) (*models.FriendRequest, error) {
	// Check if sender and receiver are the same
	if senderID == receiverID {
		return nil, errors.New("cannot send friend request to yourself")
	}

	// Check if they are already friends
	status, _, err := s.repo.CheckFriendship(senderID, receiverID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("Failed to check friendship", err)
		return nil, err
	}

	if status == "friends" {
		return nil, errors.New("already friends")
	}

	if status == "pending" {
		return nil, errors.New("friend request already sent")
	}

	if status == "blocked" {
		return nil, errors.New("cannot send friend request to blocked user")
	}

	// Create friend request
	request := &models.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
	}

	err = s.repo.CreateFriendRequest(request)
	if err != nil {
		s.logger.Error("Failed to create friend request", err)
		return nil, err
	}

	return request, nil
}

// GetFriendRequests gets friend requests for a user
func (s *friendService) GetFriendRequests(ctx context.Context, userID, status string, page, limit int) ([]*models.FriendRequest, int64, int32, error) {
	// Get friend requests
	requests, count, err := s.repo.GetFriendRequestsByReceiverID(userID, status, page, limit)
	if err != nil {
		s.logger.Error("Failed to get friend requests", err)
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return requests, count, totalPages, nil
}

// AcceptFriendRequest accepts a friend request
func (s *friendService) AcceptFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequest, error) {
	// Get friend request
	request, err := s.repo.GetFriendRequestByID(requestID)
	if err != nil {
		s.logger.Error("Failed to get friend request", err)
		return nil, err
	}

	// Check if the user is the receiver of the request
	if request.ReceiverID != userID {
		return nil, errors.New("not authorized to accept this friend request")
	}

	// Check if the request is pending
	if request.Status != "pending" {
		return nil, errors.New("friend request is not pending")
	}

	// Update request status
	err = s.repo.UpdateFriendRequestStatus(requestID, "accepted")
	if err != nil {
		s.logger.Error("Failed to update friend request status", err)
		return nil, err
	}

	// Create friendship (both ways)
	friendship1 := &models.Friendship{
		UserID:   request.SenderID,
		FriendID: request.ReceiverID,
	}
	err = s.repo.CreateFriendship(friendship1)
	if err != nil {
		s.logger.Error("Failed to create friendship", err)
		return nil, err
	}

	friendship2 := &models.Friendship{
		UserID:   request.ReceiverID,
		FriendID: request.SenderID,
	}
	err = s.repo.CreateFriendship(friendship2)
	if err != nil {
		s.logger.Error("Failed to create friendship", err)
		return nil, err
	}

	// Update request
	request.Status = "accepted"
	request.UpdatedAt = time.Now()

	return request, nil
}

// RejectFriendRequest rejects a friend request
func (s *friendService) RejectFriendRequest(ctx context.Context, requestID, userID string) (*models.FriendRequest, error) {
	// Get friend request
	request, err := s.repo.GetFriendRequestByID(requestID)
	if err != nil {
		s.logger.Error("Failed to get friend request", err)
		return nil, err
	}

	// Check if the user is the receiver of the request
	if request.ReceiverID != userID {
		return nil, errors.New("not authorized to reject this friend request")
	}

	// Check if the request is pending
	if request.Status != "pending" {
		return nil, errors.New("friend request is not pending")
	}

	// Update request status
	err = s.repo.UpdateFriendRequestStatus(requestID, "rejected")
	if err != nil {
		s.logger.Error("Failed to update friend request status", err)
		return nil, err
	}

	// Update request
	request.Status = "rejected"
	request.UpdatedAt = time.Now()

	return request, nil
}

// GetFriends gets friends for a user
func (s *friendService) GetFriends(ctx context.Context, userID string, page, limit int) ([]*models.Friendship, int64, int32, error) {
	// Get friendships
	friendships, count, err := s.repo.GetFriendshipsByUserID(userID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get friendships", err)
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return friendships, count, totalPages, nil
}

// RemoveFriend removes a friend
func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID string) error {
	// Check if they are friends
	status, _, err := s.repo.CheckFriendship(userID, friendID)
	if err != nil {
		s.logger.Error("Failed to check friendship", err)
		return err
	}

	if status != "friends" {
		return errors.New("not friends")
	}

	// Delete friendship
	err = s.repo.DeleteFriendship(userID, friendID)
	if err != nil {
		s.logger.Error("Failed to delete friendship", err)
		return err
	}

	return nil
}

// BlockUser blocks a user
func (s *friendService) BlockUser(ctx context.Context, userID, blockedUserID string) error {
	// Check if user is trying to block themselves
	if userID == blockedUserID {
		return errors.New("cannot block yourself")
	}

	// Check if already blocked
	isBlocked, err := s.repo.IsUserBlocked(userID, blockedUserID)
	if err != nil {
		s.logger.Error("Failed to check if user is blocked", err)
		return err
	}

	if isBlocked {
		return errors.New("user is already blocked")
	}

	// Remove friendship if they are friends
	status, _, err := s.repo.CheckFriendship(userID, blockedUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("Failed to check friendship", err)
		return err
	}

	if status == "friends" {
		err = s.repo.DeleteFriendship(userID, blockedUserID)
		if err != nil {
			s.logger.Error("Failed to delete friendship", err)
			return err
		}
	}

	// Block user
	blockedUser := &models.BlockedUser{
		UserID:        userID,
		BlockedUserID: blockedUserID,
	}

	err = s.repo.BlockUser(blockedUser)
	if err != nil {
		s.logger.Error("Failed to block user", err)
		return err
	}

	return nil
}

// UnblockUser unblocks a user
func (s *friendService) UnblockUser(ctx context.Context, userID, blockedUserID string) error {
	// Check if blocked
	isBlocked, err := s.repo.IsUserBlocked(userID, blockedUserID)
	if err != nil {
		s.logger.Error("Failed to check if user is blocked", err)
		return err
	}

	if !isBlocked {
		return errors.New("user is not blocked")
	}

	// Unblock user
	err = s.repo.UnblockUser(userID, blockedUserID)
	if err != nil {
		s.logger.Error("Failed to unblock user", err)
		return err
	}

	return nil
}

// GetBlockedUsers gets blocked users for a user
func (s *friendService) GetBlockedUsers(ctx context.Context, userID string, page, limit int) ([]*models.BlockedUser, int64, int32, error) {
	// Get blocked users
	blockedUsers, count, err := s.repo.GetBlockedUsers(userID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get blocked users", err)
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return blockedUsers, count, totalPages, nil
}

// CheckFriendship checks the friendship status between two users
func (s *friendService) CheckFriendship(ctx context.Context, userID, friendID string) (string, string, error) {
	// Check if user is trying to check friendship with themselves
	if userID == friendID {
		return "self", "", nil
	}

	// Check friendship status
	status, requestID, err := s.repo.CheckFriendship(userID, friendID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.logger.Error("Failed to check friendship", err)
		return "", "", err
	}

	return status, requestID, nil
}