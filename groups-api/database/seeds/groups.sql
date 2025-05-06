-- Seed data for groups table
INSERT INTO groups (id, name, description, avatar, creator_id, created_at, updated_at)
VALUES
    ('group-001', 'Tech Enthusiasts', 'A group for tech lovers', 'https://example.com/avatars/tech.jpg', 'user-001', NOW(), NOW()),
    ('group-002', 'Book Club', 'Discuss your favorite books', 'https://example.com/avatars/books.jpg', 'user-002', NOW(), NOW()),
    ('group-003', 'Fitness Fanatics', 'Share your fitness journey', 'https://example.com/avatars/fitness.jpg', 'user-003', NOW(), NOW());

-- Seed data for group_members table
INSERT INTO group_members (id, group_id, user_id, role, joined_at, created_at, updated_at)
VALUES
    ('gm-001', 'group-001', 'user-001', 'creator', NOW(), NOW(), NOW()),
    ('gm-002', 'group-001', 'user-002', 'member', NOW(), NOW(), NOW()),
    ('gm-003', 'group-001', 'user-003', 'member', NOW(), NOW(), NOW()),
    ('gm-004', 'group-002', 'user-002', 'creator', NOW(), NOW(), NOW()),
    ('gm-005', 'group-002', 'user-001', 'member', NOW(), NOW(), NOW()),
    ('gm-006', 'group-003', 'user-003', 'creator', NOW(), NOW(), NOW()),
    ('gm-007', 'group-003', 'user-004', 'admin', NOW(), NOW(), NOW()),
    ('gm-008', 'group-003', 'user-005', 'member', NOW(), NOW(), NOW());

-- Seed data for group_posts table
INSERT INTO group_posts (id, group_id, author_id, content, created_at, updated_at)
VALUES
    ('gp-001', 'group-001', 'user-001', 'Welcome to the Tech Enthusiasts group!', NOW(), NOW()),
    ('gp-002', 'group-001', 'user-002', 'Has anyone tried the new iPhone?', NOW(), NOW()),
    ('gp-003', 'group-002', 'user-002', 'Welcome to the Book Club!', NOW(), NOW()),
    ('gp-004', 'group-002', 'user-001', 'I just finished reading "The Great Gatsby"', NOW(), NOW()),
    ('gp-005', 'group-003', 'user-003', 'Welcome to Fitness Fanatics!', NOW(), NOW()),
    ('gp-006', 'group-003', 'user-004', 'My workout routine for this week', NOW(), NOW());

-- Seed data for group_post_media table
INSERT INTO group_post_media (id, post_id, media_url, created_at)
VALUES
    ('gpm-001', 'gp-002', 'https://example.com/media/iphone.jpg', NOW()),
    ('gpm-002', 'gp-004', 'https://example.com/media/gatsby.jpg', NOW()),
    ('gpm-003', 'gp-006', 'https://example.com/media/workout.jpg', NOW()),
    ('gpm-004', 'gp-006', 'https://example.com/media/workout2.jpg', NOW());

-- Seed data for group_post_likes table
INSERT INTO group_post_likes (id, post_id, user_id, created_at)
VALUES
    ('gpl-001', 'gp-001', 'user-002', NOW()),
    ('gpl-002', 'gp-001', 'user-003', NOW()),
    ('gpl-003', 'gp-002', 'user-001', NOW()),
    ('gpl-004', 'gp-003', 'user-001', NOW()),
    ('gpl-005', 'gp-005', 'user-004', NOW()),
    ('gpl-006', 'gp-005', 'user-005', NOW()),
    ('gpl-007', 'gp-006', 'user-003', NOW());

-- Seed data for group_post_comments table
INSERT INTO group_post_comments (id, post_id, user_id, content, created_at, updated_at)
VALUES
    ('gpc-001', 'gp-001', 'user-002', 'Excited to be here!', NOW(), NOW()),
    ('gpc-002', 'gp-002', 'user-001', 'Yes, it''s amazing!', NOW(), NOW()),
    ('gpc-003', 'gp-002', 'user-003', 'I prefer Android, but the new iPhone looks good.', NOW(), NOW()),
    ('gpc-004', 'gp-003', 'user-001', 'Looking forward to discussing books!', NOW(), NOW()),
    ('gpc-005', 'gp-004', 'user-002', 'That''s a classic! What did you think of it?', NOW(), NOW()),
    ('gpc-006', 'gp-005', 'user-004', 'Let''s get fit together!', NOW(), NOW()),
    ('gpc-007', 'gp-006', 'user-005', 'Thanks for sharing your routine!', NOW(), NOW());