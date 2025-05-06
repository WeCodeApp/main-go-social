package services

import (
	"context"

	pb "common/pb/common/proto/posts"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// PostService defines the interface for post-related operations
type PostService interface {
	// CreatePost creates a new post
	CreatePost(ctx context.Context, userID string, request models.PostCreateRequest) (*models.Post, error)

	// GetPost retrieves a post by ID
	GetPost(ctx context.Context, postID, userID string) (*models.Post, error)

	// GetPosts retrieves posts with pagination and filtering
	GetPosts(ctx context.Context, userID, authorID, groupID, visibility string, page, limit int) (*models.PostsResponse, error)

	// UpdatePost updates a post
	UpdatePost(ctx context.Context, postID, userID string, request models.PostUpdateRequest) (*models.Post, error)

	// DeletePost deletes a post
	DeletePost(ctx context.Context, postID, userID string) (bool, error)

	// GetComments retrieves comments for a post
	GetComments(ctx context.Context, postID string, page, limit int) (*models.CommentsResponse, error)

	// AddComment adds a comment to a post
	AddComment(ctx context.Context, postID, userID string, request models.CommentCreateRequest) (*models.Comment, error)

	// DeleteComment deletes a comment
	DeleteComment(ctx context.Context, postID, commentID, userID string) (bool, error)

	// LikePost likes a post
	LikePost(ctx context.Context, postID, userID string) (*models.LikeResponse, error)

	// UnlikePost unlikes a post
	UnlikePost(ctx context.Context, postID, userID string) (*models.LikeResponse, error)
}

// postService implements the PostService interface
type postService struct {
	cfg    *config.Config
	logger *logger.Logger
	client pb.PostServiceClient
}

// NewPostService creates a new post service
func NewPostService(cfg *config.Config, logger *logger.Logger) PostService {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.PostsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to posts service", err)
	}

	// Create a client
	client := pb.NewPostServiceClient(conn)

	return &postService{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(ctx context.Context, userID string, request models.PostCreateRequest) (*models.Post, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.CreatePost(ctxWithToken, &pb.CreatePostRequest{
		UserId:     userID,
		Content:    request.Content,
		Visibility: request.Visibility,
		GroupId:    "", // GroupID is not in the PostCreateRequest model, we'll need to update the model
		Media:      request.Media,
	})

	if err != nil {
		s.logger.Error("Failed to create post", err)
		return nil, err
	}

	return &models.Post{
		PostID:        resp.PostId,
		AuthorID:      resp.AuthorId,
		AuthorName:    resp.AuthorName,
		AuthorAvatar:  resp.AuthorAvatar,
		Content:       resp.Content,
		Visibility:    resp.Visibility,
		Media:         resp.Media,
		LikesCount:    resp.LikesCount,
		CommentsCount: resp.CommentsCount,
		IsLiked:       resp.IsLiked,
		CreatedAt:     resp.CreatedAt,
		UpdatedAt:     resp.UpdatedAt,
	}, nil
}

// GetPost retrieves a post by ID
func (s *postService) GetPost(ctx context.Context, postID, userID string) (*models.Post, error) {
	// Call the gRPC service
	resp, err := s.client.GetPost(context.Background(), &pb.GetPostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		s.logger.Error("Failed to get post", err)
		return nil, err
	}

	return &models.Post{
		PostID:        resp.PostId,
		AuthorID:      resp.AuthorId,
		AuthorName:    resp.AuthorName,
		AuthorAvatar:  resp.AuthorAvatar,
		Content:       resp.Content,
		Visibility:    resp.Visibility,
		Media:         resp.Media,
		LikesCount:    resp.LikesCount,
		CommentsCount: resp.CommentsCount,
		IsLiked:       resp.IsLiked,
		CreatedAt:     resp.CreatedAt,
		UpdatedAt:     resp.UpdatedAt,
	}, nil
}

// GetPosts retrieves posts with pagination and filtering
func (s *postService) GetPosts(ctx context.Context, userID, authorID, groupID, visibility string, page, limit int) (*models.PostsResponse, error) {
	// Call the gRPC service
	resp, err := s.client.GetPosts(context.Background(), &pb.GetPostsRequest{
		UserId:     userID,
		AuthorId:   authorID,
		GroupId:    groupID,
		Visibility: visibility,
		Page:       int32(page),
		Limit:      int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get posts", err)
		return nil, err
	}

	// Convert posts to model format
	posts := make([]models.Post, len(resp.Posts))
	for i, post := range resp.Posts {
		posts[i] = models.Post{
			PostID:        post.PostId,
			AuthorID:      post.AuthorId,
			AuthorName:    post.AuthorName,
			AuthorAvatar:  post.AuthorAvatar,
			Content:       post.Content,
			Visibility:    post.Visibility,
			Media:         post.Media,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			IsLiked:       post.IsLiked,
			CreatedAt:     post.CreatedAt,
			UpdatedAt:     post.UpdatedAt,
		}
	}

	return &models.PostsResponse{
		Posts:      posts,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// UpdatePost updates a post
func (s *postService) UpdatePost(ctx context.Context, postID, userID string, request models.PostUpdateRequest) (*models.Post, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.UpdatePost(ctxWithToken, &pb.UpdatePostRequest{
		PostId:     postID,
		UserId:     userID,
		Content:    request.Content,
		Visibility: request.Visibility,
		Media:      request.Media,
	})

	if err != nil {
		s.logger.Error("Failed to update post", err)
		return nil, err
	}

	return &models.Post{
		PostID:        resp.PostId,
		AuthorID:      resp.AuthorId,
		AuthorName:    resp.AuthorName,
		AuthorAvatar:  resp.AuthorAvatar,
		Content:       resp.Content,
		Visibility:    resp.Visibility,
		Media:         resp.Media,
		LikesCount:    resp.LikesCount,
		CommentsCount: resp.CommentsCount,
		IsLiked:       resp.IsLiked,
		CreatedAt:     resp.CreatedAt,
		UpdatedAt:     resp.UpdatedAt,
	}, nil
}

// DeletePost deletes a post
func (s *postService) DeletePost(ctx context.Context, postID, userID string) (bool, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.DeletePost(ctxWithToken, &pb.DeletePostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		s.logger.Error("Failed to delete post", err)
		return false, err
	}

	return resp.Success, nil
}

// GetComments retrieves comments for a post
func (s *postService) GetComments(ctx context.Context, postID string, page, limit int) (*models.CommentsResponse, error) {
	// Call the gRPC service
	resp, err := s.client.GetComments(context.Background(), &pb.GetCommentsRequest{
		PostId: postID,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get comments", err)
		return nil, err
	}

	// Convert comments to model format
	comments := make([]models.Comment, len(resp.Comments))
	for i, comment := range resp.Comments {
		comments[i] = models.Comment{
			CommentID:    comment.CommentId,
			PostID:       comment.PostId,
			AuthorID:     comment.AuthorId,
			AuthorName:   comment.AuthorName,
			AuthorAvatar: comment.AuthorAvatar,
			Content:      comment.Content,
			CreatedAt:    comment.CreatedAt,
		}
	}

	return &models.CommentsResponse{
		Comments:   comments,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// AddComment adds a comment to a post
func (s *postService) AddComment(ctx context.Context, postID, userID string, request models.CommentCreateRequest) (*models.Comment, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.AddComment(ctxWithToken, &pb.AddCommentRequest{
		PostId:  postID,
		UserId:  userID,
		Content: request.Content,
	})

	if err != nil {
		s.logger.Error("Failed to add comment", err)
		return nil, err
	}

	return &models.Comment{
		CommentID:    resp.CommentId,
		PostID:       resp.PostId,
		AuthorID:     resp.AuthorId,
		AuthorName:   resp.AuthorName,
		AuthorAvatar: resp.AuthorAvatar,
		Content:      resp.Content,
		CreatedAt:    resp.CreatedAt,
	}, nil
}

// DeleteComment deletes a comment
func (s *postService) DeleteComment(ctx context.Context, postID, commentID, userID string) (bool, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.DeleteComment(ctxWithToken, &pb.DeleteCommentRequest{
		CommentId: commentID,
		UserId:    userID,
		PostId:    postID,
	})

	if err != nil {
		s.logger.Error("Failed to delete comment", err)
		return false, err
	}

	return resp.Success, nil
}

// LikePost likes a post
func (s *postService) LikePost(ctx context.Context, postID, userID string) (*models.LikeResponse, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.LikePost(ctxWithToken, &pb.LikePostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		s.logger.Error("Failed to like post", err)
		return nil, err
	}

	return &models.LikeResponse{
		Success:    resp.Success,
		LikesCount: int(resp.LikesCount),
	}, nil
}

// UnlikePost unlikes a post
func (s *postService) UnlikePost(ctx context.Context, postID, userID string) (*models.LikeResponse, error) {
	// Get JWT token from context
	token, _ := ctx.Value("jwt_token").(string)

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(ctx, md)

	// Call the gRPC service with the context containing the token
	resp, err := s.client.UnlikePost(ctxWithToken, &pb.UnlikePostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		s.logger.Error("Failed to unlike post", err)
		return nil, err
	}

	return &models.LikeResponse{
		Success:    resp.Success,
		LikesCount: int(resp.LikesCount),
	}, nil
}
