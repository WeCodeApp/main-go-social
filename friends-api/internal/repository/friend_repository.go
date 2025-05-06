package repository

import (
	"friends-api/internal/models"

	"gorm.io/gorm"
)

// FriendRepository is the interface for friend-related operations
type FriendRepository interface {
	// Friend requests
	CreateFriendRequest(request *models.FriendRequest) error
	GetFriendRequestByID(id string) (*models.FriendRequest, error)
	GetFriendRequestsBySenderID(senderID string, status string, page, limit int) ([]*models.FriendRequest, int64, error)
	GetFriendRequestsByReceiverID(receiverID string, status string, page, limit int) ([]*models.FriendRequest, int64, error)
	UpdateFriendRequestStatus(id string, status string) error
	DeleteFriendRequest(id string) error

	// Friendships
	CreateFriendship(friendship *models.Friendship) error
	GetFriendshipByID(id string) (*models.Friendship, error)
	GetFriendshipsByUserID(userID string, page, limit int) ([]*models.Friendship, int64, error)
	DeleteFriendship(userID, friendID string) error

	// Blocked users
	BlockUser(blockedUser *models.BlockedUser) error
	UnblockUser(userID, blockedUserID string) error
	GetBlockedUsers(userID string, page, limit int) ([]*models.BlockedUser, int64, error)
	IsUserBlocked(userID, blockedUserID string) (bool, error)

	// Check friendship status
	CheckFriendship(userID, friendID string) (string, string, error)
}

// friendRepository is the implementation of FriendRepository
type friendRepository struct {
	db *gorm.DB
}

// NewFriendRepository creates a new friend repository
func NewFriendRepository(db *gorm.DB) FriendRepository {
	return &friendRepository{db: db}
}

// CreateFriendRequest creates a new friend request
func (r *friendRepository) CreateFriendRequest(request *models.FriendRequest) error {
	return r.db.Create(request).Error
}

// GetFriendRequestByID gets a friend request by ID
func (r *friendRepository) GetFriendRequestByID(id string) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Where("id = ?", id).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// GetFriendRequestsBySenderID gets friend requests by sender ID
func (r *friendRepository) GetFriendRequestsBySenderID(senderID string, status string, page, limit int) ([]*models.FriendRequest, int64, error) {
	var requests []*models.FriendRequest
	var count int64

	query := r.db.Model(&models.FriendRequest{}).Where("sender_id = ?", senderID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, count, nil
}

// GetFriendRequestsByReceiverID gets friend requests by receiver ID
func (r *friendRepository) GetFriendRequestsByReceiverID(receiverID string, status string, page, limit int) ([]*models.FriendRequest, int64, error) {
	var requests []*models.FriendRequest
	var count int64

	query := r.db.Model(&models.FriendRequest{}).Where("receiver_id = ?", receiverID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, count, nil
}

// UpdateFriendRequestStatus updates the status of a friend request
func (r *friendRepository) UpdateFriendRequestStatus(id string, status string) error {
	return r.db.Model(&models.FriendRequest{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteFriendRequest deletes a friend request
func (r *friendRepository) DeleteFriendRequest(id string) error {
	return r.db.Delete(&models.FriendRequest{}, "id = ?", id).Error
}

// CreateFriendship creates a new friendship
func (r *friendRepository) CreateFriendship(friendship *models.Friendship) error {
	return r.db.Create(friendship).Error
}

// GetFriendshipByID gets a friendship by ID
func (r *friendRepository) GetFriendshipByID(id string) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.Where("id = ?", id).First(&friendship).Error
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}

// GetFriendshipsByUserID gets friendships by user ID
func (r *friendRepository) GetFriendshipsByUserID(userID string, page, limit int) ([]*models.Friendship, int64, error) {
	var friendships []*models.Friendship
	var count int64

	query := r.db.Model(&models.Friendship{}).Where("user_id = ?", userID)

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&friendships).Error
	if err != nil {
		return nil, 0, err
	}

	return friendships, count, nil
}

// DeleteFriendship deletes a friendship
func (r *friendRepository) DeleteFriendship(userID, friendID string) error {
	// Delete both directions of the friendship
	err := r.db.Delete(&models.Friendship{}, "user_id = ? AND friend_id = ?", userID, friendID).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&models.Friendship{}, "user_id = ? AND friend_id = ?", friendID, userID).Error
}

// BlockUser blocks a user
func (r *friendRepository) BlockUser(blockedUser *models.BlockedUser) error {
	return r.db.Create(blockedUser).Error
}

// UnblockUser unblocks a user
func (r *friendRepository) UnblockUser(userID, blockedUserID string) error {
	return r.db.Delete(&models.BlockedUser{}, "user_id = ? AND blocked_user_id = ?", userID, blockedUserID).Error
}

// GetBlockedUsers gets blocked users
func (r *friendRepository) GetBlockedUsers(userID string, page, limit int) ([]*models.BlockedUser, int64, error) {
	var blockedUsers []*models.BlockedUser
	var count int64

	query := r.db.Model(&models.BlockedUser{}).Where("user_id = ?", userID)

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&blockedUsers).Error
	if err != nil {
		return nil, 0, err
	}

	return blockedUsers, count, nil
}

// IsUserBlocked checks if a user is blocked
func (r *friendRepository) IsUserBlocked(userID, blockedUserID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlockedUser{}).Where("user_id = ? AND blocked_user_id = ?", userID, blockedUserID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckFriendship checks the friendship status between two users
func (r *friendRepository) CheckFriendship(userID, friendID string) (string, string, error) {
	// Check if they are friends
	var friendshipCount int64
	err := r.db.Model(&models.Friendship{}).Where("user_id = ? AND friend_id = ?", userID, friendID).Count(&friendshipCount).Error
	if err != nil {
		return "", "", err
	}
	if friendshipCount > 0 {
		return "friends", "", nil
	}

	// Check if there's a pending friend request
	var request models.FriendRequest
	err = r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, friendID, friendID, userID).
		Where("status = ?", "pending").
		First(&request).Error
	if err == nil {
		return "pending", request.ID, nil
	} else if err != gorm.ErrRecordNotFound {
		return "", "", err
	}

	// Check if one has blocked the other
	var blockedCount int64
	err = r.db.Model(&models.BlockedUser{}).
		Where("(user_id = ? AND blocked_user_id = ?) OR (user_id = ? AND blocked_user_id = ?)", userID, friendID, friendID, userID).
		Count(&blockedCount).Error
	if err != nil {
		return "", "", err
	}
	if blockedCount > 0 {
		return "blocked", "", nil
	}

	// No relationship
	return "none", "", nil
}