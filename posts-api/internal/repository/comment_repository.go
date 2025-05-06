package repository

import (
	"context"
	"post-api/internal/models"

	"gorm.io/gorm"
)

// CommentRepository defines the interface for comment repository operations
type CommentRepository interface {
	// Create creates a new comment
	Create(ctx context.Context, comment *models.Comment) error

	// FindByID finds a comment by ID
	FindByID(ctx context.Context, id string) (*models.Comment, error)

	// FindByPost finds comments for a post with pagination
	FindByPost(ctx context.Context, postID string, page, limit int) ([]*models.Comment, int64, error)

	// Update updates a comment
	Update(ctx context.Context, comment *models.Comment) error

	// Delete deletes a comment
	Delete(ctx context.Context, id string) error
}

// commentRepository implements the CommentRepository interface
type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

// Create creates a new comment
func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// FindByID finds a comment by ID
func (r *commentRepository) FindByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByPost finds comments for a post with pagination
func (r *commentRepository) FindByPost(ctx context.Context, postID string, page, limit int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var count int64

	offset := (page - 1) * limit

	// Count total comments for the post
	if err := r.db.WithContext(ctx).Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get comments for the post with pagination
	if err := r.db.WithContext(ctx).Where("post_id = ?", postID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, count, nil
}

// Update updates a comment
func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete deletes a comment
func (r *commentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, "id = ?", id).Error
}