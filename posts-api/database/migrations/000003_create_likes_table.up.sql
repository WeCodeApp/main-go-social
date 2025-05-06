CREATE TABLE IF NOT EXISTS likes (
    id VARCHAR(36) PRIMARY KEY,
    post_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_likes_post_id (post_id),
    INDEX idx_likes_user_id (user_id),
    INDEX idx_likes_deleted_at (deleted_at),
    UNIQUE INDEX idx_likes_post_user (post_id, user_id),
    CONSTRAINT fk_likes_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);