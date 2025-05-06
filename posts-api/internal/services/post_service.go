package services

import (
	"context"
	"post-api/internal/models"
	"post-api/internal/repository"
	"post-api/internal/utils/logger"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PostService defines the interface for post-related operations
type PostService interface {
	// CreatePost creates a new post
	CreatePost(ctx context.Context, userID, content, visibility, groupID string, media []string) (*models.Post, error)

	// GetPost retrieves a post by ID
	GetPost(ctx context.Context, postID, userID string) (*models.Post, bool, error)

	// GetPosts retrieves posts with pagination and filtering
	GetPosts(ctx context.Context, userID, authorID, groupID, visibility string, page, limit int, friendIDs []string) ([]*models.Post, int64, int32, error)

	// UpdatePost updates a post
	UpdatePost(ctx context.Context, postID, userID, content, visibility string, media []string) (*models.Post, error)

	// DeletePost deletes a post
	DeletePost(ctx context.Context, postID, userID string) error

	// AddComment adds a comment to a post
	AddComment(ctx context.Context, postID, userID, authorName, authorAvatar, content string) (*models.Comment, error)

	// GetComments retrieves comments for a post
	GetComments(ctx context.Context, postID string, page, limit int) ([]*models.Comment, int64, int32, error)

	// DeleteComment deletes a comment
	DeleteComment(ctx context.Context, commentID, postID, userID string) error

	// LikePost likes a post
	LikePost(ctx context.Context, postID, userID string) (int32, error)

	// UnlikePost unlikes a post
	UnlikePost(ctx context.Context, postID, userID string) (int32, error)

	// IsLiked checks if a post is liked by a user
	IsLiked(ctx context.Context, postID, userID string) (bool, error)
}

// postService implements the PostService interface
type postService struct {
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	likeRepo    repository.LikeRepository
	logger      *logger.Logger
}

// NewPostService creates a new post service
func NewPostService(
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	logger *logger.Logger,
) PostService {
	return &postService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		logger:      logger,
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(ctx context.Context, userID, content, visibility, groupID string, media []string) (*models.Post, error) {
	// Validate input
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	if visibility != "public" && visibility != "private" {
		return nil, status.Error(codes.InvalidArgument, "visibility must be 'public' or 'private'")
	}

	// TODO: Get user info from users service
	authorName := "User " + userID // Placeholder
	authorAvatar := ""             // Placeholder

	// TODO: If groupID is provided, validate that the user is a member of the group
	var groupName string
	if groupID != "" {
		// TODO: Get group info from groups service
		groupName = "Group " + groupID // Placeholder
	}

	// Create post
	post := &models.Post{
		AuthorID:      userID,
		AuthorName:    authorName,
		AuthorAvatar:  authorAvatar,
		Content:       content,
		Visibility:    visibility,
		GroupID:       groupID,
		GroupName:     groupName,
		MediaArray:    media,
		LikesCount:    0,
		CommentsCount: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save post to database
	if err := s.postRepo.Create(ctx, post); err != nil {
		s.logger.Error("Failed to create post", err)
		return nil, status.Error(codes.Internal, "failed to create post")
	}

	return post, nil
}

// GetPost retrieves a post by ID
func (s *postService) GetPost(ctx context.Context, postID, userID string) (*models.Post, bool, error) {
	// Validate input
	if postID == "" {
		return nil, false, status.Error(codes.InvalidArgument, "post ID is required")
	}

	// Get post from database
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return nil, false, status.Error(codes.NotFound, "post not found")
	}

	// Check if the post is visible to the user
	isVisible := s.isPostVisibleToUser(post, userID, nil)
	if !isVisible {
		return nil, false, status.Error(codes.PermissionDenied, "you don't have permission to view this post")
	}

	// Check if the post is liked by the user
	isLiked := false
	if userID != "" {
		isLiked, _ = s.IsLiked(ctx, postID, userID)
	}

	return post, isLiked, nil
}

// GetPosts retrieves posts with pagination and filtering
func (s *postService) GetPosts(ctx context.Context, userID, authorID, groupID, visibility string, page, limit int, friendIDs []string) ([]*models.Post, int64, int32, error) {
	// Validate input
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var posts []*models.Post
	var count int64
	var err error

	// Get posts based on filters
	if authorID != "" {
		// Get posts by author
		posts, count, err = s.postRepo.FindByAuthor(ctx, authorID, page, limit)
	} else if groupID != "" {
		// Get posts by group
		posts, count, err = s.postRepo.FindByGroup(ctx, groupID, page, limit)
	} else if userID == "" || visibility == "public" {
		// Get public posts
		posts, count, err = s.postRepo.FindPublic(ctx, page, limit)
	} else {
		// Get posts visible to the user
		posts, count, err = s.postRepo.FindVisible(ctx, userID, friendIDs, page, limit)
	}

	if err != nil {
		s.logger.Error("Failed to get posts", err)
		return nil, 0, 0, status.Error(codes.Internal, "failed to get posts")
	}

	// Filter posts based on visibility
	visiblePosts := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if s.isPostVisibleToUser(post, userID, friendIDs) {
			// Check if the post is liked by the user
			if userID != "" {
				isLiked, _ := s.IsLiked(ctx, post.ID, userID)
				if isLiked {
					// Set isLiked flag (this is not stored in the database)
					// We'll need to add this field to the Post struct or response
				}
			}
			visiblePosts = append(visiblePosts, post)
		}
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return visiblePosts, count, totalPages, nil
}

// UpdatePost updates a post
func (s *postService) UpdatePost(ctx context.Context, postID, userID, content, visibility string, media []string) (*models.Post, error) {
	// Validate input
	if postID == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	if visibility != "" && visibility != "public" && visibility != "private" {
		return nil, status.Error(codes.InvalidArgument, "visibility must be 'public' or 'private'")
	}

	// Get post from database
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return nil, status.Error(codes.NotFound, "post not found")
	}

	// Check if the user is the author of the post
	if post.AuthorID != userID {
		return nil, status.Error(codes.PermissionDenied, "you don't have permission to update this post")
	}

	// Update post fields
	post.Content = content
	if visibility != "" {
		post.Visibility = visibility
	}
	if media != nil {
		post.MediaArray = media
	}
	post.UpdatedAt = time.Now()

	// Save post to database
	if err := s.postRepo.Update(ctx, post); err != nil {
		s.logger.Error("Failed to update post", err)
		return nil, status.Error(codes.Internal, "failed to update post")
	}

	return post, nil
}

// DeletePost deletes a post
func (s *postService) DeletePost(ctx context.Context, postID, userID string) error {
	// Validate input
	if postID == "" {
		return status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get post from database
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return status.Error(codes.NotFound, "post not found")
	}

	// Check if the user is the author of the post
	if post.AuthorID != userID {
		return status.Error(codes.PermissionDenied, "you don't have permission to delete this post")
	}

	// Delete post from database
	if err := s.postRepo.Delete(ctx, postID); err != nil {
		s.logger.Error("Failed to delete post", err)
		return status.Error(codes.Internal, "failed to delete post")
	}

	return nil
}

// AddComment adds a comment to a post
func (s *postService) AddComment(ctx context.Context, postID, userID, authorName, authorAvatar, content string) (*models.Comment, error) {
	// Validate input
	if postID == "" {
		return nil, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}
	if content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	// Get post from database
	_, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return nil, status.Error(codes.NotFound, "post not found")
	}

	// TODO: Check if the user can comment on the post (e.g., is a friend or group member)

	// Create comment
	comment := &models.Comment{
		PostID:       postID,
		AuthorID:     userID,
		AuthorName:   authorName,
		AuthorAvatar: authorAvatar,
		Content:      content,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save comment to database
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		s.logger.Error("Failed to create comment", err)
		return nil, status.Error(codes.Internal, "failed to create comment")
	}

	// Increment comments count for the post
	if err := s.postRepo.IncrementCommentsCount(ctx, postID); err != nil {
		s.logger.Error("Failed to increment comments count", err)
		// Don't return an error here, just log it
	}

	return comment, nil
}

// GetComments retrieves comments for a post
func (s *postService) GetComments(ctx context.Context, postID string, page, limit int) ([]*models.Comment, int64, int32, error) {
	// Validate input
	if postID == "" {
		return nil, 0, 0, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get comments from database
	comments, count, err := s.commentRepo.FindByPost(ctx, postID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get comments", err)
		return nil, 0, 0, status.Error(codes.Internal, "failed to get comments")
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return comments, count, totalPages, nil
}

// DeleteComment deletes a comment
func (s *postService) DeleteComment(ctx context.Context, commentID, postID, userID string) error {
	// Validate input
	if commentID == "" {
		return status.Error(codes.InvalidArgument, "comment ID is required")
	}
	if postID == "" {
		return status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Get comment from database
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		s.logger.Error("Failed to get comment", err)
		return status.Error(codes.NotFound, "comment not found")
	}

	// Check if the comment belongs to the post
	if comment.PostID != postID {
		return status.Error(codes.InvalidArgument, "comment does not belong to the post")
	}

	// Get post from database
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return status.Error(codes.NotFound, "post not found")
	}

	// Check if the user is the author of the comment or the post
	if comment.AuthorID != userID && post.AuthorID != userID {
		return status.Error(codes.PermissionDenied, "you don't have permission to delete this comment")
	}

	// Delete comment from database
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		s.logger.Error("Failed to delete comment", err)
		return status.Error(codes.Internal, "failed to delete comment")
	}

	// Decrement comments count for the post
	if err := s.postRepo.DecrementCommentsCount(ctx, postID); err != nil {
		s.logger.Error("Failed to decrement comments count", err)
		// Don't return an error here, just log it
	}

	return nil
}

// LikePost likes a post
func (s *postService) LikePost(ctx context.Context, postID, userID string) (int32, error) {
	// Validate input
	if postID == "" {
		return 0, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return 0, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Check if the post exists
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return 0, status.Error(codes.NotFound, "post not found")
	}

	// Check if the user has already liked the post
	_, err = s.likeRepo.FindByPostAndUser(ctx, postID, userID)
	if err == nil {
		// User has already liked the post
		return int32(post.LikesCount), status.Error(codes.AlreadyExists, "you have already liked this post")
	}

	// Create like
	like := &models.Like{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	// Save like to database
	if err := s.likeRepo.Create(ctx, like); err != nil {
		s.logger.Error("Failed to create like", err)
		return 0, status.Error(codes.Internal, "failed to like post")
	}

	// Increment likes count for the post
	if err := s.postRepo.IncrementLikesCount(ctx, postID); err != nil {
		s.logger.Error("Failed to increment likes count", err)
		// Don't return an error here, just log it
	}

	// Get updated likes count
	updatedPost, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get updated post", err)
		return int32(post.LikesCount + 1), nil // Return estimated count
	}

	return int32(updatedPost.LikesCount), nil
}

// UnlikePost unlikes a post
func (s *postService) UnlikePost(ctx context.Context, postID, userID string) (int32, error) {
	// Validate input
	if postID == "" {
		return 0, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return 0, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Check if the post exists
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get post", err)
		return 0, status.Error(codes.NotFound, "post not found")
	}

	// Check if the user has liked the post
	like, err := s.likeRepo.FindByPostAndUser(ctx, postID, userID)
	if err != nil {
		// User has not liked the post
		return int32(post.LikesCount), status.Error(codes.NotFound, "you have not liked this post")
	}

	// Delete like from database
	if err := s.likeRepo.Delete(ctx, like.ID); err != nil {
		s.logger.Error("Failed to delete like", err)
		return 0, status.Error(codes.Internal, "failed to unlike post")
	}

	// Decrement likes count for the post
	if err := s.postRepo.DecrementLikesCount(ctx, postID); err != nil {
		s.logger.Error("Failed to decrement likes count", err)
		// Don't return an error here, just log it
	}

	// Get updated likes count
	updatedPost, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Error("Failed to get updated post", err)
		return int32(post.LikesCount - 1), nil // Return estimated count
	}

	return int32(updatedPost.LikesCount), nil
}

// IsLiked checks if a post is liked by a user
func (s *postService) IsLiked(ctx context.Context, postID, userID string) (bool, error) {
	// Validate input
	if postID == "" {
		return false, status.Error(codes.InvalidArgument, "post ID is required")
	}
	if userID == "" {
		return false, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Check if the user has liked the post
	_, err := s.likeRepo.FindByPostAndUser(ctx, postID, userID)
	if err != nil {
		// User has not liked the post
		return false, nil
	}

	return true, nil
}

// isPostVisibleToUser checks if a post is visible to a user
func (s *postService) isPostVisibleToUser(post *models.Post, userID string, friendIDs []string) bool {
	// Public posts are visible to everyone
	if post.Visibility == "public" {
		return true
	}

	// Private posts are visible to the author
	if userID == post.AuthorID {
		return true
	}

	// Private posts are visible to friends of the author
	if post.Visibility == "private" && friendIDs != nil {
		for _, friendID := range friendIDs {
			if friendID == post.AuthorID {
				return true
			}
		}
	}

	// Group posts are visible to group members
	if post.GroupID != "" {
		// TODO: Check if the user is a member of the group
		// For now, assume all group posts are visible to everyone
		return true
	}

	return false
}
