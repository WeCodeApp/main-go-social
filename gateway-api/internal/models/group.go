package models

// GroupCreateRequest represents a group creation request
type GroupCreateRequest struct {
	Name        string `json:"name" binding:"required" example:"Tech Enthusiasts"`
	Description string `json:"description" example:"A group for tech enthusiasts"`
	Avatar      string `json:"avatar" example:"https://example.com/group-avatar.jpg"`
}

// GroupUpdateRequest represents a group update request
type GroupUpdateRequest struct {
	Name        string `json:"name,omitempty" example:"Tech Enthusiasts Updated"`
	Description string `json:"description,omitempty" example:"An updated group for tech enthusiasts"`
	Avatar      string `json:"avatar,omitempty" example:"https://example.com/updated-group-avatar.jpg"`
}

// Group represents a group
type Group struct {
	GroupID      string `json:"group_id" example:"group123"`
	Name         string `json:"name" example:"Tech Enthusiasts"`
	Description  string `json:"description" example:"A group for tech enthusiasts"`
	Avatar       string `json:"avatar" example:"https://example.com/group-avatar.jpg"`
	CreatorID    string `json:"creator_id" example:"user123"`
	CreatorName  string `json:"creator_name" example:"John Doe"`
	MembersCount int32  `json:"members_count" example:"42"`
	CreatedAt    string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt    string `json:"updated_at" example:"2023-01-02T12:00:00Z"`
}

// GroupsResponse represents a list of groups with pagination
type GroupsResponse struct {
	Groups     []Group `json:"groups"`
	TotalCount int32   `json:"total_count" example:"42"`
	Page       int32   `json:"page" example:"1"`
	TotalPages int32   `json:"total_pages" example:"5"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	UserID   string `json:"user_id" example:"user123"`
	Name     string `json:"name" example:"John Doe"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Role     string `json:"role" example:"admin"`
	JoinedAt string `json:"joined_at" example:"2023-01-01T12:00:00Z"`
}

// GroupMembersResponse represents a list of group members with pagination
type GroupMembersResponse struct {
	Members    []GroupMember `json:"members"`
	TotalCount int32         `json:"total_count" example:"42"`
	Page       int32         `json:"page" example:"1"`
	TotalPages int32         `json:"total_pages" example:"5"`
}

// GroupPostRequest represents a group post creation request
type GroupPostRequest struct {
	Content string   `json:"content" binding:"required" example:"This is a post in the group"`
	Media   []string `json:"media,omitempty" example:"[\"https://example.com/image1.jpg\"]"`
}