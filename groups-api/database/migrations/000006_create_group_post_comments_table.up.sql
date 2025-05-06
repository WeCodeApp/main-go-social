CREATE TABLE IF NOT EXISTS `group_post_comments` (
    id VARCHAR(36) PRIMARY KEY,
    post_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (post_id) REFERENCES `group_posts`(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_post_comments_post_id ON `group_post_comments`(post_id);
CREATE INDEX idx_group_post_comments_user_id ON `group_post_comments`(user_id);
