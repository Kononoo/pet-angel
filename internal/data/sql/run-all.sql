-- Pet Angel 数据库初始化与数据填充脚本
-- 按顺序执行所有SQL文件

-- 1. 初始化表结构
SOURCE 1-init-tables.sql;

-- 2. 基础种子数据
SOURCE 2-seed-data.sql;

-- 3. 海量数据生成
SOURCE 3-massive-data.sql;

-- 4. 补充数据
SOURCE 4-more-data.sql;

-- 完成提示
SELECT 'Pet Angel 数据库初始化完成！' AS message;
SELECT COUNT(*) AS total_users FROM users;
SELECT COUNT(*) AS total_posts FROM posts;
SELECT COUNT(*) AS total_messages FROM messages;
SELECT COUNT(*) AS total_comments FROM comments;
SELECT COUNT(*) AS total_likes FROM likes; 