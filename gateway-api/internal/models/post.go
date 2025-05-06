package models

// PostCreateRequest represents a post creation request
type PostCreateRequest struct {
	Content    string   `json:"content" binding:"required" example:"This is a post"`
	Visibility string   `json:"visibility" binding:"required,oneof=public private" example:"public"`
	Media      []string `json:"media,omitempty" example:"[\"https://example.com/image1.jpg\"]"`
}

// PostUpdateRequest represents a post update request
type PostUpdateRequest struct {
	Content    string   `json:"content,omitempty" example:"This is an updated post"`
	Visibility string   `json:"visibility,omitempty" binding:"omitempty,oneof=public private" example:"private"`
	Media      []string `json:"media,omitempty" example:"[\"https://example.com/image2.jpg\"]"`
}

// Post represents a post
type Post struct {
	PostID        string   `json:"post_id" example:"post123"`
	AuthorID      string   `json:"author_id" example:"user123"`
	AuthorName    string   `json:"author_name" example:"John Doe"`
	AuthorAvatar  string   `json:"author_avatar" example:"https://example.com/avatar.jpg"`
	Content       string   `json:"content" example:"This is a post"`
	Media         []string `json:"media" example:"[\"https://example.com/image1.jpg\"]"`
	Visibility    string   `json:"visibility" example:"public"`
	LikesCount    int32    `json:"likes_count" example:"42"`
	CommentsCount int32    `json:"comments_count" example:"10"`
	IsLiked       bool     `json:"is_liked" example:"false"`
	CreatedAt     string   `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt     string   `json:"updated_at" example:"2023-01-02T12:00:00Z"`
}

// PostsResponse represents a list of posts with pagination
type PostsResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int32  `json:"total_count" example:"42"`
	Page       int32  `json:"page" example:"1"`
	TotalPages int32  `json:"total_pages" example:"5"`
}

// CommentCreateRequest represents a comment creation request
type CommentCreateRequest struct {
	Content string `json:"content" binding:"required" example:"This is a comment"`
}

// Comment represents a comment
type Comment struct {
	CommentID    string `json:"comment_id" example:"comment123"`
	PostID       string `json:"post_id" example:"post123"`
	AuthorID     string `json:"author_id" example:"user123"`
	AuthorName   string `json:"author_name" example:"John Doe"`
	AuthorAvatar string `json:"author_avatar" example:"https://example.com/avatar.jpg"`
	Content      string `json:"content" example:"This is a comment"`
	CreatedAt    string `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt    string `json:"updated_at" example:"2023-01-02T12:00:00Z"`
}

// CommentsResponse represents a list of comments with pagination
type CommentsResponse struct {
	Comments   []Comment `json:"comments"`
	TotalCount int32     `json:"total_count" example:"42"`
	Page       int32     `json:"page" example:"1"`
	TotalPages int32     `json:"total_pages" example:"5"`
}