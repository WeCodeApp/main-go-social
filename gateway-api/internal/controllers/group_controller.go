package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "common/pb/common/proto/groups"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
)

// GroupController handles group-related requests
type GroupController struct {
	cfg    *config.Config
	logger *logger.Logger
	client pb.GroupServiceClient
}

// NewGroupController creates a new group controller
func NewGroupController(cfg *config.Config, logger *logger.Logger) *GroupController {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.GroupsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to groups service", err)
	}

	// Create a client
	client := pb.NewGroupServiceClient(conn)

	return &GroupController{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// CreateGroup handles group creation
// @Summary Create a group
// @Description Create a new group
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.GroupCreateRequest true "Group creation request"
// @Success 201 {object} models.Group "Group created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups [post]
func (c *GroupController) CreateGroup(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	var request models.GroupCreateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.CreateGroup(ctxWithToken, &pb.CreateGroupRequest{
		UserId:      userID,
		Name:        request.Name,
		Description: request.Description,
		Avatar:      request.Avatar,
	})

	if err != nil {
		c.logger.Error("Failed to create group", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create group",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"group_id":      resp.GroupId,
		"name":          resp.Name,
		"description":   resp.Description,
		"avatar":        resp.Avatar,
		"creator_id":    resp.CreatorId,
		"creator_name":  resp.CreatorName,
		"members_count": resp.MembersCount,
		"posts_count":   resp.PostsCount,
		"is_member":     resp.IsMember,
		"created_at":    resp.CreatedAt,
		"updated_at":    resp.UpdatedAt,
	})
}

// GetGroup handles retrieving a group by ID
// @Summary Get a group
// @Description Get a group by ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} models.Group "Group"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id} [get]
func (c *GroupController) GetGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID") // May be empty if not authenticated
	token := ctx.GetString("jwt_token")

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.GetGroup(ctxWithToken, &pb.GetGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		c.logger.Error("Failed to get group", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get group",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Group{
		GroupID:      resp.GroupId,
		Name:         resp.Name,
		Description:  resp.Description,
		Avatar:       resp.Avatar,
		CreatorID:    resp.CreatorId,
		CreatorName:  resp.CreatorName,
		MembersCount: resp.MembersCount,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
	})
}

// GetGroups handles retrieving groups with pagination and filtering
// @Summary Get groups
// @Description Get groups with pagination and filtering
// @Tags groups
// @Produce json
// @Param query query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of groups per page" default(10)
// @Success 200 {object} models.GroupsResponse "Groups"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups [get]
func (c *GroupController) GetGroups(ctx *gin.Context) {
	userID := ctx.GetString("userID") // May be empty if not authenticated
	query := ctx.Query("query")
	token := ctx.GetString("jwt_token")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.GetGroups(ctxWithToken, &pb.GetGroupsRequest{
		UserId: userID,
		Query:  query,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		c.logger.Error("Failed to get groups", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get groups",
		})
		return
	}

	// Convert groups to model format
	groups := make([]models.Group, len(resp.Groups))
	for i, group := range resp.Groups {
		groups[i] = models.Group{
			GroupID:      group.GroupId,
			Name:         group.Name,
			Description:  group.Description,
			Avatar:       group.Avatar,
			CreatorID:    group.CreatorId,
			CreatorName:  group.CreatorName,
			MembersCount: group.MembersCount,
			CreatedAt:    group.CreatedAt,
			UpdatedAt:    group.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, models.GroupsResponse{
		Groups:     groups,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	})
}

// UpdateGroup handles updating a group
// @Summary Update a group
// @Description Update a group by ID
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Param request body models.GroupUpdateRequest true "Update group request"
// @Success 200 {object} models.Group "Group updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id} [put]
func (c *GroupController) UpdateGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	var request models.GroupUpdateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.UpdateGroup(ctxWithToken, &pb.UpdateGroupRequest{
		GroupId:     groupID,
		UserId:      userID,
		Name:        request.Name,
		Description: request.Description,
		Avatar:      request.Avatar,
	})

	if err != nil {
		c.logger.Error("Failed to update group", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update group",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Group{
		GroupID:      resp.GroupId,
		Name:         resp.Name,
		Description:  resp.Description,
		Avatar:       resp.Avatar,
		CreatorID:    resp.CreatorId,
		CreatorName:  resp.CreatorName,
		MembersCount: resp.MembersCount,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
	})
}

// DeleteGroup handles deleting a group
// @Summary Delete a group
// @Description Delete a group by ID
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} models.SuccessResponse "Group deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id} [delete]
func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.DeleteGroup(ctxWithToken, &pb.DeleteGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		c.logger.Error("Failed to delete group", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete group",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Success: resp.Success,
	})
}

// JoinGroup handles joining a group
// @Summary Join a group
// @Description Join a group
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} models.SuccessWithCountResponse "Group joined successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id}/members [post]
func (c *GroupController) JoinGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.JoinGroup(ctxWithToken, &pb.JoinGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		c.logger.Error("Failed to join group", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to join group",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessWithCountResponse{
		Success:      resp.Success,
		MembersCount: int(resp.MembersCount),
	})
}

// LeaveGroup handles leaving a group
// @Summary Leave a group
// @Description Leave a group
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Success 200 {object} models.SuccessWithCountResponse "Group left successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id}/members [delete]
func (c *GroupController) LeaveGroup(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.LeaveGroup(ctxWithToken, &pb.LeaveGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		c.logger.Error("Failed to leave group", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to leave group",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessWithCountResponse{
		Success:      resp.Success,
		MembersCount: int(resp.MembersCount),
	})
}

// GetGroupMembers handles retrieving members of a group
// @Summary Get group members
// @Description Get members of a group with pagination
// @Tags groups
// @Produce json
// @Param id path string true "Group ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of members per page" default(10)
// @Success 200 {object} models.GroupMembersResponse "Group members with pagination"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id}/members [get]
func (c *GroupController) GetGroupMembers(ctx *gin.Context) {
	groupID := ctx.Param("id")
	token := ctx.GetString("jwt_token")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.GetGroupMembers(ctxWithToken, &pb.GetGroupMembersRequest{
		GroupId: groupID,
		Page:    int32(page),
		Limit:   int32(limit),
	})

	if err != nil {
		c.logger.Error("Failed to get group members", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get group members",
		})
		return
	}

	// Convert members to model format
	members := make([]models.GroupMember, len(resp.Members))
	for i, member := range resp.Members {
		members[i] = models.GroupMember{
			UserID:   member.UserId,
			Name:     member.Name,
			Avatar:   member.Avatar,
			Role:     member.Role,
			JoinedAt: member.JoinedAt,
		}
	}

	ctx.JSON(http.StatusOK, models.GroupMembersResponse{
		Members:    members,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	})
}

// CreateGroupPost handles creating a post in a group
// @Summary Create a post in a group
// @Description Create a new post in a group
// @Tags groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Param request body models.GroupPostRequest true "Post content"
// @Success 201 {object} models.Post "Created post"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Not a member of the group"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id}/posts [post]
func (c *GroupController) CreateGroupPost(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID")
	token := ctx.GetString("jwt_token")

	var request models.GroupPostRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.CreateGroupPost(ctxWithToken, &pb.CreateGroupPostRequest{
		GroupId: groupID,
		UserId:  userID,
		Content: request.Content,
		Media:   request.Media,
	})

	if err != nil {
		c.logger.Error("Failed to create group post", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create group post",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Post{
		PostID:        resp.PostId,
		AuthorID:      resp.AuthorId,
		AuthorName:    resp.AuthorName,
		AuthorAvatar:  resp.AuthorAvatar,
		Content:       resp.Content,
		Media:         resp.Media,
		LikesCount:    resp.LikesCount,
		CommentsCount: resp.CommentsCount,
		IsLiked:       resp.IsLiked,
		CreatedAt:     resp.CreatedAt,
		UpdatedAt:     resp.UpdatedAt,
	})
}

// GetGroupPosts handles retrieving posts in a group
// @Summary Get posts in a group
// @Description Get posts in a group with pagination
// @Tags groups
// @Produce json
// @Security BearerAuth
// @Param id path string true "Group ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of posts per page" default(10)
// @Success 200 {object} models.PostsResponse "Group posts with pagination"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Not a member of the group"
// @Failure 404 {object} models.ErrorResponse "Group not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /groups/{id}/posts [get]
func (c *GroupController) GetGroupPosts(ctx *gin.Context) {
	groupID := ctx.Param("id")
	userID := ctx.GetString("userID") // May be empty if not authenticated
	token := ctx.GetString("jwt_token")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Create metadata with authorization token
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create new context with metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)

	// Call the gRPC service with the context containing the token
	resp, err := c.client.GetGroupPosts(ctxWithToken, &pb.GetGroupPostsRequest{
		GroupId: groupID,
		UserId:  userID,
		Page:    int32(page),
		Limit:   int32(limit),
	})

	if err != nil {
		c.logger.Error("Failed to get group posts", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get group posts",
		})
		return
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
			Media:         post.Media,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			IsLiked:       post.IsLiked,
			CreatedAt:     post.CreatedAt,
			UpdatedAt:     post.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, models.PostsResponse{
		Posts:      posts,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	})
}
