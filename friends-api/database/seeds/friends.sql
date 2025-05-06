-- Seed data for friend_requests table
INSERT INTO friend_requests (id, sender_id, receiver_id, status, created_at, updated_at)
VALUES
    ('fr-001', 'user-001', 'user-002', 'pending', NOW(), NOW()),
    ('fr-002', 'user-003', 'user-001', 'pending', NOW(), NOW()),
    ('fr-003', 'user-004', 'user-001', 'accepted', NOW(), NOW()),
    ('fr-004', 'user-001', 'user-005', 'accepted', NOW(), NOW()),
    ('fr-005', 'user-006', 'user-001', 'rejected', NOW(), NOW()),
    ('fr-006', 'user-001', 'user-007', 'rejected', NOW(), NOW());

-- Seed data for friendships table
INSERT INTO friendships (id, user_id, friend_id, created_at, updated_at)
VALUES
    ('f-001', 'user-001', 'user-004', NOW(), NOW()),
    ('f-002', 'user-004', 'user-001', NOW(), NOW()),
    ('f-003', 'user-001', 'user-005', NOW(), NOW()),
    ('f-004', 'user-005', 'user-001', NOW(), NOW());

-- Seed data for blocked_users table
INSERT INTO blocked_users (id, user_id, blocked_user_id, created_at, updated_at)
VALUES
    ('bu-001', 'user-001', 'user-008', NOW(), NOW()),
    ('bu-002', 'user-009', 'user-001', NOW(), NOW());