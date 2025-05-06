package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/services"
	"gateway-api/internal/utils/logger"
)

// PostController handles post-related requests
type PostController struct {
	cfg         *config.Config
	logger      *logger.Logger
	postService services.PostService
}

// NewPostController creates a new post controller
func NewPostController(cfg *config.Config, logger *logger.Logger) *PostController {
	postService := services.NewPostService(cfg, logger)

	return &PostController{
		cfg:         cfg,
		logger:      logger,
		postService: postService,
	}
}

// CreatePost handles post creation
// @Summary Create a post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.PostCreateRequest true "Post creation request"
// @Success 201 {object} models.Post "Post created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts [post]
func (c *PostController) CreatePost(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var request models.PostCreateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Call the post service
	resp, err := c.postService.CreatePost(ctx, userID, request)

	if err != nil {
		c.logger.Error("Failed to create post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create post",
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetPost handles retrieving a post by ID
// @Summary Get a post
// @Description Get a post by ID
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.Post "Post"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id} [get]
func (c *PostController) GetPost(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID") // May be empty if not authenticated

	// Call the post service
	resp, err := c.postService.GetPost(ctx, postID, userID)

	if err != nil {
		c.logger.Error("Failed to get post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get post",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetPosts handles retrieving posts with pagination and filtering
// @Summary Get posts
// @Description Get posts with pagination and filtering
// @Tags posts
// @Produce json
// @Param author_id query string false "Filter posts by author ID"
// @Param group_id query string false "Filter posts by group ID"
// @Param visibility query string false "Filter posts by visibility" Enums(public, private)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of posts per page" default(10)
// @Success 200 {object} models.PostsResponse "Posts"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts [get]
func (c *PostController) GetPosts(ctx *gin.Context) {
	userID := ctx.GetString("userID") // May be empty if not authenticated
	authorID := ctx.Query("author_id")
	groupID := ctx.Query("group_id")
	visibility := ctx.Query("visibility")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Call the post service
	resp, err := c.postService.GetPosts(ctx, userID, authorID, groupID, visibility, page, limit)

	if err != nil {
		c.logger.Error("Failed to get posts", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get posts",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdatePost handles updating a post
// @Summary Update a post
// @Description Update a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body models.PostUpdateRequest true "Update post request"
// @Success 200 {object} models.Post "Post updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id} [put]
func (c *PostController) UpdatePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID")

	var request models.PostUpdateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Call the post service
	resp, err := c.postService.UpdatePost(ctx, postID, userID, request)

	if err != nil {
		c.logger.Error("Failed to update post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update post",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeletePost handles deleting a post
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} models.SuccessResponse "Post deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id} [delete]
func (c *PostController) DeletePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID")

	// Call the post service
	success, err := c.postService.DeletePost(ctx, postID, userID)

	if err != nil {
		c.logger.Error("Failed to delete post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete post",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: success,
	})
}

// GetComments handles retrieving comments for a post
// @Summary Get comments for a post
// @Description Get comments for a post with pagination
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of comments per page" default(10)
// @Success 200 {object} models.CommentsResponse "Comments"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments [get]
func (c *PostController) GetComments(ctx *gin.Context) {
	postID := ctx.Param("id")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Call the post service
	resp, err := c.postService.GetComments(ctx, postID, page, limit)

	if err != nil {
		c.logger.Error("Failed to get comments", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get comments",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// AddComment handles adding a comment to a post
// @Summary Add a comment to a post
// @Description Add a comment to a post
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body models.CommentCreateRequest true "Comment request"
// @Success 201 {object} models.Comment "Comment added successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments [post]
func (c *PostController) AddComment(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID")

	var request models.CommentCreateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Call the post service
	resp, err := c.postService.AddComment(ctx, postID, userID, request)

	if err != nil {
		c.logger.Error("Failed to add comment", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to add comment",
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteComment handles deleting a comment
// @Summary Delete a comment
// @Description Delete a comment from a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param commentId path string true "Comment ID"
// @Success 200 {object} models.SuccessResponse "Comment deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post or comment not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id}/comments/{commentId} [delete]
func (c *PostController) DeleteComment(ctx *gin.Context) {
	postID := ctx.Param("id")
	commentID := ctx.Param("commentId")
	userID := ctx.GetString("userID")

	// Call the post service
	success, err := c.postService.DeleteComment(ctx, postID, commentID, userID)

	if err != nil {
		c.logger.Error("Failed to delete comment", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete comment",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: success,
	})
}

// LikePost handles liking a post
// @Summary Like a post
// @Description Like a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} models.LikeResponse "Post liked successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id}/like [post]
func (c *PostController) LikePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID")

	// Call the post service
	resp, err := c.postService.LikePost(ctx, postID, userID)

	if err != nil {
		c.logger.Error("Failed to like post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to like post",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UnlikePost handles unliking a post
// @Summary Unlike a post
// @Description Unlike a post
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Success 200 {object} models.LikeResponse "Post unliked successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /posts/{id}/like [delete]
func (c *PostController) UnlikePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	userID := ctx.GetString("userID")

	// Call the post service
	resp, err := c.postService.UnlikePost(ctx, postID, userID)

	if err != nil {
		c.logger.Error("Failed to unlike post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to unlike post",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
