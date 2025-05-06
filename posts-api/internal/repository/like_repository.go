package repository

import (
	"context"
	"post-api/internal/models"

	"gorm.io/gorm"
)

// LikeRepository defines the interface for like repository operations
type LikeRepository interface {
	// Create creates a new like
	Create(ctx context.Context, like *models.Like) error

	// FindByID finds a like by ID
	FindByID(ctx context.Context, id string) (*models.Like, error)

	// FindByPostAndUser finds a like by post ID and user ID
	FindByPostAndUser(ctx context.Context, postID, userID string) (*models.Like, error)

	// FindByPost finds likes for a post
	FindByPost(ctx context.Context, postID string) ([]*models.Like, error)

	// CountByPost counts likes for a post
	CountByPost(ctx context.Context, postID string) (int64, error)

	// Delete deletes a like
	Delete(ctx context.Context, id string) error

	// DeleteByPostAndUser deletes a like by post ID and user ID
	DeleteByPostAndUser(ctx context.Context, postID, userID string) error
}

// likeRepository implements the LikeRepository interface
type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository creates a new like repository
func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// Create creates a new like
func (r *likeRepository) Create(ctx context.Context, like *models.Like) error {
	return r.db.WithContext(ctx).Create(like).Error
}

// FindByID finds a like by ID
func (r *likeRepository) FindByID(ctx context.Context, id string) (*models.Like, error) {
	var like models.Like
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// FindByPostAndUser finds a like by post ID and user ID
func (r *likeRepository) FindByPostAndUser(ctx context.Context, postID, userID string) (*models.Like, error) {
	var like models.Like
	err := r.db.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// FindByPost finds likes for a post
func (r *likeRepository) FindByPost(ctx context.Context, postID string) ([]*models.Like, error) {
	var likes []*models.Like
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// CountByPost counts likes for a post
func (r *likeRepository) CountByPost(ctx context.Context, postID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete deletes a like
func (r *likeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Like{}, "id = ?", id).Error
}

// DeleteByPostAndUser deletes a like by post ID and user ID
func (r *likeRepository) DeleteByPostAndUser(ctx context.Context, postID, userID string) error {
	return r.db.WithContext(ctx).Delete(&models.Like{}, "post_id = ? AND user_id = ?", postID, userID).Error
}