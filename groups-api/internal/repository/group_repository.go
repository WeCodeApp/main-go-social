package repository

import (
	"context"
	"groups-api/internal/models"

	"gorm.io/gorm"
)

// GroupRepository defines the interface for group-related database operations
type GroupRepository interface {
	// Group operations
	CreateGroup(ctx context.Context, group *models.Group) error
	GetGroupByID(ctx context.Context, id string) (*models.Group, error)
	GetGroups(ctx context.Context, query string, page, limit int) ([]*models.Group, int64, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, id string) error

	// Group member operations
	AddMember(ctx context.Context, member *models.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	GetGroupMembers(ctx context.Context, groupID string, page, limit int) ([]*models.GroupMember, int64, error)
	GetMemberByID(ctx context.Context, groupID, userID string) (*models.GroupMember, error)
	UpdateMember(ctx context.Context, member *models.GroupMember) error
	IsMember(ctx context.Context, groupID, userID string) (bool, error)

	// Group post operations
	CreatePost(ctx context.Context, post *models.GroupPost) error
	GetPostByID(ctx context.Context, id string) (*models.GroupPost, error)
	GetGroupPosts(ctx context.Context, groupID string, page, limit int) ([]*models.GroupPost, int64, error)
	UpdatePost(ctx context.Context, post *models.GroupPost) error
	DeletePost(ctx context.Context, id string) error

	// Group post media operations
	AddPostMedia(ctx context.Context, media *models.GroupPostMedia) error
	GetPostMedia(ctx context.Context, postID string) ([]*models.GroupPostMedia, error)
	DeletePostMedia(ctx context.Context, id string) error

	// Group post like operations
	LikePost(ctx context.Context, like *models.GroupPostLike) error
	UnlikePost(ctx context.Context, postID, userID string) error
	GetPostLikes(ctx context.Context, postID string) ([]*models.GroupPostLike, error)
	IsPostLiked(ctx context.Context, postID, userID string) (bool, error)

	// Group post comment operations
	CreateComment(ctx context.Context, comment *models.GroupPostComment) error
	GetCommentByID(ctx context.Context, id string) (*models.GroupPostComment, error)
	GetPostComments(ctx context.Context, postID string, page, limit int) ([]*models.GroupPostComment, int64, error)
	UpdateComment(ctx context.Context, comment *models.GroupPostComment) error
	DeleteComment(ctx context.Context, id string) error
}

// groupRepository implements the GroupRepository interface
type groupRepository struct {
	db *gorm.DB
}

// NewGroupRepository creates a new group repository
func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

// CreateGroup creates a new group
func (r *groupRepository) CreateGroup(ctx context.Context, group *models.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// GetGroupByID gets a group by ID
func (r *groupRepository) GetGroupByID(ctx context.Context, id string) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetGroups gets groups with pagination and filtering
func (r *groupRepository) GetGroups(ctx context.Context, query string, page, limit int) ([]*models.Group, int64, error) {
	var groups []*models.Group
	var count int64

	db := r.db.WithContext(ctx)
	if query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	err := db.Model(&models.Group{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = db.Offset(offset).Limit(limit).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	return groups, count, nil
}

// UpdateGroup updates a group
func (r *groupRepository) UpdateGroup(ctx context.Context, group *models.Group) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// DeleteGroup deletes a group
func (r *groupRepository) DeleteGroup(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Group{}, "id = ?", id).Error
}

// AddMember adds a member to a group
func (r *groupRepository) AddMember(ctx context.Context, member *models.GroupMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

// RemoveMember removes a member from a group
func (r *groupRepository) RemoveMember(ctx context.Context, groupID, userID string) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error
}

// GetGroupMembers gets members of a group with pagination
func (r *groupRepository) GetGroupMembers(ctx context.Context, groupID string, page, limit int) ([]*models.GroupMember, int64, error) {
	var members []*models.GroupMember
	var count int64

	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).Where("group_id = ?", groupID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.WithContext(ctx).Where("group_id = ?", groupID).Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, count, nil
}

// GetMemberByID gets a member by group ID and user ID
func (r *groupRepository) GetMemberByID(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	var member models.GroupMember
	err := r.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateMember updates a group member
func (r *groupRepository) UpdateMember(ctx context.Context, member *models.GroupMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

// IsMember checks if a user is a member of a group
func (r *groupRepository) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreatePost creates a new post in a group
func (r *groupRepository) CreatePost(ctx context.Context, post *models.GroupPost) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// GetPostByID gets a post by ID
func (r *groupRepository) GetPostByID(ctx context.Context, id string) (*models.GroupPost, error) {
	var post models.GroupPost
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetGroupPosts gets posts in a group with pagination
func (r *groupRepository) GetGroupPosts(ctx context.Context, groupID string, page, limit int) ([]*models.GroupPost, int64, error) {
	var posts []*models.GroupPost
	var count int64

	err := r.db.WithContext(ctx).Model(&models.GroupPost{}).Where("group_id = ?", groupID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.WithContext(ctx).Where("group_id = ?", groupID).Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

// UpdatePost updates a post
func (r *groupRepository) UpdatePost(ctx context.Context, post *models.GroupPost) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// DeletePost deletes a post
func (r *groupRepository) DeletePost(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.GroupPost{}, "id = ?", id).Error
}

// AddPostMedia adds media to a post
func (r *groupRepository) AddPostMedia(ctx context.Context, media *models.GroupPostMedia) error {
	return r.db.WithContext(ctx).Create(media).Error
}

// GetPostMedia gets media for a post
func (r *groupRepository) GetPostMedia(ctx context.Context, postID string) ([]*models.GroupPostMedia, error) {
	var media []*models.GroupPostMedia
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&media).Error
	if err != nil {
		return nil, err
	}
	return media, nil
}

// DeletePostMedia deletes media from a post
func (r *groupRepository) DeletePostMedia(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.GroupPostMedia{}, "id = ?", id).Error
}

// LikePost likes a post
func (r *groupRepository) LikePost(ctx context.Context, like *models.GroupPostLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

// UnlikePost unlikes a post
func (r *groupRepository) UnlikePost(ctx context.Context, postID, userID string) error {
	return r.db.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.GroupPostLike{}).Error
}

// GetPostLikes gets likes for a post
func (r *groupRepository) GetPostLikes(ctx context.Context, postID string) ([]*models.GroupPostLike, error) {
	var likes []*models.GroupPostLike
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// IsPostLiked checks if a post is liked by a user
func (r *groupRepository) IsPostLiked(ctx context.Context, postID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupPostLike{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateComment creates a new comment on a post
func (r *groupRepository) CreateComment(ctx context.Context, comment *models.GroupPostComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetCommentByID gets a comment by ID
func (r *groupRepository) GetCommentByID(ctx context.Context, id string) (*models.GroupPostComment, error) {
	var comment models.GroupPostComment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetPostComments gets comments for a post with pagination
func (r *groupRepository) GetPostComments(ctx context.Context, postID string, page, limit int) ([]*models.GroupPostComment, int64, error) {
	var comments []*models.GroupPostComment
	var count int64

	err := r.db.WithContext(ctx).Model(&models.GroupPostComment{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.WithContext(ctx).Where("post_id = ?", postID).Offset(offset).Limit(limit).Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, count, nil
}

// UpdateComment updates a comment
func (r *groupRepository) UpdateComment(ctx context.Context, comment *models.GroupPostComment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// DeleteComment deletes a comment
func (r *groupRepository) DeleteComment(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.GroupPostComment{}, "id = ?", id).Error
}