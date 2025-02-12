use todo;
CREATE TABLE community_posts (
                                 id INT AUTO_INCREMENT PRIMARY KEY,  -- 设置为自增
                                 user_id VARCHAR(255) NOT NULL,
                                 content TEXT NOT NULL,
                                 image_url VARCHAR(255),
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- 插入第一条数据
INSERT INTO community_posts (id, user_id, content, image_url, created_at, updated_at)
VALUES ('1', '12345', 'This is the first community post.', 'https://example.com/image1.jpg', NOW(), NOW());

-- 插入第二条数据
INSERT INTO community_posts (id, user_id, content, image_url, created_at, updated_at)
VALUES ('2', '67890', 'This is the second community post.', 'https://example.com/image2.jpg', NOW(), NOW());
