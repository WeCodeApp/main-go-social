package models

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a post in the system
type Post struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	AuthorID      string         `gorm:"type:varchar(36);not null;index" json:"author_id"`
	AuthorName    string         `gorm:"type:varchar(255);not null" json:"author_name"`
	AuthorAvatar  string         `gorm:"type:varchar(255)" json:"author_avatar"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Visibility    string         `gorm:"type:enum('public','private');not null;default:'public'" json:"visibility"`
	GroupID       string         `gorm:"type:varchar(36);index" json:"group_id"`
	GroupName     string         `gorm:"type:varchar(255)" json:"group_name"`
	Media         string         `gorm:"type:text" json:"-"` // Stored as JSON array in database
	MediaArray    []string       `gorm:"-" json:"media"`     // Used in application
	LikesCount    int            `gorm:"default:0" json:"likes_count"`
	CommentsCount int            `gorm:"default:0" json:"comments_count"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Post model
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate is a hook that is called before creating a post
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}

// Comment represents a comment on a post
type Comment struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID       string         `gorm:"type:varchar(36);not null;index" json:"post_id"`
	AuthorID     string         `gorm:"type:varchar(36);not null;index" json:"author_id"`
	AuthorName   string         `gorm:"type:varchar(255);not null" json:"author_name"`
	AuthorAvatar string         `gorm:"type:varchar(255)" json:"author_avatar"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate is a hook that is called before creating a comment
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}

// Like represents a like on a post
type Like struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID    string         `gorm:"type:varchar(36);not null;index" json:"post_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Like model
func (Like) TableName() string {
	return "likes"
}

// BeforeCreate is a hook that is called before creating a like
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = generateUUID()
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
