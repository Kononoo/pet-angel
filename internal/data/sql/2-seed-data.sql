-- Pet Angel 示例种子数据（可重复执行）
-- 说明：
-- 1) 严格使用字符串/DATETIME，不使用 JSON/ENUM；资源全部是 URL 字符串
-- 2) 通过唯一键/INSERT ... SELECT 保证幂等，避免主键冲突
-- 3) 如需重置，可先执行 1-init-tables.sql 重新建表

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =========================
-- 宠物模型（猫/狗，多款，含默认）
-- =========================
INSERT INTO `pet_models`(`name`,`path`,`type`,`is_default`,`sort_order`)
SELECT '天使猫-默认','/models/cat/angel_default.glb',0,1,1 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使猫-默认')
UNION ALL SELECT '天使猫-粉','/models/cat/angel_pink.glb',0,0,2 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使猫-粉')
UNION ALL SELECT '天使猫-蓝','/models/cat/angel_blue.glb',0,0,3 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使猫-蓝')
UNION ALL SELECT '天使狗-默认','/models/dog/angel_default.glb',1,1,1 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使狗-默认')
UNION ALL SELECT '天使狗-金毛','/models/dog/golden.glb',1,0,2 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使狗-金毛')
UNION ALL SELECT '天使狗-柯基','/models/dog/corgi.glb',1,0,3 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM pet_models WHERE name='天使狗-柯基');

-- =========================
-- 道具（金币消耗数值对齐产品方案）
-- =========================
INSERT INTO `items`(`name`,`description`,`icon_path`,`coin_cost`)
SELECT '普通猫粮（单次）','满足饥饿度','/items/food_basic.png',20 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='普通猫粮（单次）')
UNION ALL SELECT '营养套餐（3次）','心情提升','/items/food_pack.png',50 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='营养套餐（3次）')
UNION ALL SELECT '豪华鲜食（周卡）','每日自动喂养+特效','/items/food_week.png',120 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='豪华鲜食（周卡）')
UNION ALL SELECT '限定节日食物','特殊外形+全屏特效','/items/food_festival.png',80 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='限定节日食物')
UNION ALL SELECT '基础猫砂（单次）','清洁度恢复至70%','/items/litter_basic.png',30 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='基础猫砂（单次）')
UNION ALL SELECT '除臭猫砂（3次）','清洁度100%+24小时除臭','/items/litter_odor.png',80 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='除臭猫砂（3次）')
UNION ALL SELECT '智能厕所（永久）','自动清洁+每日维持','/items/toilet_smart.png',300 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='智能厕所（永久）')
UNION ALL SELECT '逗猫棒（单次）','心情提升','/items/toy_wand.png',15 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='逗猫棒（单次）')
UNION ALL SELECT '电动老鼠（3次）','心情提升','/items/toy_mouse.png',40 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='电动老鼠（3次）')
UNION ALL SELECT '猫爬架（永久）','心情提升','/items/cat_tree.png',200 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='猫爬架（永久）')
UNION ALL SELECT '智能摄像头（月卡）','远程查看与录像','/items/camera_month.png',180 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='智能摄像头（月卡）')
UNION ALL SELECT '专属头像框','个人主页外观','/items/avatar_frame.png',60 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM items WHERE name='专属头像框');

-- =========================
-- 社区分类（丧宠关怀 + 通用养宠）
-- 注意：1-init 中已含 日常/知识/信息/种草，这里补充其余标签
-- =========================
INSERT INTO `categories`(`name`,`sort_order`)
SELECT '喵星登陆',11 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='喵星登陆')
UNION ALL SELECT '情感互助',12 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='情感互助')
UNION ALL SELECT '回忆分享',13 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='回忆分享')
UNION ALL SELECT '治愈日常',14 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='治愈日常')
UNION ALL SELECT '走出阴霾',15 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='走出阴霾')
UNION ALL SELECT '猫猫',21 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='猫猫')
UNION ALL SELECT '狗狗',22 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='狗狗')
UNION ALL SELECT '新手',23 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='新手')
UNION ALL SELECT '健康',24 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='健康')
UNION ALL SELECT '训练',25 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='训练')
UNION ALL SELECT '用品',26 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='用品')
UNION ALL SELECT '食物',27 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='食物')
UNION ALL SELECT '情感',28 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='情感')
UNION ALL SELECT '隐私日记',29 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='隐私日记')
UNION ALL SELECT '萌宠摄影',30 FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM categories WHERE name='萌宠摄影');

-- =========================
-- 用户（示例数据，密码明文 123456，便于联调；生产必须改为 bcrypt）
-- model_url 以当前所选 pet_models.path 为准，可后续由服务端刷新
-- =========================
INSERT INTO `users`(`nickname`,`username`,`password`,`avatar`,`model_id`,`model_url`,`pet_name`,`pet_avatar`,`pet_sex`,`kind`,`weight`,`hobby`,`description`,`coins`,`created_at`,`updated_at`)
SELECT 'Lemon','lemon','123456','/avatars/u1.png',1,'/models/cat/angel_default.glb','柚子','/pets/p1.png',2,'cat',4,'晒太阳','爱笑的猫咪',300,'2025-03-01 10:00:00','2025-03-01 10:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='lemon')
UNION ALL SELECT 'Ache','ache','123456','/avatars/u2.png',4,'/models/dog/angel_default.glb','阿奇','/pets/p2.png',1,'dog',9,'跑步','元气小狗',260,'2025-03-02 11:00:00','2025-03-02 11:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='ache')
UNION ALL SELECT 'momo','momo','123456','/avatars/u3.png',2,'/models/cat/angel_pink.glb','桃桃','/pets/p3.png',2,'cat',5,'午睡','慵懒可爱',180,'2025-03-03 12:00:00','2025-03-03 12:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='momo')
UNION ALL SELECT 'Dev','dev','123456','/avatars/u4.png',5,'/models/dog/golden.glb','可可','/pets/p4.png',1,'dog',18,'捡球','忠诚陪伴',120,'2025-03-04 13:00:00','2025-03-04 13:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='dev')
UNION ALL SELECT 'Nana','nana','123456','/avatars/u5.png',3,'/models/cat/angel_blue.glb','奶糖','/pets/p5.png',2,'cat',3,'梳毛','软萌仔',90,'2025-03-05 14:00:00','2025-03-05 14:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='nana')
UNION ALL SELECT 'Yuki','yuki','123456','/avatars/u6.png',6,'/models/dog/corgi.glb','雪球','/pets/p6.png',1,'dog',10,'摄影','短腿行侠',60,'2025-03-06 15:00:00','2025-03-06 15:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='yuki')
UNION ALL SELECT 'Kiki','kiki','123456','/avatars/u7.png',1,'/models/cat/angel_default.glb','奇奇','/pets/p7.png',2,'cat',6,'晒太阳','爱吃罐头',40,'2025-03-07 16:20:00','2025-03-07 16:20:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='kiki')
UNION ALL SELECT '七七','seven','123456','/avatars/u8.png',1,'/models/cat/angel_default.glb','七七','/pets/p8.png',2,'cat',7,'爬窗台','温柔粘人',500,'2025-03-08 09:30:00','2025-03-08 09:30:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='seven')
UNION ALL SELECT 'Lemon-粉丝','lemon_fan','123456','/avatars/u9.png',1,'/models/cat/angel_default.glb','小柚','/pets/p9.png',2,'cat',4,'收藏','忠实粉丝',80,'2025-03-09 09:10:00','2025-03-09 09:10:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='lemon_fan')
UNION ALL SELECT '摄影师','photographer','123456','/avatars/u10.png',4,'/models/dog/angel_default.glb','阿布','/pets/p10.png',1,'dog',14,'摄影','镜头里的宠物',220,'2025-03-09 10:20:00','2025-03-09 10:20:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='photographer');

-- =========================
-- 关注关系（去重幂等）
-- =========================
INSERT INTO `user_follows` (`follower_id`,`followee_id`,`created_at`)
SELECT u1.id, u2.id, '2025-03-10 10:00:00'
FROM users u1, users u2
WHERE u1.username IN ('lemon','ache','momo','dev')
  AND u2.username IN ('nana','yuki','kiki','seven')
  AND u1.id <> u2.id
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);

-- 额外关注：粉丝关注 Lemon、Kiki
INSERT INTO `user_follows`(`follower_id`,`followee_id`,`created_at`)
SELECT f.id, t.id, '2025-03-10 12:00:00'
FROM users f, users t
WHERE f.username IN ('lemon_fan') AND t.username IN ('lemon','kiki') AND f.id<>t.id
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);

-- =========================
-- 帖子（图文/视频混合）
-- 使用唯一标题便于后续通过标题定位点赞/评论
-- =========================
INSERT INTO `posts`(`user_id`,`category_id`,`title`,`content`,`type`,`image_urls`,`video_url`,`cover_url`,`locate`,`tags`,`liked_count`,`comment_count`,`is_private`,`created_at`,`updated_at`)
SELECT u.id, c.id, '初见·小天使猫','第一次相遇，愿治愈每一天',0,
       '/images/c1_1.jpg,/images/c1_2.jpg',NULL,'/covers/c1.jpg','杭州','猫,治愈',37,5,0,'2025-03-11 09:00:00','2025-03-11 09:00:00'
FROM users u JOIN categories c ON u.username='lemon' AND c.name='治愈日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='初见·小天使猫')
UNION ALL
SELECT u.id, c.id, '训练的第五种技巧？','大家有更好的方法嘛',0,
       '/images/c2_1.jpg,/images/c2_2.jpg',NULL,'/covers/c2.jpg','上海','狗,训练',58,9,0,'2025-03-12 10:00:00','2025-03-12 10:00:00'
FROM users u JOIN categories c ON u.username='ache' AND c.name='训练'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='训练的第五种技巧？')
UNION ALL
SELECT u.id, c.id, '私密回忆·那一天','我会慢慢走出阴霾',0,
       '/images/c3_1.jpg',NULL,'/covers/c3.jpg','南京','回忆,情感',12,3,1,'2025-03-13 20:00:00','2025-03-13 20:00:00'
FROM users u JOIN categories c ON u.username='momo' AND c.name='走出阴霾'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='私密回忆·那一天')
UNION ALL
SELECT u.id, c.id, '春日公园VLOG','和可可在草地奔跑',1,
       NULL,'/videos/v1.mp4','/covers/v1.jpg','北京','狗,日常',1713,20,0,'2025-03-14 08:30:00','2025-03-14 08:30:00'
FROM users u JOIN categories c ON u.username='dev' AND c.name='日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='春日公园VLOG')
UNION ALL
SELECT u.id, c.id, '晒晒TA的睡姿','今天也睡成可爱的面包',0,
       '/images/c4_1.jpg,/images/c4_2.jpg',NULL,'/covers/c4.jpg','深圳','猫,日常',89,6,0,'2025-03-15 21:00:00','2025-03-15 21:00:00'
FROM users u JOIN categories c ON u.username='kiki' AND c.name='日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='晒晒TA的睡姿')
UNION ALL
SELECT u.id, c.id, '新手养猫清单','给新手的 10 条建议',0,
       '/images/c5_1.jpg',NULL,'/covers/c5.jpg','广州','新手,清单',132,14,0,'2025-03-16 09:30:00','2025-03-16 09:30:00'
FROM users u JOIN categories c ON u.username='nana' AND c.name='新手'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='新手养猫清单')
UNION ALL
SELECT u.id, c.id, '年度体检记录','记录健康每一步',0,
       '/images/c6_1.jpg',NULL,'/covers/c6.jpg','杭州','健康,体检',45,8,0,'2025-03-17 10:00:00','2025-03-17 10:00:00'
FROM users u JOIN categories c ON u.username='lemon' AND c.name='健康'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='年度体检记录');

-- =========================
-- 评论（通过标题+用户定位 post_id）
-- =========================
INSERT INTO `comments`(`post_id`,`user_id`,`content`,`liked_count`,`created_at`)
SELECT p.id, u.id, '好可爱，愿你每天被温柔以待～', 6, '2025-03-11 10:00:00'
FROM posts p JOIN users u ON p.title='初见·小天使猫' AND u.username='nana'
WHERE NOT EXISTS(
  SELECT 1 FROM comments c WHERE c.post_id=p.id AND c.user_id=u.id AND c.content='好可爱，愿你每天被温柔以待～')
UNION ALL
SELECT p.id, u.id, '训练分享很实用，我也试试', 3, '2025-03-12 12:00:00'
FROM posts p JOIN users u ON p.title='训练的第五种技巧？' AND u.username='yuki'
WHERE NOT EXISTS(
  SELECT 1 FROM comments c WHERE c.post_id=p.id AND c.user_id=u.id AND c.content='训练分享很实用，我也试试')
UNION ALL
SELECT p.id, u.id, '抱抱你，会越来越好的', 9, '2025-03-13 21:00:00'
FROM posts p JOIN users u ON p.title='私密回忆·那一天' AND u.username='lemon'
WHERE NOT EXISTS(
  SELECT 1 FROM comments c WHERE c.post_id=p.id AND c.user_id=u.id AND c.content='抱抱你，会越来越好的')
UNION ALL
SELECT p.id, u.id, 'vlog镜头好稳！', 2, '2025-03-14 09:00:00'
FROM posts p JOIN users u ON p.title='春日公园VLOG' AND u.username='kiki'
WHERE NOT EXISTS(
  SELECT 1 FROM comments c WHERE c.post_id=p.id AND c.user_id=u.id AND c.content='vlog镜头好稳！')
UNION ALL
SELECT p.id, u.id, '收藏了，太实用了', 1, '2025-03-16 10:10:00'
FROM posts p JOIN users u ON p.title='新手养猫清单' AND u.username='lemon_fan'
WHERE NOT EXISTS(
  SELECT 1 FROM comments c WHERE c.post_id=p.id AND c.user_id=u.id AND c.content='收藏了，太实用了');

-- =========================
-- 点赞（帖子/评论）
-- =========================
-- 帖子点赞
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 0, p.id, '2025-03-15 10:00:00' FROM users u JOIN posts p ON u.username='seven' AND p.title='春日公园VLOG'
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 0, p.id, '2025-03-15 11:00:00' FROM users u JOIN posts p ON u.username='nana' AND p.title='初见·小天使猫'
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);
-- 评论点赞
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 1, c.id, '2025-03-15 12:00:00' FROM users u JOIN comments c ON u.username='momo' AND c.content='好可爱，愿你每天被温柔以待～'
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);
-- 额外点赞
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 0, p.id, '2025-03-16 10:20:00' FROM users u JOIN posts p ON u.username='photographer' AND p.title='晒晒TA的睡姿'
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);

-- =========================
-- 消息/小纸条（message_type：0聊天 1小纸条）
-- =========================
-- 聊天示例：用户→AI 两条
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`) 
SELECT u.id, 0, 0, 0, 0, '今天有点想它了……', '2025-03-16 20:00:00' FROM users u WHERE u.username='lemon'
UNION ALL
SELECT u.id, 1, 0, 0, 0, '我在呢，会一直陪着你～', '2025-03-16 20:00:05' FROM users u WHERE u.username='lemon';

-- 小纸条示例：每日晚间隐私留言（锁定，20金币）
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
SELECT u.id, 1, 1, 1, 20, '今晚也别熬夜，我会在梦里见到你', '2025-03-16 21:30:00' FROM users u WHERE u.username='lemon'
UNION ALL
SELECT u.id, 1, 1, 1, 20, '别担心，我一直在你身边', '2025-03-16 21:31:00' FROM users u WHERE u.username='momo';

-- 更多消息/小纸条：为多位用户生成样例
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
SELECT u.id, 1, 1, 1, 20, '今晚我会在窗边等你回来', '2025-03-17 21:10:00' FROM users u WHERE u.username='kiki'
UNION ALL
SELECT u.id, 1, 1, 1, 20, '别忘了早点休息，明天一起晒太阳', '2025-03-17 21:15:00' FROM users u WHERE u.username='nana'
UNION ALL
SELECT u.id, 0, 0, 0, 0, '今天拍到它打哈欠啦', '2025-03-17 18:20:00' FROM users u WHERE u.username='photographer'
UNION ALL
SELECT u.id, 1, 0, 0, 0, '你拍得真好看～', '2025-03-17 18:20:05' FROM users u WHERE u.username='photographer';

SET FOREIGN_KEY_CHECKS = 1;

