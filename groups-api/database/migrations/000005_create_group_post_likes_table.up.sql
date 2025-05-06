CREATE TABLE IF NOT EXISTS `group_post_likes` (
    id VARCHAR(36) PRIMARY KEY,
    post_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (post_id) REFERENCES `group_posts`(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_post_likes_post_id ON `group_post_likes`(post_id);
CREATE INDEX idx_group_post_likes_user_id ON `group_post_likes`(user_id);
CREATE UNIQUE INDEX idx_group_post_likes_post_user ON `group_post_likes`(post_id, user_id, deleted_at);
