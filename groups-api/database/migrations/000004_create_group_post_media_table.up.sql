CREATE TABLE IF NOT EXISTS `group_post_media` (
    id VARCHAR(36) PRIMARY KEY,
    post_id VARCHAR(36) NOT NULL,
    media_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (post_id) REFERENCES `group_posts`(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_post_media_post_id ON `group_post_media`(post_id);
