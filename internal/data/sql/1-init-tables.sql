-- Pet Angel 数据库初始化脚本

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `gender` enum('male','female','other') DEFAULT NULL COMMENT '性别',
  `region` varchar(100) DEFAULT NULL COMMENT '地区',
  `partner` varchar(100) DEFAULT NULL COMMENT '伴侣',
  `coin_balance` int(11) DEFAULT 0 COMMENT '金币余额',
  `total_earned` int(11) DEFAULT 0 COMMENT '总获得金币',
  `total_spent` int(11) DEFAULT 0 COMMENT '总消费金币',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 宠物表
CREATE TABLE IF NOT EXISTS `pets` (
  `pet_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '主人ID',
  `name` varchar(50) NOT NULL COMMENT '宠物名称',
  `species` varchar(50) NOT NULL COMMENT '物种',
  `gender` enum('male','female') DEFAULT NULL COMMENT '性别',
  `weight` decimal(5,2) DEFAULT NULL COMMENT '体重(kg)',
  `hobbies` json DEFAULT NULL COMMENT '爱好列表',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `adoption_date` date DEFAULT NULL COMMENT '收养日期',
  `passed_away_date` date DEFAULT NULL COMMENT '去世日期',
  `memorial_words` text DEFAULT NULL COMMENT '纪念语',
  `avatar_id` bigint(20) DEFAULT NULL COMMENT '关联虚拟形象ID',
  `background_image` varchar(255) DEFAULT NULL COMMENT '背景图片',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`pet_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_avatar_id` (`avatar_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='宠物表';

-- 虚拟形象表
CREATE TABLE IF NOT EXISTS `avatars` (
  `avatar_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '形象名称',
  `resource_path` varchar(255) NOT NULL COMMENT '资源路径',
  `idle_animation` varchar(255) DEFAULT NULL COMMENT '待机动画',
  `switch_animation` varchar(255) DEFAULT NULL COMMENT '切换动画',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认',
  `is_preset` tinyint(1) DEFAULT 1 COMMENT '是否预设',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`avatar_id`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='虚拟形象表';

-- 用户当前形象表
CREATE TABLE IF NOT EXISTS `user_avatars` (
  `user_id` bigint(20) NOT NULL,
  `avatar_id` bigint(20) NOT NULL,
  `last_switch_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  KEY `idx_avatar_id` (`avatar_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户当前形象表';

-- 道具分类表
CREATE TABLE IF NOT EXISTS `prop_categories` (
  `category_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='道具分类表';

-- 道具表
CREATE TABLE IF NOT EXISTS `props` (
  `prop_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '道具名称',
  `category_id` bigint(20) NOT NULL COMMENT '分类ID',
  `icon_path` varchar(255) NOT NULL COMMENT '图标路径',
  `coin_cost` int(11) DEFAULT 0 COMMENT '消耗金币数',
  `effect_description` text DEFAULT NULL COMMENT '使用效果描述',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`prop_id`),
  KEY `idx_category_id` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='道具表';

-- 用户道具持有表
CREATE TABLE IF NOT EXISTS `user_props` (
  `user_id` bigint(20) NOT NULL,
  `prop_id` bigint(20) NOT NULL,
  `quantity` int(11) DEFAULT 0 COMMENT '持有数量',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `prop_id`),
  KEY `idx_prop_id` (`prop_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户道具持有表';

-- 社区标签表
CREATE TABLE IF NOT EXISTS `tags` (
  `tag_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '标签名称',
  `icon` varchar(255) DEFAULT NULL COMMENT '标签图标',
  `post_count` int(11) DEFAULT 0 COMMENT '帖子数量',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认展示',
  `sort_order` int(11) DEFAULT 0 COMMENT '排序序号',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='社区标签表';

-- 帖子表
CREATE TABLE IF NOT EXISTS `posts` (
  `post_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '发布者ID',
  `title` varchar(200) NOT NULL COMMENT '标题',
  `content` text DEFAULT NULL COMMENT '内容',
  `images` json DEFAULT NULL COMMENT '图片列表',
  `videos` json DEFAULT NULL COMMENT '视频列表',
  `tag_ids` json DEFAULT NULL COMMENT '标签ID列表',
  `like_count` int(11) DEFAULT 0 COMMENT '点赞数',
  `comment_count` int(11) DEFAULT 0 COMMENT '评论数',
  `view_count` int(11) DEFAULT 0 COMMENT '浏览数',
  `status` enum('normal','draft','deleted') DEFAULT 'normal' COMMENT '状态',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帖子表';

-- 评论表
CREATE TABLE IF NOT EXISTS `comments` (
  `comment_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `post_id` bigint(20) NOT NULL COMMENT '帖子ID',
  `user_id` bigint(20) NOT NULL COMMENT '评论者ID',
  `content` text NOT NULL COMMENT '评论内容',
  `like_count` int(11) DEFAULT 0 COMMENT '点赞数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`comment_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

-- 用户关系表（关注）
CREATE TABLE IF NOT EXISTS `user_relations` (
  `user_id` bigint(20) NOT NULL COMMENT '关注者ID',
  `target_user_id` bigint(20) NOT NULL COMMENT '被关注者ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `target_user_id`),
  KEY `idx_target_user_id` (`target_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关系表';

-- 用户互动表（点赞）
CREATE TABLE IF NOT EXISTS `user_interactions` (
  `user_id` bigint(20) NOT NULL,
  `target_type` enum('post','comment') NOT NULL COMMENT '目标类型',
  `target_id` bigint(20) NOT NULL COMMENT '目标ID',
  `interaction_type` enum('like','unlike') NOT NULL COMMENT '互动类型',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`, `target_type`, `target_id`),
  KEY `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户互动表';

-- 聊天消息表
CREATE TABLE IF NOT EXISTS `chat_messages` (
  `message_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `sender` enum('user','avatar') NOT NULL COMMENT '发送方',
  `content` text NOT NULL COMMENT '消息内容',
  `message_type` enum('text','prop_use','note') DEFAULT 'text' COMMENT '消息类型',
  `related_id` bigint(20) DEFAULT NULL COMMENT '关联ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`message_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息表';

-- 小纸条表
CREATE TABLE IF NOT EXISTS `messages` (
  `message_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL COMMENT '小纸条内容',
  `message_type` enum('free','paid') DEFAULT 'free' COMMENT '类型',
  `unlock_coins` int(11) DEFAULT 0 COMMENT '解锁所需金币',
  `pet_id` bigint(20) DEFAULT NULL COMMENT '关联宠物ID',
  `is_ai_generated` tinyint(1) DEFAULT 1 COMMENT '是否AI生成',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`message_id`),
  KEY `idx_pet_id` (`pet_id`),
  KEY `idx_message_type` (`message_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小纸条表';

-- 用户解锁记录表
CREATE TABLE IF NOT EXISTS `user_unlock_records` (
  `record_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `message_id` bigint(20) NOT NULL COMMENT '小纸条ID',
  `coins_spent` int(11) DEFAULT 0 COMMENT '消耗金币',
  `unlock_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`record_id`),
  UNIQUE KEY `uk_user_message` (`user_id`, `message_id`),
  KEY `idx_message_id` (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户解锁记录表';

-- 插入默认数据
INSERT INTO `avatars` (`name`, `resource_path`, `sort_order`, `is_default`, `is_preset`) VALUES
('默认猫咪', '/avatars/default_cat.png', 1, 1, 1),
('默认狗狗', '/avatars/default_dog.png', 2, 0, 1),
('可爱猫咪', '/avatars/cute_cat.png', 3, 0, 1),
('忠诚狗狗', '/avatars/loyal_dog.png', 4, 0, 1);

INSERT INTO `prop_categories` (`name`, `sort_order`) VALUES
('食物', 1),
('玩具', 2),
('护理', 3),
('装饰', 4);

INSERT INTO `props` (`name`, `category_id`, `icon_path`, `coin_cost`, `effect_description`, `sort_order`) VALUES
('猫粮', 1, '/props/cat_food.png', 5, '让猫咪吃饱饱', 1),
('狗粮', 1, '/props/dog_food.png', 5, '让狗狗吃饱饱', 2),
('逗猫棒', 2, '/props/cat_wand.png', 10, '和猫咪玩耍', 3),
('飞盘', 2, '/props/frisbee.png', 10, '和狗狗玩耍', 4),
('梳子', 3, '/props/comb.png', 8, '给宠物梳毛', 5),
('项圈', 4, '/props/collar.png', 15, '给宠物戴上项圈', 6);

INSERT INTO `tags` (`name`, `icon`, `is_default`, `sort_order`) VALUES
('日常', '/tags/daily.png', 1, 1),
('知识', '/tags/knowledge.png', 1, 2),
('信息', '/tags/info.png', 1, 3),
('种草', '/tags/recommend.png', 1, 4); 