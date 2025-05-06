package services

import (
	"context"
	"errors"
	"groups-api/internal/models"
	"groups-api/internal/repository"
	"groups-api/internal/utils/logger"
)

// GroupService defines the interface for group-related business logic
type GroupService interface {
	// Group operations
	CreateGroup(ctx context.Context, userID, name, description, avatar string) (*models.Group, error)
	GetGroup(ctx context.Context, id string, userID string) (*models.Group, int32, int32, bool, error)
	GetGroups(ctx context.Context, userID, query string, page, limit int) ([]*models.Group, int64, int32, error)
	UpdateGroup(ctx context.Context, id, userID, name, description, avatar string) (*models.Group, error)
	DeleteGroup(ctx context.Context, id, userID string) error

	// Group member operations
	JoinGroup(ctx context.Context, groupID, userID string) (bool, int32, error)
	LeaveGroup(ctx context.Context, groupID, userID string) (bool, int32, error)
	GetGroupMembers(ctx context.Context, groupID string, page, limit int) ([]*models.GroupMember, int64, int32, error)

	// Group post operations
	CreateGroupPost(ctx context.Context, groupID, userID, content string, mediaURLs []string) (*models.GroupPost, error)
	GetGroupPosts(ctx context.Context, groupID, userID string, page, limit int) ([]*models.GroupPost, int64, int32, error)
}

// groupService implements the GroupService interface
type groupService struct {
	repo   repository.GroupRepository
	logger *logger.Logger
}

// NewGroupService creates a new group service
func NewGroupService(repo repository.GroupRepository, logger *logger.Logger) GroupService {
	return &groupService{
		repo:   repo,
		logger: logger,
	}
}

// CreateGroup creates a new group
func (s *groupService) CreateGroup(ctx context.Context, userID, name, description, avatar string) (*models.Group, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if name == "" {
		return nil, errors.New("group name is required")
	}

	// Create group
	group := &models.Group{
		Name:        name,
		Description: description,
		Avatar:      avatar,
		CreatorID:   userID,
	}

	// Save group to database
	err := s.repo.CreateGroup(ctx, group)
	if err != nil {
		s.logger.Error("Failed to create group", err)
		return nil, err
	}

	// Add creator as a member with creator role
	member := &models.GroupMember{
		GroupID: group.ID,
		UserID:  userID,
		Role:    "creator",
	}

	err = s.repo.AddMember(ctx, member)
	if err != nil {
		s.logger.Error("Failed to add creator as member", err)
		// Don't return error here, as the group was created successfully
	}

	return group, nil
}

// GetGroup gets a group by ID
func (s *groupService) GetGroup(ctx context.Context, id string, userID string) (*models.Group, int32, int32, bool, error) {
	// Get group from database
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, 0, 0, false, err
	}

	// Get member count
	_, count, err := s.repo.GetGroupMembers(ctx, id, 1, 1)
	if err != nil {
		s.logger.Error("Failed to get group members", err)
		// Don't return error here, as we can still return the group
	}

	// Get post count
	_, postCount, err := s.repo.GetGroupPosts(ctx, id, 1, 1)
	if err != nil {
		s.logger.Error("Failed to get group posts", err)
		// Don't return error here, as we can still return the group
	}

	// Check if user is a member
	isMember := false
	if userID != "" {
		isMember, err = s.repo.IsMember(ctx, id, userID)
		if err != nil {
			s.logger.Error("Failed to check if user is a member", err)
			// Don't return error here, as we can still return the group
		}
	}

	return group, int32(count), int32(postCount), isMember, nil
}

// GetGroups gets groups with pagination and filtering
func (s *groupService) GetGroups(ctx context.Context, userID, query string, page, limit int) ([]*models.Group, int64, int32, error) {
	// Get groups from database
	groups, count, err := s.repo.GetGroups(ctx, query, page, limit)
	if err != nil {
		s.logger.Error("Failed to get groups", err)
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return groups, count, totalPages, nil
}

// UpdateGroup updates a group
func (s *groupService) UpdateGroup(ctx context.Context, id, userID, name, description, avatar string) (*models.Group, error) {
	// Get group from database
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, err
	}

	// Check if user is the creator or an admin
	member, err := s.repo.GetMemberByID(ctx, id, userID)
	if err != nil {
		s.logger.Error("Failed to get member", err)
		return nil, errors.New("not authorized to update this group")
	}

	if member.Role != "creator" && member.Role != "admin" {
		return nil, errors.New("not authorized to update this group")
	}

	// Update group
	if name != "" {
		group.Name = name
	}
	if description != "" {
		group.Description = description
	}
	if avatar != "" {
		group.Avatar = avatar
	}

	// Save group to database
	err = s.repo.UpdateGroup(ctx, group)
	if err != nil {
		s.logger.Error("Failed to update group", err)
		return nil, err
	}

	return group, nil
}

// DeleteGroup deletes a group
func (s *groupService) DeleteGroup(ctx context.Context, id, userID string) error {
	// Get group from database
	group, err := s.repo.GetGroupByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return err
	}

	// Check if user is the creator
	if group.CreatorID != userID {
		return errors.New("not authorized to delete this group")
	}

	// Delete group from database
	err = s.repo.DeleteGroup(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete group", err)
		return err
	}

	return nil
}

// JoinGroup adds a user to a group
func (s *groupService) JoinGroup(ctx context.Context, groupID, userID string) (bool, int32, error) {
	// Check if group exists
	_, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return false, 0, err
	}

	// Check if user is already a member
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check if user is a member", err)
		return false, 0, err
	}

	if isMember {
		return false, 0, errors.New("already a member of this group")
	}

	// Add user as a member
	member := &models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    "member",
	}

	err = s.repo.AddMember(ctx, member)
	if err != nil {
		s.logger.Error("Failed to add member", err)
		return false, 0, err
	}

	// Get updated member count
	_, count, err := s.repo.GetGroupMembers(ctx, groupID, 1, 1)
	if err != nil {
		s.logger.Error("Failed to get group members", err)
		// Don't return error here, as the user was added successfully
		return true, 0, nil
	}

	return true, int32(count), nil
}

// LeaveGroup removes a user from a group
func (s *groupService) LeaveGroup(ctx context.Context, groupID, userID string) (bool, int32, error) {
	// Check if group exists
	group, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return false, 0, err
	}

	// Check if user is a member
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check if user is a member", err)
		return false, 0, err
	}

	if !isMember {
		return false, 0, errors.New("not a member of this group")
	}

	// Check if user is the creator
	if group.CreatorID == userID {
		return false, 0, errors.New("creator cannot leave the group")
	}

	// Remove user from group
	err = s.repo.RemoveMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Error("Failed to remove member", err)
		return false, 0, err
	}

	// Get updated member count
	_, count, err := s.repo.GetGroupMembers(ctx, groupID, 1, 1)
	if err != nil {
		s.logger.Error("Failed to get group members", err)
		// Don't return error here, as the user was removed successfully
		return true, 0, nil
	}

	return true, int32(count), nil
}

// GetGroupMembers gets members of a group with pagination
func (s *groupService) GetGroupMembers(ctx context.Context, groupID string, page, limit int) ([]*models.GroupMember, int64, int32, error) {
	// Check if group exists
	_, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, 0, 0, err
	}

	// Get members from database
	members, count, err := s.repo.GetGroupMembers(ctx, groupID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get group members", err)
		return nil, 0, 0, err
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return members, count, totalPages, nil
}

// CreateGroupPost creates a new post in a group
func (s *groupService) CreateGroupPost(ctx context.Context, groupID, userID, content string, mediaURLs []string) (*models.GroupPost, error) {
	// Check if group exists
	_, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, err
	}

	// Check if user is a member
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		s.logger.Error("Failed to check if user is a member", err)
		return nil, err
	}

	if !isMember {
		return nil, errors.New("not a member of this group")
	}

	// Create post
	post := &models.GroupPost{
		GroupID:  groupID,
		AuthorID: userID,
		Content:  content,
	}

	// Save post to database
	err = s.repo.CreatePost(ctx, post)
	if err != nil {
		s.logger.Error("Failed to create post", err)
		return nil, err
	}

	// Add media to post
	for _, mediaURL := range mediaURLs {
		media := &models.GroupPostMedia{
			PostID:   post.ID,
			MediaURL: mediaURL,
		}

		err = s.repo.AddPostMedia(ctx, media)
		if err != nil {
			s.logger.Error("Failed to add media to post", err)
			// Don't return error here, as the post was created successfully
		}
	}

	// Get media for post
	media, err := s.repo.GetPostMedia(ctx, post.ID)
	if err != nil {
		s.logger.Error("Failed to get post media", err)
		// Don't return error here, as the post was created successfully
	} else {
		post.Media = media
	}

	return post, nil
}

// GetGroupPosts gets posts in a group with pagination
func (s *groupService) GetGroupPosts(ctx context.Context, groupID, userID string, page, limit int) ([]*models.GroupPost, int64, int32, error) {
	// Check if group exists
	_, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		s.logger.Error("Failed to get group", err)
		return nil, 0, 0, err
	}

	// Check if user is a member
	isMember := false
	if userID != "" {
		isMember, err = s.repo.IsMember(ctx, groupID, userID)
		if err != nil {
			s.logger.Error("Failed to check if user is a member", err)
			return nil, 0, 0, err
		}
	}

	// Only members can see posts
	if !isMember {
		return nil, 0, 0, errors.New("not a member of this group")
	}

	// Get posts from database
	posts, count, err := s.repo.GetGroupPosts(ctx, groupID, page, limit)
	if err != nil {
		s.logger.Error("Failed to get group posts", err)
		return nil, 0, 0, err
	}

	// Get media, likes, and comments for each post
	for _, post := range posts {
		// Get media
		media, err := s.repo.GetPostMedia(ctx, post.ID)
		if err != nil {
			s.logger.Error("Failed to get post media", err)
			// Don't return error here, as we can still return the posts
		} else {
			post.Media = media
		}

		// Get likes
		likes, err := s.repo.GetPostLikes(ctx, post.ID)
		if err != nil {
			s.logger.Error("Failed to get post likes", err)
			// Don't return error here, as we can still return the posts
		} else {
			post.Likes = likes
		}

		// Get comments
		comments, _, err := s.repo.GetPostComments(ctx, post.ID, 1, 100)
		if err != nil {
			s.logger.Error("Failed to get post comments", err)
			// Don't return error here, as we can still return the posts
		} else {
			post.Comments = comments
		}
	}

	// Calculate total pages
	totalPages := int32((count + int64(limit) - 1) / int64(limit))

	return posts, count, totalPages, nil
}
