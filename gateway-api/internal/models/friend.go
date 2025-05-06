package models

// FriendRequest represents a friend request
type FriendRequest struct {
	FriendID string `json:"friend_id" binding:"required" example:"user456"`
}

// Friend represents a friend
type Friend struct {
	UserID       string `json:"user_id" example:"user456"`
	Name         string `json:"name" example:"Jane Doe"`
	Avatar       string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Email        string `json:"email" example:"jane.doe@example.com"`
	FriendsSince string `json:"friends_since" example:"2023-01-01T12:00:00Z"`
}

// FriendsResponse represents a list of friends with pagination
type FriendsResponse struct {
	Friends     []Friend `json:"friends"`
	TotalCount  int32    `json:"total_count" example:"42"`
	Page        int32    `json:"page" example:"1"`
	TotalPages  int32    `json:"total_pages" example:"5"`
}

// FriendRequestDetails represents a friend request with details
type FriendRequestDetails struct {
	RequestID      string `json:"request_id" example:"req123"`
	SenderID       string `json:"sender_id" example:"user123"`
	SenderName     string `json:"sender_name" example:"John Doe"`
	SenderAvatar   string `json:"sender_avatar" example:"https://example.com/avatar.jpg"`
	ReceiverID     string `json:"receiver_id" example:"user456"`
	ReceiverName   string `json:"receiver_name" example:"Jane Doe"`
	ReceiverAvatar string `json:"receiver_avatar" example:"https://example.com/avatar2.jpg"`
	Status         string `json:"status" example:"pending"`
	CreatedAt      string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt      string `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}

// FriendRequestsResponse represents a list of friend requests with pagination
type FriendRequestsResponse struct {
	Requests    []FriendRequestDetails `json:"requests"`
	TotalCount  int32                  `json:"total_count" example:"42"`
	Page        int32                  `json:"page" example:"1"`
	TotalPages  int32                  `json:"total_pages" example:"5"`
}