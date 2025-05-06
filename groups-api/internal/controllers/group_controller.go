package controllers

import (
	pb "common/pb/common/proto/groups"
	"context"
	"groups-api/internal/services"
	"groups-api/internal/utils/errors"
	"groups-api/internal/utils/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GroupController handles gRPC requests for group-related operations
type GroupController struct {
	pb.UnimplementedGroupServiceServer
	service services.GroupService
	logger  *logger.Logger
}

// NewGroupController creates a new group controller
func NewGroupController(service services.GroupService, logger *logger.Logger) *GroupController {
	return &GroupController{
		service: service,
		logger:  logger,
	}
}

// CreateGroup creates a new group
func (c *GroupController) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "group name is required")
	}

	// Create group
	group, err := c.service.CreateGroup(ctx, userID, req.Name, req.Description, req.Avatar)
	if err != nil {
		c.logger.Error("Failed to create group", err)
		return nil, status.Error(codes.Internal, "failed to create group")
	}

	// Create response
	return &pb.GroupResponse{
		GroupId:      group.ID,
		Name:         group.Name,
		Description:  group.Description,
		Avatar:       group.Avatar,
		CreatorId:    group.CreatorID,
		CreatorName:  "", // Would need to fetch from users service
		MembersCount: 1,  // Creator is the only member initially
		PostsCount:   0,  // No posts initially
		IsMember:     true,
		CreatedAt:    group.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    group.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetGroup retrieves a group by ID
func (c *GroupController) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GroupResponse, error) {
	// Get user ID from context or request
	userID := req.UserId
	if userID == "" {
		if id, ok := ctx.Value("userID").(string); ok {
			userID = id
		}
	}

	// Get group
	group, membersCount, postsCount, isMember, err := c.service.GetGroup(ctx, req.GroupId, userID)
	if err != nil {
		c.logger.Error("Failed to get group", err)
		return nil, status.Error(codes.NotFound, "group not found")
	}

	// Create response
	return &pb.GroupResponse{
		GroupId:      group.ID,
		Name:         group.Name,
		Description:  group.Description,
		Avatar:       group.Avatar,
		CreatorId:    group.CreatorID,
		CreatorName:  "", // Would need to fetch from users service
		MembersCount: membersCount,
		PostsCount:   postsCount,
		IsMember:     isMember,
		CreatedAt:    group.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    group.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetGroups retrieves groups with pagination and filtering
func (c *GroupController) GetGroups(ctx context.Context, req *pb.GetGroupsRequest) (*pb.GetGroupsResponse, error) {
	// Get user ID from context or request
	userID := req.UserId
	if userID == "" {
		if id, ok := ctx.Value("userID").(string); ok {
			userID = id
		}
	}

	// Get groups
	groups, totalCount, totalPages, err := c.service.GetGroups(ctx, userID, req.Query, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get groups", err)
		return nil, status.Error(codes.Internal, "failed to get groups")
	}

	// Create response
	response := &pb.GetGroupsResponse{
		Groups:     make([]*pb.GroupResponse, 0, len(groups)),
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}

	// Add groups to response
	for _, group := range groups {
		// Get group details
		groupDetails, membersCount, postsCount, isMember, err := c.service.GetGroup(ctx, group.ID, userID)
		if err != nil {
			c.logger.Error("Failed to get group details", err)
			continue
		}

		response.Groups = append(response.Groups, &pb.GroupResponse{
			GroupId:      groupDetails.ID,
			Name:         groupDetails.Name,
			Description:  groupDetails.Description,
			Avatar:       groupDetails.Avatar,
			CreatorId:    groupDetails.CreatorID,
			CreatorName:  "", // Would need to fetch from users service
			MembersCount: membersCount,
			PostsCount:   postsCount,
			IsMember:     isMember,
			CreatedAt:    groupDetails.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    groupDetails.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// UpdateGroup updates a group
func (c *GroupController) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Update group
	group, err := c.service.UpdateGroup(ctx, req.GroupId, userID, req.Name, req.Description, req.Avatar)
	if err != nil {
		c.logger.Error("Failed to update group", err)
		return nil, status.Error(codes.Internal, "failed to update group")
	}

	// Get group details
	groupDetails, membersCount, postsCount, isMember, err := c.service.GetGroup(ctx, group.ID, userID)
	if err != nil {
		c.logger.Error("Failed to get group details", err)
		return nil, status.Error(codes.Internal, "failed to get group details")
	}

	// Create response
	return &pb.GroupResponse{
		GroupId:      groupDetails.ID,
		Name:         groupDetails.Name,
		Description:  groupDetails.Description,
		Avatar:       groupDetails.Avatar,
		CreatorId:    groupDetails.CreatorID,
		CreatorName:  "", // Would need to fetch from users service
		MembersCount: membersCount,
		PostsCount:   postsCount,
		IsMember:     isMember,
		CreatedAt:    groupDetails.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    groupDetails.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteGroup deletes a group
func (c *GroupController) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Delete group
	err := c.service.DeleteGroup(ctx, req.GroupId, userID)
	if err != nil {
		c.logger.Error("Failed to delete group", err)
		return nil, status.Error(codes.Internal, "failed to delete group")
	}

	// Create response
	return &pb.DeleteGroupResponse{
		Success: true,
	}, nil
}

// JoinGroup adds a user to a group
func (c *GroupController) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Join group
	success, membersCount, err := c.service.JoinGroup(ctx, req.GroupId, userID)
	if err != nil {
		c.logger.Error("Failed to join group", err)
		return nil, status.Error(codes.Internal, "failed to join group")
	}

	// Create response
	return &pb.JoinGroupResponse{
		Success:      success,
		MembersCount: membersCount,
	}, nil
}

// LeaveGroup removes a user from a group
func (c *GroupController) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Leave group
	success, membersCount, err := c.service.LeaveGroup(ctx, req.GroupId, userID)
	if err != nil {
		c.logger.Error("Failed to leave group", err)
		return nil, status.Error(codes.Internal, "failed to leave group")
	}

	// Create response
	return &pb.LeaveGroupResponse{
		Success:      success,
		MembersCount: membersCount,
	}, nil
}

// GetGroupMembers retrieves members of a group
func (c *GroupController) GetGroupMembers(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error) {
	// Get members
	members, totalCount, totalPages, err := c.service.GetGroupMembers(ctx, req.GroupId, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get group members", err)
		return nil, status.Error(codes.Internal, "failed to get group members")
	}

	// Create response
	response := &pb.GetGroupMembersResponse{
		Members:    make([]*pb.GroupMemberResponse, 0, len(members)),
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}

	// Add members to response
	for _, member := range members {
		response.Members = append(response.Members, &pb.GroupMemberResponse{
			UserId:   member.UserID,
			Name:     "", // Would need to fetch from users service
			Avatar:   "", // Would need to fetch from users service
			Role:     member.Role,
			JoinedAt: member.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

// CreateGroupPost creates a post in a group
func (c *GroupController) CreateGroupPost(ctx context.Context, req *pb.CreateGroupPostRequest) (*pb.GroupPostResponse, error) {
	// Get user ID from context
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		c.logger.Error("Failed to get user ID from context", nil)
		return nil, errors.ErrUnauthenticated
	}

	// Create post
	post, err := c.service.CreateGroupPost(ctx, req.GroupId, userID, req.Content, req.Media)
	if err != nil {
		c.logger.Error("Failed to create group post", err)
		return nil, status.Error(codes.Internal, "failed to create group post")
	}

	// Create response
	response := &pb.GroupPostResponse{
		PostId:        post.ID,
		GroupId:       post.GroupID,
		AuthorId:      post.AuthorID,
		AuthorName:    "", // Would need to fetch from users service
		AuthorAvatar:  "", // Would need to fetch from users service
		Content:       post.Content,
		Media:         make([]string, 0, len(post.Media)),
		LikesCount:    int32(len(post.Likes)),
		CommentsCount: int32(len(post.Comments)),
		IsLiked:       false, // Would need to check if the user liked the post
		CreatedAt:     post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add media to response
	for _, media := range post.Media {
		response.Media = append(response.Media, media.MediaURL)
	}

	return response, nil
}

// GetGroupPosts retrieves posts in a group
func (c *GroupController) GetGroupPosts(ctx context.Context, req *pb.GetGroupPostsRequest) (*pb.GetGroupPostsResponse, error) {
	// Get user ID from context or request
	userID := req.UserId
	if userID == "" {
		if id, ok := ctx.Value("userID").(string); ok {
			userID = id
		}
	}

	// Get posts
	posts, totalCount, totalPages, err := c.service.GetGroupPosts(ctx, req.GroupId, userID, int(req.Page), int(req.Limit))
	if err != nil {
		c.logger.Error("Failed to get group posts", err)
		return nil, status.Error(codes.Internal, "failed to get group posts")
	}

	// Create response
	response := &pb.GetGroupPostsResponse{
		Posts:      make([]*pb.GroupPostResponse, 0, len(posts)),
		TotalCount: int32(totalCount),
		Page:       req.Page,
		TotalPages: totalPages,
	}

	// Add posts to response
	for _, post := range posts {
		postResponse := &pb.GroupPostResponse{
			PostId:        post.ID,
			GroupId:       post.GroupID,
			AuthorId:      post.AuthorID,
			AuthorName:    "", // Would need to fetch from users service
			AuthorAvatar:  "", // Would need to fetch from users service
			Content:       post.Content,
			Media:         make([]string, 0, len(post.Media)),
			LikesCount:    int32(len(post.Likes)),
			CommentsCount: int32(len(post.Comments)),
			IsLiked:       false, // Would need to check if the user liked the post
			CreatedAt:     post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Add media to response
		for _, media := range post.Media {
			postResponse.Media = append(postResponse.Media, media.MediaURL)
		}

		// Check if the user liked the post
		for _, like := range post.Likes {
			if like.UserID == userID {
				postResponse.IsLiked = true
				break
			}
		}

		response.Posts = append(response.Posts, postResponse)
	}

	return response, nil
}
