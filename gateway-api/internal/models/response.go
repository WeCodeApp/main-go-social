package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response with a boolean flag
type SuccessResponse struct {
	Success bool `json:"success"`
}

// SuccessWithCountResponse represents a success response with a count
type SuccessWithCountResponse struct {
	Success      bool `json:"success"`
	MembersCount int  `json:"members_count"`
}

// LikeResponse represents a response for like/unlike operations
type LikeResponse struct {
	Success    bool `json:"success"`
	LikesCount int  `json:"likes_count"`
}