package models

import (
	"time"

	"gorm.io/gorm"
)

// FriendRequest represents a friend request in the system
type FriendRequest struct {
	ID         string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	SenderID   string         `gorm:"type:varchar(36);not null;index" json:"sender_id"`
	ReceiverID string         `gorm:"type:varchar(36);not null;index" json:"receiver_id"`
	Status     string         `gorm:"type:enum('pending','accepted','rejected');default:'pending';not null" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the FriendRequest model
func (FriendRequest) TableName() string {
	return "friend_requests"
}

// BeforeCreate is a hook that is called before creating a friend request
func (fr *FriendRequest) BeforeCreate(tx *gorm.DB) error {
	if fr.ID == "" {
		fr.ID = generateUUID()
	}
	return nil
}

// Friendship represents a friendship between two users
type Friendship struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	FriendID  string         `gorm:"type:varchar(36);not null;index" json:"friend_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Friendship model
func (Friendship) TableName() string {
	return "friendships"
}

// BeforeCreate is a hook that is called before creating a friendship
func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = generateUUID()
	}
	return nil
}

// BlockedUser represents a blocked user in the system
type BlockedUser struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	BlockedUserID string        `gorm:"type:varchar(36);not null;index" json:"blocked_user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the BlockedUser model
func (BlockedUser) TableName() string {
	return "blocked_users"
}

// BeforeCreate is a hook that is called before creating a blocked user
func (bu *BlockedUser) BeforeCreate(tx *gorm.DB) error {
	if bu.ID == "" {
		bu.ID = generateUUID()
	}
	return nil
}

// generateUUID generates a UUID
func generateUUID() string {
	// This is a simple implementation for demonstration purposes
	// In a real application, you would use a proper UUID library
	return time.Now().Format("20060102150405") + randomString(10)
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure different values for each character
	}
	return string(result)
}