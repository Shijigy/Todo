-- 添加打卡任务表
-- 创建打卡任务表
use todo;
CREATE TABLE IF NOT EXISTS todos (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY,                     -- 打卡任务ID
                                     user_id BIGINT NOT NULL,                                  -- 用户ID
                                     title VARCHAR(255) NOT NULL,                               -- 打卡任务名称
                                     description TEXT,                                          -- 任务描述 (新增字段)
                                     status ENUM('pending', 'completed', 'failed') DEFAULT 'pending',  -- 打卡状态
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,            -- 创建时间
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                                     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  -- 外键约束
);

-- 创建打卡历史记录表
CREATE TABLE IF NOT EXISTS todos_history (
                                             id BIGINT AUTO_INCREMENT PRIMARY KEY,                          -- 历史记录ID
                                             todo_id BIGINT NOT NULL,                                        -- 打卡任务ID
                                             checkin_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,               -- 打卡时间
                                             status ENUM('pending', 'completed', 'failed') DEFAULT 'pending', -- 打卡状态
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                 -- 创建时间
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- 更新时间
                                             FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE   -- 外键约束
);


-- 插入一些初始打卡任务
INSERT INTO todos (user_id, title, description, status, created_at, updated_at) VALUES
                                                                                    (1, '晨间打卡', '完成早晨的任务打卡', 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
                                                                                    (2, '工作结束打卡', '记录下班时间打卡', 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
