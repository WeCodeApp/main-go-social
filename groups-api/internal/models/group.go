package models

import (
	"time"

	"gorm.io/gorm"
)

// Group represents a group in the system
type Group struct {
	ID          string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	CreatorID   string         `gorm:"type:varchar(36);not null;index" json:"creator_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Group model
func (Group) TableName() string {
	return "groups"
}

// BeforeCreate is a hook that is called before creating a group
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = generateUUID()
	}
	return nil
}

// GroupMember represents a member of a group
type GroupMember struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	GroupID   string         `gorm:"type:varchar(36);not null;index" json:"group_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Role      string         `gorm:"type:enum('creator','admin','member');default:'member';not null" json:"role"`
	JoinedAt  time.Time      `json:"joined_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Group     *Group         `gorm:"foreignKey:GroupID" json:"-"`
}

// TableName returns the table name for the GroupMember model
func (GroupMember) TableName() string {
	return "group_members"
}

// BeforeCreate is a hook that is called before creating a group member
func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == "" {
		gm.ID = generateUUID()
	}
	if gm.JoinedAt.IsZero() {
		gm.JoinedAt = time.Now()
	}
	return nil
}

// GroupPost represents a post in a group
type GroupPost struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	GroupID   string         `gorm:"type:varchar(36);not null;index" json:"group_id"`
	AuthorID  string         `gorm:"type:varchar(36);not null;index" json:"author_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Group     *Group         `gorm:"foreignKey:GroupID" json:"-"`
	Media     []*GroupPostMedia `gorm:"foreignKey:PostID" json:"media,omitempty"`
	Likes     []*GroupPostLike  `gorm:"foreignKey:PostID" json:"likes,omitempty"`
	Comments  []*GroupPostComment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// TableName returns the table name for the GroupPost model
func (GroupPost) TableName() string {
	return "group_posts"
}

// BeforeCreate is a hook that is called before creating a group post
func (gp *GroupPost) BeforeCreate(tx *gorm.DB) error {
	if gp.ID == "" {
		gp.ID = generateUUID()
	}
	return nil
}

// GroupPostMedia represents media attached to a group post
type GroupPostMedia struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID    string         `gorm:"type:varchar(36);not null;index" json:"post_id"`
	MediaURL  string         `gorm:"type:varchar(255);not null" json:"media_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Post      *GroupPost     `gorm:"foreignKey:PostID" json:"-"`
}

// TableName returns the table name for the GroupPostMedia model
func (GroupPostMedia) TableName() string {
	return "group_post_media"
}

// BeforeCreate is a hook that is called before creating group post media
func (gpm *GroupPostMedia) BeforeCreate(tx *gorm.DB) error {
	if gpm.ID == "" {
		gpm.ID = generateUUID()
	}
	return nil
}

// GroupPostLike represents a like on a group post
type GroupPostLike struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID    string         `gorm:"type:varchar(36);not null;index" json:"post_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Post      *GroupPost     `gorm:"foreignKey:PostID" json:"-"`
}

// TableName returns the table name for the GroupPostLike model
func (GroupPostLike) TableName() string {
	return "group_post_likes"
}

// BeforeCreate is a hook that is called before creating a group post like
func (gpl *GroupPostLike) BeforeCreate(tx *gorm.DB) error {
	if gpl.ID == "" {
		gpl.ID = generateUUID()
	}
	return nil
}

// GroupPostComment represents a comment on a group post
type GroupPostComment struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID    string         `gorm:"type:varchar(36);not null;index" json:"post_id"`
	UserID    string         `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Post      *GroupPost     `gorm:"foreignKey:PostID" json:"-"`
}

// TableName returns the table name for the GroupPostComment model
func (GroupPostComment) TableName() string {
	return "group_post_comments"
}

// BeforeCreate is a hook that is called before creating a group post comment
func (gpc *GroupPostComment) BeforeCreate(tx *gorm.DB) error {
	if gpc.ID == "" {
		gpc.ID = generateUUID()
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
