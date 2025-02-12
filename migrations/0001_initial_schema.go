-- 初始数据库架构
-- 创建用户表、角色表等核心表

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY,   -- 用户ID
                                     username VARCHAR(255) NOT NULL UNIQUE,   -- 用户名
                                     email VARCHAR(255) NOT NULL UNIQUE,      -- 用户邮箱
                                     password VARCHAR(255) NOT NULL,     -- 密码哈希
                                     first_name VARCHAR(100),                 -- 姓名
                                     last_name VARCHAR(100),                  -- 姓氏
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
);

-- 创建角色表
CREATE TABLE IF NOT EXISTS roles (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY,  -- 角色ID
                                     name VARCHAR(50) NOT NULL UNIQUE,      -- 角色名称
                                     description TEXT,                      -- 角色描述
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- 更新时间
);

-- 创建用户-角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
                                          user_id BIGINT NOT NULL,               -- 用户ID
                                          role_id BIGINT NOT NULL,               -- 角色ID
                                          PRIMARY KEY (user_id, role_id),        -- 复合主键
                                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,  -- 外键约束
                                          FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE   -- 外键约束
);

-- 创建权限表
CREATE TABLE IF NOT EXISTS permissions (
                                           id BIGINT AUTO_INCREMENT PRIMARY KEY,  -- 权限ID
                                           name VARCHAR(100) NOT NULL UNIQUE,     -- 权限名称
                                           description TEXT,                      -- 权限描述
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- 更新时间
);

-- 创建角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
                                                role_id BIGINT NOT NULL,               -- 角色ID
                                                permission_id BIGINT NOT NULL,         -- 权限ID
                                                PRIMARY KEY (role_id, permission_id),  -- 复合主键
                                                FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,  -- 外键约束
                                                FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE -- 外键约束
);

-- 创建用户最后一次登录记录表
CREATE TABLE IF NOT EXISTS user_logins (
                                           id BIGINT AUTO_INCREMENT PRIMARY KEY,   -- 登录记录ID
                                           user_id BIGINT NOT NULL,                -- 用户ID
                                           login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 登录时间
                                           ip_address VARCHAR(50),                 -- 用户登录IP地址
                                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  -- 外键约束
);

-- 示例数据插入
-- 插入一些初始角色
INSERT INTO roles (name, description) VALUES
                                          ('admin', '管理员，拥有所有权限'),
                                          ('user', '普通用户，有限权限');

-- 插入一些初始权限
INSERT INTO permissions (name, description) VALUES
                                                ('read', '读取权限'),
                                                ('write', '写入权限'),
                                                ('delete', '删除权限');

-- 插入用户数据
INSERT INTO users (username, email, password, first_name, last_name) VALUES
                                                ('john_doe', 'john.doe@example.com', 'hashed_password_123', 'John', 'Doe'),
                                                ('jane_smith', 'jane.smith@example.com', 'hashed_password_456', 'Jane', 'Smith');

