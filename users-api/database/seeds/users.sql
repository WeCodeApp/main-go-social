-- Seed data for users table
INSERT INTO users (id, name, email, avatar, provider, created_at, updated_at)
VALUES 
    ('user1', 'John Doe', 'john.doe@example.com', 'https://example.com/avatars/john.jpg', 'google', NOW(), NOW()),
    ('user2', 'Jane Smith', 'jane.smith@example.com', 'https://example.com/avatars/jane.jpg', 'microsoft', NOW(), NOW());