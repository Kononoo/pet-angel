-- Pet Angel 数据库初始化脚本（优化版）
-- 变更要点：
-- 1) 聊天不再使用 会话+消息，两表；仅保留 messages（按 user_id 维度存储全量消息）
-- 2) 移除所有 JSON/ENUM/TIMESTAMP，统一使用：
--    - 整型 tinyint 表达枚举语义（业务层转义）
--    - 字符串/文本存储列表（如图片URL逗号分隔）、标签等
--    - DATETIME 统一记录时间（默认 CURRENT_TIMESTAMP）
-- 3) 资源全部保存为可访问的 URL 字符串（Minio 上的对象链接）

-- =========================
-- 用户表
-- =========================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id`           bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `nickname`     varchar(50)  DEFAULT NULL COMMENT '昵称',
  `username`     varchar(64)  DEFAULT '' COMMENT '登录用户名',
  `password`     varchar(255) DEFAULT '' COMMENT '密码哈希（bcrypt）',
  `avatar`       varchar(255) DEFAULT NULL COMMENT '用户头像URL',
  `model_id`     bigint(20)   NOT NULL COMMENT '当前宠物模型ID（关联 pet_models.id）',
  `model_url`    varchar(255) NOT NULL DEFAULT '/models/Dog_1.glb' COMMENT '当前宠物模型URL',
  `pet_name`     varchar(50)  DEFAULT NULL COMMENT '宠物名称',
  `pet_avatar`   varchar(255) DEFAULT NULL COMMENT '宠物头像URL',
  `pet_sex`      tinyint(1)   DEFAULT 0 COMMENT '宠物性别 0-未知 1-男 2-女',
  `kind`         varchar(50)  DEFAULT NULL COMMENT '宠物种类（中文/英文均可）',
  `weight`       int(11)      DEFAULT 0 COMMENT '宠物体重(kg)，整数保存',
  `hobby`        varchar(255) DEFAULT NULL COMMENT '宠物/用户爱好摘要',
  `description`  text         DEFAULT NULL COMMENT '个人/宠物简介',
  `coins`        int(11)      DEFAULT 0 COMMENT '金币余额',
  `created_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- =========================
-- 宠物模型表
-- =========================
DROP TABLE IF EXISTS `pet_models`;
CREATE TABLE `pet_models` (
  `id`         bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '模型ID',
  `name`       varchar(100) NOT NULL COMMENT '模型名称',
  `path`       varchar(255) NOT NULL COMMENT '模型资源URL（图片/三方存储）',
  `type`       tinyint(1)   NOT NULL COMMENT '宠物类型 0-猫 1-狗（业务层转义）',
  `is_default` tinyint(1)   DEFAULT 0 COMMENT '是否默认模型（每类可有一个默认）',
  `sort_order` int(11)      DEFAULT 0 COMMENT '排序序号（越小越靠前）',
  `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='宠物模型表';

-- =========================
-- 道具表
-- =========================
DROP TABLE IF EXISTS `items`;
CREATE TABLE `items` (
  `id`          bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '道具ID',
  `name`        varchar(100) NOT NULL COMMENT '道具名称',
  `description` varchar(255) DEFAULT NULL COMMENT '道具描述',
  `icon_path`   varchar(255) DEFAULT NULL COMMENT '道具图标URL',
  `coin_cost`   int(11)      DEFAULT 0 COMMENT '使用/解锁所需金币',
  `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='道具表';

-- =========================
-- 聊天：按用户维度的消息表（含小纸条）
-- =========================
-- sender: 0-用户 1-AI
-- message_type: 0-普通聊天 1-小纸条
-- 小纸条使用 is_locked/unlock_coins 控制解锁；普通消息两字段恒为 0
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id`           bigint(20)  NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `user_id`      bigint(20)  NOT NULL COMMENT '归属用户ID（每个用户一个“会话”概念）',
  `sender`       tinyint(1)  NOT NULL COMMENT '发送方 0-用户 1-AI',
  `message_type` tinyint(1)  NOT NULL DEFAULT 0 COMMENT '消息类型 0-聊天 1-小纸条',
  `is_locked`    tinyint(1)  NOT NULL DEFAULT 0 COMMENT '是否锁定（仅小纸条使用）',
  `unlock_coins` int(11)     NOT NULL DEFAULT 0 COMMENT '解锁所需金币（仅小纸条使用）',
  `content`      text         NOT NULL COMMENT '消息内容',
  `created_at`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_time` (`user_id`,`id`),
  KEY `idx_user_type` (`user_id`,`message_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表（用户与AI聊天+小纸条）';

-- 小纸条解锁记录
DROP TABLE IF EXISTS `user_unlock_records`;
CREATE TABLE `user_unlock_records` (
  `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `user_id`     bigint(20) NOT NULL COMMENT '用户ID',
  `message_id`  bigint(20) NOT NULL COMMENT '小纸条消息ID（message_type=1）',
  `coins_spent` int(11)    NOT NULL DEFAULT 0 COMMENT '消耗金币',
  `created_at`  datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_message` (`user_id`,`message_id`),
  KEY `idx_message_id` (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户小纸条解锁记录';

-- =========================
-- 社区：分类 + 帖子 + 评论 + 点赞 + 关注
-- =========================
DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id`         bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name`       varchar(50)  NOT NULL COMMENT '分类名称（如：日常/知识/信息/种草）',
  `sort_order` int(11)      DEFAULT 0 COMMENT '排序序号',
  `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子分类表';

-- 统一帖子表：type 0-图文 1-视频
-- 图文：image_urls 逗号分隔URL；视频：video_url + cover_url
DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id`            bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '帖子ID',
  `user_id`       bigint(20)   NOT NULL COMMENT '发布者ID',
  `category_id`   bigint(20)   DEFAULT NULL COMMENT '分类ID（可为空）',
  `title`         varchar(200) NOT NULL COMMENT '标题',
  `content`       text         DEFAULT NULL COMMENT '内容（富文本由前端渲染）',
  `type`          tinyint(1)   NOT NULL DEFAULT 0 COMMENT '帖子类型 0-图文 1-视频',
  `image_urls`    text         DEFAULT NULL COMMENT '图文：图片URL列表，逗号分隔（type=0）',
  `video_url`     varchar(255) DEFAULT NULL COMMENT '视频URL（type=1）',
  `cover_url`     varchar(255) DEFAULT NULL COMMENT '封面URL（视频/图文均可使用）',
  `locate`        varchar(100) DEFAULT NULL COMMENT '地理位置/文本',
  `tags`          varchar(255) DEFAULT NULL COMMENT '标签列表（英文逗号分隔）',
  `liked_count`   int(11)      NOT NULL DEFAULT 0 COMMENT '点赞数',
  `comment_count` int(11)      NOT NULL DEFAULT 0 COMMENT '评论数',
  `is_private`    tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否私密（仅自己可见）',
  `created_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_type_created` (`type`,`created_at`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_liked_count` (`liked_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子表（图文/视频统一）';

DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
  `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `post_id`     bigint(20) NOT NULL COMMENT '帖子ID',
  `user_id`     bigint(20) NOT NULL COMMENT '评论者ID',
  `content`     text       NOT NULL COMMENT '评论内容',
  `liked_count` int(11)    NOT NULL DEFAULT 0 COMMENT '点赞数',
  `created_at`  datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

-- 点赞表：target_type 0-帖子 1-评论
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes` (
  `user_id`     bigint(20) NOT NULL COMMENT '用户ID',
  `target_type` tinyint(1) NOT NULL COMMENT '目标类型 0-帖子 1-评论',
  `target_id`   bigint(20) NOT NULL COMMENT '目标ID',
  `created_at`  datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`user_id`,`target_type`,`target_id`),
  KEY `idx_target` (`target_type`,`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表（帖子/评论通用）';

-- 关注表（用户关系）
DROP TABLE IF EXISTS `user_follows`;
CREATE TABLE `user_follows` (
  `follower_id` bigint(20) NOT NULL COMMENT '关注者ID',
  `followee_id` bigint(20) NOT NULL COMMENT '被关注者ID',
  `created_at`  datetime   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`follower_id`,`followee_id`),
  KEY `idx_followee` (`followee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关注关系表';

-- =========================
-- 默认数据
-- =========================
INSERT INTO `pet_models` (`name`, `path`, `type`, `is_default`, `sort_order`) VALUES
('默认猫咪', '/models/cat/default.png', 0, 1, 1),
('可爱猫咪', '/models/cat/cute.png',    0, 0, 2),
('默认狗狗', '/models/dog/default.png', 1, 1, 1),
('忠诚狗狗', '/models/dog/loyal.png',   1, 0, 2)
ON DUPLICATE KEY UPDATE `path`=VALUES(`path`);

INSERT INTO `items` (`name`, `description`, `icon_path`, `coin_cost`) VALUES
('猫粮',  '让猫咪吃饱饱', '/items/cat_food.png', 5),
('狗粮',  '让狗狗吃饱饱', '/items/dog_food.png', 5),
('逗猫棒', '和猫咪玩耍', '/items/cat_wand.png', 10),
('飞盘',  '和狗狗玩耍', '/items/frisbee.png', 10),
('梳子',  '给宠物梳毛', '/items/comb.png', 8),
('项圈',  '给宠物戴上项圈', '/items/collar.png', 15)
ON DUPLICATE KEY UPDATE `icon_path`=VALUES(`icon_path`);

INSERT INTO `categories` (`name`, `sort_order`) VALUES
('日常', 1),('知识', 2),('信息', 3),('种草', 4)
ON DUPLICATE KEY UPDATE `sort_order`=VALUES(`sort_order`);