-- Pet Angel 海量数据生成脚本
-- 生成真实的生产环境数据：用户、帖子、小纸条、聊天记录等

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =========================
-- 1. 更新宠物模型路径（使用实际的GLB文件）
-- =========================
UPDATE pet_models SET path = '/models/Cat_1.glb' WHERE name LIKE '%猫%' AND sort_order = 1;
UPDATE pet_models SET path = '/models/Cat_2.glb' WHERE name LIKE '%猫%' AND sort_order = 2;
UPDATE pet_models SET path = '/models/Cat_3.glb' WHERE name LIKE '%猫%' AND sort_order = 3;
UPDATE pet_models SET path = '/models/Cat_4.glb' WHERE name LIKE '%猫%' AND sort_order = 4;
UPDATE pet_models SET path = '/models/Cat_5.glb' WHERE name LIKE '%猫%' AND sort_order = 5;
UPDATE pet_models SET path = '/models/Dog_1.glb' WHERE name LIKE '%狗%' AND sort_order = 1;
UPDATE pet_models SET path = '/models/Dog_2.glb' WHERE name LIKE '%狗%' AND sort_order = 2;
UPDATE pet_models SET path = '/models/Dog_3.glb' WHERE name LIKE '%狗%' AND sort_order = 3;
UPDATE pet_models SET path = '/models/Dog_4.glb' WHERE name LIKE '%狗%' AND sort_order = 4;
UPDATE pet_models SET path = '/models/Dog_5.glb' WHERE name LIKE '%狗%' AND sort_order = 5;

-- =========================
-- 2. 生成更多用户（15个用户）
-- =========================
INSERT INTO `users`(`nickname`,`username`,`password`,`avatar`,`model_id`,`model_url`,`pet_name`,`pet_avatar`,`pet_sex`,`kind`,`weight`,`hobby`,`description`,`coins`,`created_at`,`updated_at`)
SELECT '小天使','angel001','123456','/static/avatar/angel001.jpg',1,'/models/Cat_1.glb','小天使','/static/pet/angel001.jpg',2,'cat',4,'晒太阳,睡觉','治愈系小天使',450,'2025-01-15 10:00:00','2025-01-15 10:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='angel001')
UNION ALL SELECT '治愈师','healer002','123456','/static/avatar/healer002.jpg',6,'/models/Dog_1.glb','治愈','/static/pet/healer002.jpg',1,'dog',12,'陪伴,治愈','专业治愈师',380,'2025-01-20 11:00:00','2025-01-20 11:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='healer002')
UNION ALL SELECT '回忆录','memory003','123456','/static/avatar/memory003.jpg',2,'/models/Cat_2.glb','回忆','/static/pet/memory003.jpg',2,'cat',5,'回忆,记录','记录美好回忆',520,'2025-02-01 12:00:00','2025-02-01 12:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='memory003')
UNION ALL SELECT '守护者','guardian004','123456','/static/avatar/guardian004.jpg',7,'/models/Dog_2.glb','守护','/static/pet/guardian004.jpg',1,'dog',15,'守护,陪伴','永远的守护者',290,'2025-02-10 13:00:00','2025-02-10 13:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='guardian004')
UNION ALL SELECT '温暖家','warm005','123456','/static/avatar/warm005.jpg',3,'/models/Cat_3.glb','温暖','/static/pet/warm005.jpg',2,'cat',6,'温暖,陪伴','温暖的家',180,'2025-02-15 14:00:00','2025-02-15 14:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='warm005')
UNION ALL SELECT '快乐狗','happy006','123456','/static/avatar/happy006.jpg',8,'/models/Dog_3.glb','快乐','/static/pet/happy006.jpg',1,'dog',8,'快乐,奔跑','快乐的小短腿',320,'2025-02-20 15:00:00','2025-02-20 15:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='happy006')
UNION ALL SELECT '温柔猫','gentle007','123456','/static/avatar/gentle007.jpg',4,'/models/Cat_4.glb','温柔','/static/pet/gentle007.jpg',2,'cat',7,'温柔,粘人','温柔的布偶猫',410,'2025-03-01 16:00:00','2025-03-01 16:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='gentle007')
UNION ALL SELECT '忠诚狗','loyal008','123456','/static/avatar/loyal008.jpg',9,'/models/Dog_4.glb','忠诚','/static/pet/loyal008.jpg',1,'dog',14,'忠诚,陪伴','忠诚的伙伴',260,'2025-03-05 17:00:00','2025-03-05 17:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='loyal008')
UNION ALL SELECT '治愈猫','healcat009','123456','/static/avatar/healcat009.jpg',5,'/models/Cat_5.glb','治愈','/static/pet/healcat009.jpg',2,'cat',4,'治愈,陪伴','治愈系暹罗',350,'2025-03-10 18:00:00','2025-03-10 18:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='healcat009')
UNION ALL SELECT '陪伴狗','companion010','123456','/static/avatar/companion010.jpg',10,'/models/Dog_5.glb','陪伴','/static/pet/companion010.jpg',1,'dog',13,'陪伴,聪明','聪明的边牧',480,'2025-03-15 19:00:00','2025-03-15 19:00:00' FROM DUAL WHERE NOT EXISTS(SELECT 1 FROM users WHERE username='companion010');

-- =========================
-- 3. 生成海量帖子（150+ 帖子）
-- =========================
INSERT INTO `posts`(`user_id`,`category_id`,`title`,`content`,`type`,`image_urls`,`video_url`,`cover_url`,`locate`,`tags`,`liked_count`,`comment_count`,`is_private`,`created_at`,`updated_at`)
SELECT u.id, c.id, '初见·小天使猫','第一次相遇，愿治愈每一天。它用那双温柔的眼睛看着我，仿佛在说：别难过，我会一直陪着你。',0,
       '/static/image/2025/01/15/cat_angel_001.jpg,/static/image/2025/01/15/cat_angel_002.jpg',NULL,'/static/image/2025/01/15/cat_angel_cover.jpg','杭州','猫,治愈,天使',156,23,0,'2025-01-15 10:30:00','2025-01-15 10:30:00'
FROM users u JOIN categories c ON u.username='angel001' AND c.name='治愈日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='初见·小天使猫')

UNION ALL SELECT u.id, c.id, '治愈系金毛的日常','每天最治愈的时刻，就是看着它摇着尾巴向我跑来。失去爱宠的痛苦，在它的陪伴下慢慢愈合。',0,
       '/static/image/2025/01/20/golden_heal_001.jpg,/static/image/2025/01/20/golden_heal_002.jpg,/static/image/2025/01/20/golden_heal_003.jpg',NULL,'/static/image/2025/01/20/golden_heal_cover.jpg','北京','狗,治愈,陪伴',203,31,0,'2025-01-20 11:30:00','2025-01-20 11:30:00'
FROM users u JOIN categories c ON u.username='healer002' AND c.name='情感互助'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='治愈系金毛的日常')

UNION ALL SELECT u.id, c.id, '回忆录：与美短的365天','记录与美短相处的每一天，每一个温馨的瞬间都值得珍藏。它教会了我如何更好地爱与被爱。',1,
       NULL,'/static/video/2025/02/01/memory_cat_001.mp4','/static/video/2025/02/01/memory_cat_cover.jpg','上海','猫,回忆,记录',89,15,0,'2025-02-01 12:30:00','2025-02-01 12:30:00'
FROM users u JOIN categories c ON u.username='memory003' AND c.name='回忆分享'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='回忆录：与美短的365天')

UNION ALL SELECT u.id, c.id, '守护者的誓言','德牧的忠诚让我明白，真正的守护不是占有，而是无条件的陪伴与保护。',0,
       '/static/image/2025/02/10/guardian_dog_001.jpg,/static/image/2025/02/10/guardian_dog_002.jpg',NULL,'/static/image/2025/02/10/guardian_dog_cover.jpg','深圳','狗,守护,忠诚',178,28,0,'2025-02-10 13:30:00','2025-02-10 13:30:00'
FROM users u JOIN categories c ON u.username='guardian004' AND c.name='走出阴霾'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='守护者的誓言')

UNION ALL SELECT u.id, c.id, '温暖的家：橘猫的治愈时光','橘猫用它的温暖治愈了我的孤独，每天最期待的就是回家看到它。',0,
       '/static/image/2025/02/15/orange_cat_001.jpg,/static/image/2025/02/15/orange_cat_002.jpg',NULL,'/static/image/2025/02/15/orange_cat_cover.jpg','广州','猫,温暖,治愈',134,19,0,'2025-02-15 14:30:00','2025-02-15 14:30:00'
FROM users u JOIN categories c ON u.username='warm005' AND c.name='治愈日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='温暖的家：橘猫的治愈时光')

UNION ALL SELECT u.id, c.id, '快乐柯基的奔跑时光','柯基的小短腿跑起来特别可爱，每次看到它都会忘记烦恼。',1,
       NULL,'/static/video/2025/02/20/happy_corgi_001.mp4','/static/video/2025/02/20/happy_corgi_cover.jpg','成都','狗,快乐,柯基',267,42,0,'2025-02-20 15:30:00','2025-02-20 15:30:00'
FROM users u JOIN categories c ON u.username='happy006' AND c.name='日常'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='快乐柯基的奔跑时光')

UNION ALL SELECT u.id, c.id, '温柔布偶的粘人时光','布偶猫的温柔让我重新相信爱情，它总是用最温柔的方式陪伴我。',0,
       '/static/image/2025/03/01/ragdoll_gentle_001.jpg,/static/image/2025/03/01/ragdoll_gentle_002.jpg,/static/image/2025/03/01/ragdoll_gentle_003.jpg',NULL,'/static/image/2025/03/01/ragdoll_gentle_cover.jpg','武汉','猫,温柔,布偶',198,35,0,'2025-03-01 16:30:00','2025-03-01 16:30:00'
FROM users u JOIN categories c ON u.username='gentle007' AND c.name='情感'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='温柔布偶的粘人时光')

UNION ALL SELECT u.id, c.id, '忠诚拉布拉多的陪伴','拉布拉多的忠诚让我明白，真正的朋友永远不会离开。',0,
       '/static/image/2025/03/05/labrador_loyal_001.jpg,/static/image/2025/03/05/labrador_loyal_002.jpg',NULL,'/static/image/2025/03/05/labrador_loyal_cover.jpg','南京','狗,忠诚,陪伴',145,26,0,'2025-03-05 17:30:00','2025-03-05 17:30:00'
FROM users u JOIN categories c ON u.username='loyal008' AND c.name='情感互助'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='忠诚拉布拉多的陪伴')

UNION ALL SELECT u.id, c.id, '治愈系暹罗的智慧','暹罗猫的智慧让我惊讶，它总能理解我的心情。',0,
       '/static/image/2025/03/10/siamese_heal_001.jpg,/static/image/2025/03/10/siamese_heal_002.jpg',NULL,'/static/image/2025/03/10/siamese_heal_cover.jpg','西安','猫,智慧,暹罗',167,29,0,'2025-03-10 18:30:00','2025-03-10 18:30:00'
FROM users u JOIN categories c ON u.username='healcat009' AND c.name='知识'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='治愈系暹罗的智慧')

UNION ALL SELECT u.id, c.id, '聪明边牧的训练日常','边牧的聪明让我惊叹，每天的训练都充满惊喜。',1,
       NULL,'/static/video/2025/03/15/border_collie_train_001.mp4','/static/video/2025/03/15/border_collie_train_cover.jpg','重庆','狗,聪明,训练',312,48,0,'2025-03-15 19:30:00','2025-03-15 19:30:00'
FROM users u JOIN categories c ON u.username='companion010' AND c.name='训练'
WHERE NOT EXISTS(SELECT 1 FROM posts WHERE title='聪明边牧的训练日常');

-- 继续生成更多帖子（批量插入）
INSERT INTO `posts`(`user_id`,`category_id`,`title`,`content`,`type`,`image_urls`,`video_url`,`cover_url`,`locate`,`tags`,`liked_count`,`comment_count`,`is_private`,`created_at`,`updated_at`)
SELECT u.id, c.id, CONCAT('治愈日记第', n.n, '天'), CONCAT('第', n.n, '天的治愈时光，每一天都在慢慢变好。'), 0,
       CONCAT('/static/image/2025/01/', LPAD(n.n, 2, '0'), '/heal_diary_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/01/', LPAD(n.n, 2, '0'), '/heal_diary_', LPAD(n.n, 3, '0'), '_cover.jpg'), '杭州','治愈,日记', FLOOR(50 + RAND() * 150), FLOOR(5 + RAND() * 20), 0, DATE_ADD('2025-01-15 10:30:00', INTERVAL n.n DAY), DATE_ADD('2025-01-15 10:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='angel001' AND c.name='治愈日常' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('治愈日记第', n.n, '天'))

UNION ALL SELECT u.id, c.id, CONCAT('陪伴时光第', n.n, '期'), CONCAT('第', n.n, '期的陪伴时光，感谢有你。'), 0,
       CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/companion_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/companion_', LPAD(n.n, 3, '0'), '_cover.jpg'), '北京','陪伴,时光', FLOOR(60 + RAND() * 180), FLOOR(8 + RAND() * 25), 0, DATE_ADD('2025-02-01 11:30:00', INTERVAL n.n DAY), DATE_ADD('2025-02-01 11:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='healer002' AND c.name='情感互助' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('陪伴时光第', n.n, '期'));

-- =========================
-- 4. 生成海量小纸条（每个用户15+ 小纸条）
-- =========================
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
-- angel001 的小纸条
SELECT u.id, 1, 1, 0, 0, '早安，小天使！今天也要开开心心的哦～', '2025-01-15 08:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%早安，小天使%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '今天想告诉你一个小秘密：你比想象中更坚强。', '2025-01-16 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%小秘密%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！记得给自己一个温暖的拥抱。', '2025-01-17 12:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%午安%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你知道吗？每一次微笑都是对过去的勇敢告别。', '2025-01-18 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%微笑%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，愿你梦里有温暖的阳光。', '2025-01-19 22:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%晚安%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '今天的你比昨天更勇敢，这就是成长。', '2025-01-20 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%勇敢%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！新的一天，新的希望。', '2025-01-21 08:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%新的一天%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '记住，你值得被爱，也值得拥有幸福。', '2025-01-22 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%值得被爱%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！记得喝水，照顾好自己。', '2025-01-23 12:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%喝水%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一个治愈的瞬间，都是你内心的力量。', '2025-01-24 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈的瞬间%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，明天会更好。', '2025-01-25 22:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%明天会更好%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的温柔，是这个世界最美的风景。', '2025-01-26 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最美的风景%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！今天也要保持微笑。', '2025-01-27 08:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%保持微笑%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '时间会治愈一切，而你正在创造奇迹。', '2025-01-28 20:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%创造奇迹%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！记得给自己一些时间。', '2025-01-29 12:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%给自己一些时间%');

-- healer002 的小纸条
SELECT u.id, 1, 1, 0, 0, '早安，治愈师！今天也要温暖他人。', '2025-01-20 08:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈师%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的存在，本身就是一种治愈。', '2025-01-21 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%本身就是一种治愈%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的温暖正在传递。', '2025-01-22 12:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%温暖正在传递%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次陪伴，都是对生命的尊重。', '2025-01-23 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对生命的尊重%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，治愈的力量在你心中。', '2025-01-24 22:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈的力量%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的善良，是这个世界的光。', '2025-01-25 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%这个世界的光%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！新的一天，新的治愈。', '2025-01-26 08:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%新的治愈%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '陪伴是最长情的告白。', '2025-01-27 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最长情的告白%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的存在让世界更美好。', '2025-01-28 12:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%让世界更美好%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次治愈，都是对爱的诠释。', '2025-01-29 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对爱的诠释%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，明天继续传递温暖。', '2025-01-30 22:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%继续传递温暖%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的温柔，是治愈的良药。', '2025-01-31 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈的良药%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！今天也要温暖如初。', '2025-02-01 08:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%温暖如初%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次陪伴，都是对生命的礼赞。', '2025-02-02 20:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对生命的礼赞%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的光芒正在闪耀。', '2025-02-03 12:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%光芒正在闪耀%');

-- =========================
-- 5. 生成聊天对话记录
-- =========================
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
-- angel001 的聊天记录
SELECT u.id, 0, 0, 0, 0, '你好，小天使', '2025-01-15 10:00:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='你好，小天使')
UNION ALL SELECT u.id, 1, 0, 0, 0, '你好！我是你的专属小天使，很高兴认识你！', '2025-01-15 10:01:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%专属小天使%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '我今天心情不太好', '2025-01-15 10:02:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='我今天心情不太好')
UNION ALL SELECT u.id, 1, 0, 0, 0, '没关系，每个人都会有心情不好的时候。让我陪着你，给你一个温暖的拥抱吧！', '2025-01-15 10:03:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%温暖的拥抱%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '谢谢你，我感觉好多了', '2025-01-15 10:04:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='谢谢你，我感觉好多了')
UNION ALL SELECT u.id, 1, 0, 0, 0, '太好了！记住，你永远不是一个人，我会一直在这里陪着你。', '2025-01-15 10:05:00' FROM users u WHERE u.username='angel001' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%永远不是一个人%');

-- healer002 的聊天记录
SELECT u.id, 0, 0, 0, 0, '你好，治愈师', '2025-01-20 11:00:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='你好，治愈师')
UNION ALL SELECT u.id, 1, 0, 0, 0, '你好！我是你的专属治愈师，让我们一起面对生活的挑战吧！', '2025-01-20 11:01:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%专属治愈师%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '我最近失去了我的宠物', '2025-01-20 11:02:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='我最近失去了我的宠物')
UNION ALL SELECT u.id, 1, 0, 0, 0, '我理解你的痛苦。失去一个深爱的伙伴是人生中最艰难的经历之一。但请记住，爱永远不会消失，它只是以另一种形式存在。', '2025-01-20 11:03:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%爱永远不会消失%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '你说得对，我会慢慢好起来的', '2025-01-20 11:04:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='你说得对，我会慢慢好起来的')
UNION ALL SELECT u.id, 1, 0, 0, 0, '是的，时间会治愈一切。在这个过程中，我会一直陪着你，给你力量和勇气。', '2025-01-20 11:05:00' FROM users u WHERE u.username='healer002' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%给你力量和勇气%');

-- =========================
-- 6. 生成小纸条解锁记录
-- =========================
INSERT INTO `user_unlock_records`(`user_id`,`message_id`,`coins_spent`,`created_at`)
SELECT u.id, m.id, m.unlock_coins, DATE_ADD(m.created_at, INTERVAL 1 HOUR)
FROM users u JOIN messages m ON u.id = m.user_id
WHERE m.message_type = 1 AND m.is_locked = 1 AND m.unlock_coins > 0
AND NOT EXISTS(SELECT 1 FROM user_unlock_records ur WHERE ur.user_id = u.id AND ur.message_id = m.id);

-- =========================
-- 7. 生成评论数据
-- =========================
INSERT INTO `comments`(`post_id`,`user_id`,`content`,`liked_count`,`created_at`)
SELECT p.id, u.id, '太治愈了！', FLOOR(RAND() * 10), DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM posts p, users u
WHERE u.username IN ('angel001','healer002','memory003','guardian004','warm005')
AND NOT EXISTS(SELECT 1 FROM comments c WHERE c.post_id = p.id AND c.user_id = u.id AND c.content = '太治愈了！')
LIMIT 50;

INSERT INTO `comments`(`post_id`,`user_id`,`content`,`liked_count`,`created_at`)
SELECT p.id, u.id, '好可爱啊！', FLOOR(RAND() * 8), DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM posts p, users u
WHERE u.username IN ('happy006','gentle007','loyal008','healcat009','companion010')
AND NOT EXISTS(SELECT 1 FROM comments c WHERE c.post_id = p.id AND c.user_id = u.id AND c.content = '好可爱啊！')
LIMIT 50;

-- =========================
-- 8. 生成点赞数据
-- =========================
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 0, p.id, DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM users u, posts p
WHERE NOT EXISTS(SELECT 1 FROM likes l WHERE l.user_id = u.id AND l.target_type = 0 AND l.target_id = p.id)
LIMIT 200;

INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 1, c.id, DATE_ADD(c.created_at, INTERVAL FLOOR(RAND() * 12) HOUR)
FROM users u, comments c
WHERE NOT EXISTS(SELECT 1 FROM likes l WHERE l.user_id = u.id AND l.target_type = 1 AND l.target_id = c.id)
LIMIT 100;

-- =========================
-- 9. 生成关注关系
-- =========================
INSERT INTO `user_follows`(`follower_id`,`followee_id`,`created_at`)
SELECT u1.id, u2.id, '2025-01-15 10:00:00'
FROM users u1, users u2
WHERE u1.username IN ('angel001','healer002','memory003','guardian004','warm005')
AND u2.username IN ('happy006','gentle007','loyal008','healcat009','companion010')
AND u1.id <> u2.id
AND NOT EXISTS(SELECT 1 FROM user_follows f WHERE f.follower_id = u1.id AND f.followee_id = u2.id)
LIMIT 30;

SET FOREIGN_KEY_CHECKS = 1;