CREATE TABLE community_posts (
                                 id INT AUTO_INCREMENT PRIMARY KEY,       -- 自增主键
                                 user_id VARCHAR(255) NOT NULL,           -- 用户 ID
                                 content TEXT NOT NULL,                   -- 动态内容
                                 image_url VARCHAR(255),                  -- 动态附带图片的 URL
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                                 tags VARCHAR(255),                       -- 标签字段
                                 likes_count INT DEFAULT 0               -- 点赞数，默认为 0
);

CREATE TABLE likes (
                       id INT AUTO_INCREMENT PRIMARY KEY,       -- 点赞记录的自增 ID
                       user_id VARCHAR(255) NOT NULL,           -- 点赞的用户 ID
                       post_id INT NOT NULL,                    -- 点赞的社区动态 ID
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 点赞时间
);

-- 插入第一条社区动态
INSERT INTO community_posts (user_id, content, image_url, created_at, updated_at, tags, likes_count)
VALUES ('12345', 'This is the first community post.', 'https://example.com/image1.jpg', NOW(), NOW(), 'tag1, tag2', 0);

-- 插入第二条社区动态
INSERT INTO community_posts (user_id, content, image_url, created_at, updated_at, tags, likes_count)
VALUES ('67890', 'This is the second community post.', 'https://example.com/image2.jpg', NOW(), NOW(), 'tag3, tag4', 0);

-- 用户12345对帖子1点赞
INSERT INTO likes (user_id, post_id, created_at)
VALUES ('12345', 1, NOW());

-- 用户67890对帖子2点赞
INSERT INTO likes (user_id, post_id, created_at)
VALUES ('67890', 2, NOW());
