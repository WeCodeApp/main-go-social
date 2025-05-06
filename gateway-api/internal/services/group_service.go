package services

import (
	"context"
	"errors"

	pb "common/pb/common/proto/groups"
	"gateway-api/internal/config"
	"gateway-api/internal/models"
	"gateway-api/internal/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GroupService defines the interface for group-related operations
type GroupService interface {
	// CreateGroup creates a new group
	CreateGroup(ctx context.Context, userID string, request models.GroupCreateRequest) (*models.Group, error)

	// GetGroup retrieves a group by ID
	GetGroup(ctx context.Context, groupID, userID string) (*models.Group, error)

	// GetGroups retrieves groups with pagination and filtering
	GetGroups(ctx context.Context, userID, query string, page, limit int) (*models.GroupsResponse, error)

	// UpdateGroup updates a group
	UpdateGroup(ctx context.Context, groupID, userID string, request models.GroupUpdateRequest) (*models.Group, error)

	// DeleteGroup deletes a group
	DeleteGroup(ctx context.Context, groupID, userID string) (bool, error)

	// JoinGroup joins a group
	JoinGroup(ctx context.Context, groupID, userID string) (*models.SuccessWithCountResponse, error)

	// LeaveGroup leaves a group
	LeaveGroup(ctx context.Context, groupID, userID string) (*models.SuccessWithCountResponse, error)

	// GetGroupMembers retrieves members of a group
	GetGroupMembers(ctx context.Context, groupID string, page, limit int) (*models.GroupMembersResponse, error)

	// CreateGroupPost creates a post in a group
	CreateGroupPost(ctx context.Context, groupID, userID string, request models.GroupPostRequest) (*models.Post, error)

	// GetGroupPosts retrieves posts in a group
	GetGroupPosts(ctx context.Context, groupID, userID string, page, limit int) (*models.PostsResponse, error)
}

// groupService implements the GroupService interface
type groupService struct {
	cfg    *config.Config
	logger *logger.Logger
	client pb.GroupServiceClient
}

// createAuthContext creates a new context with the JWT token in the metadata
func (s *groupService) createAuthContext(ctx context.Context) (context.Context, error) {
	// Extract token from context
	token, ok := ctx.Value("jwt_token").(string)
	if !ok || token == "" {
		return ctx, errors.New("no token found in context")
	}

	// Create metadata with authorization header
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})

	// Create a new context with the metadata
	return metadata.NewOutgoingContext(ctx, md), nil
}

// NewGroupService creates a new group service
func NewGroupService(cfg *config.Config, logger *logger.Logger) GroupService {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(cfg.GroupsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to groups service", err)
	}

	// Create a client
	client := pb.NewGroupServiceClient(conn)

	return &groupService{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// CreateGroup creates a new group
func (s *groupService) CreateGroup(ctx context.Context, userID string, request models.GroupCreateRequest) (*models.Group, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.CreateGroup(authCtx, &pb.CreateGroupRequest{
		UserId:      userID,
		Name:        request.Name,
		Description: request.Description,
		Avatar:      request.Avatar,
	})

	if err != nil {
		s.logger.Error("Failed to create group", err)
		return nil, err
	}

	return &models.Group{
		GroupID:      resp.GroupId,
		Name:         resp.Name,
		Description:  resp.Description,
		Avatar:       resp.Avatar,
		CreatorID:    resp.CreatorId,
		CreatorName:  resp.CreatorName,
		MembersCount: resp.MembersCount,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
	}, nil
}

// GetGroup retrieves a group by ID
func (s *groupService) GetGroup(ctx context.Context, groupID, userID string) (*models.Group, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.GetGroup(authCtx, &pb.GetGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, err
	}

	return &models.Group{
		GroupID:      resp.GroupId,
		Name:         resp.Name,
		Description:  resp.Description,
		Avatar:       resp.Avatar,
		CreatorID:    resp.CreatorId,
		CreatorName:  resp.CreatorName,
		MembersCount: resp.MembersCount,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
	}, nil
}

// GetGroups retrieves groups with pagination and filtering
func (s *groupService) GetGroups(ctx context.Context, userID, query string, page, limit int) (*models.GroupsResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.GetGroups(authCtx, &pb.GetGroupsRequest{
		UserId: userID,
		Query:  query,
		Page:   int32(page),
		Limit:  int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get groups", err)
		return nil, err
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

	return &models.GroupsResponse{
		Groups:     groups,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// UpdateGroup updates a group
func (s *groupService) UpdateGroup(ctx context.Context, groupID, userID string, request models.GroupUpdateRequest) (*models.Group, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.UpdateGroup(authCtx, &pb.UpdateGroupRequest{
		GroupId:     groupID,
		UserId:      userID,
		Name:        request.Name,
		Description: request.Description,
		Avatar:      request.Avatar,
	})

	if err != nil {
		s.logger.Error("Failed to update group", err)
		return nil, err
	}

	return &models.Group{
		GroupID:      resp.GroupId,
		Name:         resp.Name,
		Description:  resp.Description,
		Avatar:       resp.Avatar,
		CreatorID:    resp.CreatorId,
		CreatorName:  resp.CreatorName,
		MembersCount: resp.MembersCount,
		CreatedAt:    resp.CreatedAt,
		UpdatedAt:    resp.UpdatedAt,
	}, nil
}

// DeleteGroup deletes a group
func (s *groupService) DeleteGroup(ctx context.Context, groupID, userID string) (bool, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return false, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.DeleteGroup(authCtx, &pb.DeleteGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		s.logger.Error("Failed to delete group", err)
		return false, err
	}

	return resp.Success, nil
}

// JoinGroup joins a group
func (s *groupService) JoinGroup(ctx context.Context, groupID, userID string) (*models.SuccessWithCountResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.JoinGroup(authCtx, &pb.JoinGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		s.logger.Error("Failed to join group", err)
		return nil, err
	}

	return &models.SuccessWithCountResponse{
		Success:      resp.Success,
		MembersCount: int(resp.MembersCount),
	}, nil
}

// LeaveGroup leaves a group
func (s *groupService) LeaveGroup(ctx context.Context, groupID, userID string) (*models.SuccessWithCountResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.LeaveGroup(authCtx, &pb.LeaveGroupRequest{
		GroupId: groupID,
		UserId:  userID,
	})

	if err != nil {
		s.logger.Error("Failed to leave group", err)
		return nil, err
	}

	return &models.SuccessWithCountResponse{
		Success:      resp.Success,
		MembersCount: int(resp.MembersCount),
	}, nil
}

// GetGroupMembers retrieves members of a group
func (s *groupService) GetGroupMembers(ctx context.Context, groupID string, page, limit int) (*models.GroupMembersResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.GetGroupMembers(authCtx, &pb.GetGroupMembersRequest{
		GroupId: groupID,
		Page:    int32(page),
		Limit:   int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get group members", err)
		return nil, err
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

	return &models.GroupMembersResponse{
		Members:    members,
		TotalCount: resp.TotalCount,
		Page:       resp.Page,
		TotalPages: resp.TotalPages,
	}, nil
}

// CreateGroupPost creates a post in a group
func (s *groupService) CreateGroupPost(ctx context.Context, groupID, userID string, request models.GroupPostRequest) (*models.Post, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.CreateGroupPost(authCtx, &pb.CreateGroupPostRequest{
		GroupId: groupID,
		UserId:  userID,
		Content: request.Content,
		Media:   request.Media,
	})

	if err != nil {
		s.logger.Error("Failed to create group post", err)
		return nil, err
	}

	return &models.Post{
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
	}, nil
}

// GetGroupPosts retrieves posts in a group
func (s *groupService) GetGroupPosts(ctx context.Context, groupID, userID string, page, limit int) (*models.PostsResponse, error) {
	// Create context with authorization metadata
	authCtx, err := s.createAuthContext(ctx)
	if err != nil {
		s.logger.Error("Failed to create auth context", err)
		return nil, err
	}

	// Call the gRPC service with the auth context
	resp, err := s.client.GetGroupPosts(authCtx, &pb.GetGroupPostsRequest{
		GroupId: groupID,
		UserId:  userID,
		Page:    int32(page),
		Limit:   int32(limit),
	})

	if err != nil {
		s.logger.Error("Failed to get group posts", err)
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
