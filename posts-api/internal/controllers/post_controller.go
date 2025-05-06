package controllers

import (
	pb "common/pb/common/proto/posts"
	"context"
	"post-api/internal/models"
	"post-api/internal/services"
	"post-api/internal/utils/logger"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// PostController handles gRPC requests for post-related operations
type PostController struct {
	pb.UnimplementedPostServiceServer
	postService services.PostService
	logger      *logger.Logger
}

// NewPostController creates a new post controller
func NewPostController(postService services.PostService, logger *logger.Logger) *PostController {
	return &PostController{
		postService: postService,
		logger:      logger,
	}
}

// getUserIDFromContext extracts the user ID from the context metadata
func (c *PostController) getUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("user_id")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing user ID")
	}

	return values[0], nil
}

// getFriendIDsFromContext extracts the friend IDs from the context metadata
func (c *PostController) getFriendIDsFromContext(ctx context.Context) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	values := md.Get("friend_ids")
	if len(values) == 0 {
		return nil
	}

	return values
}

// CreatePost creates a new post
func (c *PostController) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	c.logger.Info("CreatePost request received", "user_id", req.UserId, "visibility", req.Visibility)

	// Create post using the service
	post, err := c.postService.CreatePost(ctx, req.UserId, req.Content, req.Visibility, req.GroupId, req.Media)
	if err != nil {
		c.logger.Error("Failed to create post", err)
		return nil, err
	}

	// Convert post model to gRPC response
	return c.convertPostToResponse(post, false), nil
}

// GetPost retrieves a post by ID
func (c *PostController) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.PostResponse, error) {
	c.logger.Info("GetPost request received", "post_id", req.PostId, "user_id", req.UserId)

	// Get post using the service
	post, isLiked, err := c.postService.GetPost(ctx, req.PostId, req.UserId)
	if err != nil {
		c.logger.Error("Failed to get post", err)
		return nil, err
	}

	// Convert post model to gRPC response
	return c.convertPostToResponse(post, isLiked), nil
}

// GetPosts retrieves posts with pagination and filtering
func (c *PostController) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	c.logger.Info("GetPosts request received",
		"user_id", req.UserId,
		"author_id", req.AuthorId,
		"group_id", req.GroupId,
		"visibility", req.Visibility,
		"page", req.Page,
		"limit", req.Limit)

	// Get friend IDs from context
	friendIDs := c.getFriendIDsFromContext(ctx)

	// Get posts using the service
	posts, totalCount, totalPages, err := c.postService.GetPosts(
		ctx,
		req.UserId,
		req.AuthorId,
		req.GroupId,
		req.Visibility,
		int(req.Page),
		int(req.Limit),
		friendIDs,
	)
	if err != nil {
		c.logger.Error("Failed to get posts", err)
		return nil, err
	}

	// Convert post models to gRPC responses
	postResponses := make([]*pb.PostResponse, len(posts))
	for i, post := range posts {
		// Check if the post is liked by the user
		isLiked := false
		if req.UserId != "" {
			isLiked, _ = c.postService.IsLiked(ctx, post.ID, req.UserId)
		}
		postResponses[i] = c.convertPostToResponse(post, isLiked)
	}

	return &pb.GetPostsResponse{
		Posts:      postResponses,
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}

// UpdatePost updates a post
func (c *PostController) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error) {
	c.logger.Info("UpdatePost request received", "post_id", req.PostId, "user_id", req.UserId)

	// Update post using the service
	post, err := c.postService.UpdatePost(ctx, req.PostId, req.UserId, req.Content, req.Visibility, req.Media)
	if err != nil {
		c.logger.Error("Failed to update post", err)
		return nil, err
	}

	// Check if the post is liked by the user
	isLiked, _ := c.postService.IsLiked(ctx, post.ID, req.UserId)

	// Convert post model to gRPC response
	return c.convertPostToResponse(post, isLiked), nil
}

// DeletePost deletes a post
func (c *PostController) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	c.logger.Info("DeletePost request received", "post_id", req.PostId, "user_id", req.UserId)

	// Delete post using the service
	err := c.postService.DeletePost(ctx, req.PostId, req.UserId)
	if err != nil {
		c.logger.Error("Failed to delete post", err)
		return nil, err
	}

	return &pb.DeletePostResponse{
		Success: true,
	}, nil
}

// AddComment adds a comment to a post
func (c *PostController) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.CommentResponse, error) {
	c.logger.Info("AddComment request received", "post_id", req.PostId, "user_id", req.UserId)

	// TODO: Get user info from users service
	authorName := "User " + req.UserId // Placeholder
	authorAvatar := ""                 // Placeholder

	// Add comment using the service
	comment, err := c.postService.AddComment(ctx, req.PostId, req.UserId, authorName, authorAvatar, req.Content)
	if err != nil {
		c.logger.Error("Failed to add comment", err)
		return nil, err
	}

	// Convert comment model to gRPC response
	return c.convertCommentToResponse(comment), nil
}

// GetComments retrieves comments for a post
func (c *PostController) GetComments(ctx context.Context, req *pb.GetCommentsRequest) (*pb.GetCommentsResponse, error) {
	c.logger.Info("GetComments request received", "post_id", req.PostId, "page", req.Page, "limit", req.Limit)

	// Get comments using the service
	comments, totalCount, totalPages, err := c.postService.GetComments(ctx, req.PostId, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get comments", err)
		return nil, err
	}

	// Convert comment models to gRPC responses
	commentResponses := make([]*pb.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = c.convertCommentToResponse(comment)
	}

	return &pb.GetCommentsResponse{
		Comments:   commentResponses,
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}

// DeleteComment deletes a comment
func (c *PostController) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	c.logger.Info("DeleteComment request received", "comment_id", req.CommentId, "post_id", req.PostId, "user_id", req.UserId)

	// Delete comment using the service
	err := c.postService.DeleteComment(ctx, req.CommentId, req.PostId, req.UserId)
	if err != nil {
		c.logger.Error("Failed to delete comment", err)
		return nil, err
	}

	return &pb.DeleteCommentResponse{
		Success: true,
	}, nil
}

// LikePost likes a post
func (c *PostController) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	c.logger.Info("LikePost request received", "post_id", req.PostId, "user_id", req.UserId)

	// Like post using the service
	likesCount, err := c.postService.LikePost(ctx, req.PostId, req.UserId)
	if err != nil {
		c.logger.Error("Failed to like post", err)
		return nil, err
	}

	return &pb.LikePostResponse{
		Success:    true,
		LikesCount: likesCount,
	}, nil
}

// UnlikePost unlikes a post
func (c *PostController) UnlikePost(ctx context.Context, req *pb.UnlikePostRequest) (*pb.UnlikePostResponse, error) {
	c.logger.Info("UnlikePost request received", "post_id", req.PostId, "user_id", req.UserId)

	// Unlike post using the service
	likesCount, err := c.postService.UnlikePost(ctx, req.PostId, req.UserId)
	if err != nil {
		c.logger.Error("Failed to unlike post", err)
		return nil, err
	}

	return &pb.UnlikePostResponse{
		Success:    true,
		LikesCount: likesCount,
	}, nil
}

// convertPostToResponse converts a post model to a gRPC response
func (c *PostController) convertPostToResponse(post *models.Post, isLiked bool) *pb.PostResponse {
	return &pb.PostResponse{
		PostId:        post.ID,
		AuthorId:      post.AuthorID,
		AuthorName:    post.AuthorName,
		AuthorAvatar:  post.AuthorAvatar,
		Content:       post.Content,
		Visibility:    post.Visibility,
		GroupId:       post.GroupID,
		GroupName:     post.GroupName,
		Media:         post.MediaArray,
		LikesCount:    int32(post.LikesCount),
		CommentsCount: int32(post.CommentsCount),
		IsLiked:       isLiked,
		CreatedAt:     post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     post.UpdatedAt.Format(time.RFC3339),
	}
}

// convertCommentToResponse converts a comment model to a gRPC response
func (c *PostController) convertCommentToResponse(comment *models.Comment) *pb.CommentResponse {
	return &pb.CommentResponse{
		CommentId:    comment.ID,
		PostId:       comment.PostID,
		AuthorId:     comment.AuthorID,
		AuthorName:   comment.AuthorName,
		AuthorAvatar: comment.AuthorAvatar,
		Content:      comment.Content,
		CreatedAt:    comment.CreatedAt.Format(time.RFC3339),
	}
}
