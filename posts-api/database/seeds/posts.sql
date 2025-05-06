-- Sample data for posts
INSERT INTO posts (id, author_id, author_name, author_avatar, content, visibility, group_id, group_name, media, likes_count, comments_count)
VALUES
    ('post1', 'user1', 'John Doe', 'https://example.com/avatars/john.jpg', 'This is a public post by John', 'public', NULL, NULL, '["https://example.com/images/sample1.jpg"]', 2, 1),
    ('post2', 'user1', 'John Doe', 'https://example.com/avatars/john.jpg', 'This is a private post by John', 'private', NULL, NULL, NULL, 0, 0),
    ('post3', 'user2', 'Jane Smith', 'https://example.com/avatars/jane.jpg', 'This is a public post by Jane', 'public', NULL, NULL, NULL, 1, 2),
    ('post4', 'user2', 'Jane Smith', 'https://example.com/avatars/jane.jpg', 'This is a post in Group A', 'public', 'group1', 'Group A', NULL, 0, 0),
    ('post5', 'user3', 'Bob Johnson', 'https://example.com/avatars/bob.jpg', 'This is a post in Group B', 'public', 'group2', 'Group B', '["https://example.com/images/sample2.jpg", "https://example.com/images/sample3.jpg"]', 1, 1);

-- Sample data for comments
INSERT INTO comments (id, post_id, author_id, author_name, author_avatar, content)
VALUES
    ('comment1', 'post1', 'user2', 'Jane Smith', 'https://example.com/avatars/jane.jpg', 'Great post, John!'),
    ('comment2', 'post3', 'user1', 'John Doe', 'https://example.com/avatars/john.jpg', 'Thanks for sharing, Jane!'),
    ('comment3', 'post3', 'user3', 'Bob Johnson', 'https://example.com/avatars/bob.jpg', 'I agree with John!'),
    ('comment4', 'post5', 'user1', 'John Doe', 'https://example.com/avatars/john.jpg', 'Nice photos, Bob!');

-- Sample data for likes
INSERT INTO likes (id, post_id, user_id)
VALUES
    ('like1', 'post1', 'user2'),
    ('like2', 'post1', 'user3'),
    ('like3', 'post3', 'user1'),
    ('like4', 'post5', 'user2');