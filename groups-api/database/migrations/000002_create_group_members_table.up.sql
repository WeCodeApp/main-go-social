CREATE TABLE IF NOT EXISTS `group_members` (
    id VARCHAR(36) PRIMARY KEY,
    group_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    role ENUM('creator', 'admin', 'member') NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE
);

CREATE INDEX idx_group_members_group_id ON `group_members`(group_id);
CREATE INDEX idx_group_members_user_id ON `group_members`(user_id);
CREATE UNIQUE INDEX idx_group_members_group_user ON `group_members`(group_id, user_id, deleted_at);
