package models

// AuthRequest represents an authentication request
type AuthRequest struct {
	Provider    string `json:"provider" binding:"required,oneof=google microsoft" example:"google"`
	AccessToken string `json:"access_token" binding:"required" example:"ya29.a0AfB_byC..."`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	UserID      string `json:"user_id" example:"user123"`
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ProfileUpdateRequest represents a profile update request
type ProfileUpdateRequest struct {
	Name   string `json:"name" example:"John Doe"`
	Avatar string `json:"avatar" example:"https://example.com/avatar.jpg"`
}

// UserProfile represents a user profile
type UserProfile struct {
	UserID    string `json:"user_id" example:"user123"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Avatar    string `json:"avatar" example:"https://example.com/avatar.jpg"`
	CreatedAt string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-02T12:00:00Z"`
}