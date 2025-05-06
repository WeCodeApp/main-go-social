package repository

import (
	"context"
	"encoding/json"
	"post-api/internal/models"

	"gorm.io/gorm"
)

// PostRepository defines the interface for post repository operations
type PostRepository interface {
	// Create creates a new post
	Create(ctx context.Context, post *models.Post) error

	// FindByID finds a post by ID
	FindByID(ctx context.Context, id string) (*models.Post, error)

	// FindByAuthor finds posts by author ID with pagination
	FindByAuthor(ctx context.Context, authorID string, page, limit int) ([]*models.Post, int64, error)

	// FindByGroup finds posts by group ID with pagination
	FindByGroup(ctx context.Context, groupID string, page, limit int) ([]*models.Post, int64, error)

	// FindPublic finds public posts with pagination
	FindPublic(ctx context.Context, page, limit int) ([]*models.Post, int64, error)

	// FindVisible finds posts visible to a user (public or authored by friends) with pagination
	FindVisible(ctx context.Context, userID string, friendIDs []string, page, limit int) ([]*models.Post, int64, error)

	// Update updates a post
	Update(ctx context.Context, post *models.Post) error

	// Delete deletes a post
	Delete(ctx context.Context, id string) error

	// IncrementLikesCount increments the likes count for a post
	IncrementLikesCount(ctx context.Context, id string) error

	// DecrementLikesCount decrements the likes count for a post
	DecrementLikesCount(ctx context.Context, id string) error

	// IncrementCommentsCount increments the comments count for a post
	IncrementCommentsCount(ctx context.Context, id string) error

	// DecrementCommentsCount decrements the comments count for a post
	DecrementCommentsCount(ctx context.Context, id string) error
}

// postRepository implements the PostRepository interface
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post
func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	// Convert media array to JSON string if it's not empty
	if len(post.MediaArray) > 0 {
		mediaJSON, err := json.Marshal(post.MediaArray)
		if err != nil {
			return err
		}
		post.Media = string(mediaJSON)
	}

	return r.db.WithContext(ctx).Create(post).Error
}

// FindByID finds a post by ID
func (r *postRepository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}

	// Parse media JSON string to array if it's not empty
	if post.Media != "" {
		var mediaArray []string
		if err := json.Unmarshal([]byte(post.Media), &mediaArray); err != nil {
			return nil, err
		}
		post.MediaArray = mediaArray
	}

	return &post, nil
}

// FindByAuthor finds posts by author ID with pagination
func (r *postRepository) FindByAuthor(ctx context.Context, authorID string, page, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var count int64

	offset := (page - 1) * limit

	// Count total posts by author
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Where("author_id = ?", authorID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get posts by author with pagination
	if err := r.db.WithContext(ctx).Where("author_id = ?", authorID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Parse media JSON string to array for each post
	for _, post := range posts {
		if post.Media != "" {
			var mediaArray []string
			if err := json.Unmarshal([]byte(post.Media), &mediaArray); err != nil {
				return nil, 0, err
			}
			post.MediaArray = mediaArray
		}
	}

	return posts, count, nil
}

// FindByGroup finds posts by group ID with pagination
func (r *postRepository) FindByGroup(ctx context.Context, groupID string, page, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var count int64

	offset := (page - 1) * limit

	// Count total posts by group
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Where("group_id = ?", groupID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get posts by group with pagination
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Parse media JSON string to array for each post
	for _, post := range posts {
		if post.Media != "" {
			var mediaArray []string
			if err := json.Unmarshal([]byte(post.Media), &mediaArray); err != nil {
				return nil, 0, err
			}
			post.MediaArray = mediaArray
		}
	}

	return posts, count, nil
}

// FindPublic finds public posts with pagination
func (r *postRepository) FindPublic(ctx context.Context, page, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var count int64

	offset := (page - 1) * limit

	// Count total public posts
	if err := r.db.WithContext(ctx).Model(&models.Post{}).Where("visibility = ?", "public").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get public posts with pagination
	if err := r.db.WithContext(ctx).Where("visibility = ?", "public").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Parse media JSON string to array for each post
	for _, post := range posts {
		if post.Media != "" {
			var mediaArray []string
			if err := json.Unmarshal([]byte(post.Media), &mediaArray); err != nil {
				return nil, 0, err
			}
			post.MediaArray = mediaArray
		}
	}

	return posts, count, nil
}

// FindVisible finds posts visible to a user (public or authored by friends) with pagination
func (r *postRepository) FindVisible(ctx context.Context, userID string, friendIDs []string, page, limit int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var count int64

	offset := (page - 1) * limit

	// Count total visible posts
	query := r.db.WithContext(ctx).Model(&models.Post{}).Where("visibility = ? OR (visibility = ? AND author_id IN ?)", "public", "private", friendIDs)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get visible posts with pagination
	if err := r.db.WithContext(ctx).Where("visibility = ? OR (visibility = ? AND author_id IN ?)", "public", "private", friendIDs).Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// Parse media JSON string to array for each post
	for _, post := range posts {
		if post.Media != "" {
			var mediaArray []string
			if err := json.Unmarshal([]byte(post.Media), &mediaArray); err != nil {
				return nil, 0, err
			}
			post.MediaArray = mediaArray
		}
	}

	return posts, count, nil
}

// Update updates a post
func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	// Convert media array to JSON string if it's not empty
	if len(post.MediaArray) > 0 {
		mediaJSON, err := json.Marshal(post.MediaArray)
		if err != nil {
			return err
		}
		post.Media = string(mediaJSON)
	}

	return r.db.WithContext(ctx).Save(post).Error
}

// Delete deletes a post
func (r *postRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Post{}, "id = ?", id).Error
}

// IncrementLikesCount increments the likes count for a post
func (r *postRepository) IncrementLikesCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Update("likes_count", gorm.Expr("likes_count + ?", 1)).Error
}

// DecrementLikesCount decrements the likes count for a post
func (r *postRepository) DecrementLikesCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Update("likes_count", gorm.Expr("likes_count - ?", 1)).Error
}

// IncrementCommentsCount increments the comments count for a post
func (r *postRepository) IncrementCommentsCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Update("comments_count", gorm.Expr("comments_count + ?", 1)).Error
}

// DecrementCommentsCount decrements the comments count for a post
func (r *postRepository) DecrementCommentsCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Update("comments_count", gorm.Expr("comments_count - ?", 1)).Error
}
