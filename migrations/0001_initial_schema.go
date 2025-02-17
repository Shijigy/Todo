CREATE TABLE `users` (
                         `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,  -- 用户 ID，自增整数类型，作为主键
                         `username` VARCHAR(255) NOT NULL,               -- 用户名
                         `password` VARCHAR(255) NOT NULL,               -- 密码
                         `email` VARCHAR(255) NOT NULL UNIQUE,           -- 邮箱，确保唯一
                         `created_at` DATETIME NOT NULL,                 -- 创建时间
                         `updated_at` DATETIME NOT NULL,                 -- 更新时间
                         `status` INT NOT NULL DEFAULT 1                  -- 状态，默认为 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
