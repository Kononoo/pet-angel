-- Pet Angel 补充海量数据
-- 为其他用户生成更多小纸条、聊天记录和帖子

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =========================
-- 为其他用户生成小纸条（每个用户15+ 小纸条）
-- =========================
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
-- memory003 的小纸条
SELECT u.id, 1, 1, 0, 0, '早安，回忆录！今天也要记录美好。', '2025-02-01 08:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%回忆录%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一个回忆都是珍贵的宝藏。', '2025-02-02 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%珍贵的宝藏%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！记录今天的温暖瞬间。', '2025-02-03 12:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%温暖瞬间%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '时间会带走一切，但带不走我们的回忆。', '2025-02-04 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%带不走我们的回忆%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，愿梦里有美好的回忆。', '2025-02-05 22:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%美好的回忆%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的回忆，是这个世界最美的故事。', '2025-02-06 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最美的故事%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！新的一天，新的回忆。', '2025-02-07 08:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%新的回忆%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次记录，都是对生命的礼赞。', '2025-02-08 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对生命的礼赞%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！记得记录今天的感动。', '2025-02-09 12:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%记录今天的感动%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '回忆是最美的风景，你是最好的记录者。', '2025-02-10 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最好的记录者%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，明天继续记录美好。', '2025-02-11 22:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%继续记录美好%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的回忆录，是治愈的良药。', '2025-02-12 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈的良药%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！今天也要记录精彩。', '2025-02-13 08:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%记录精彩%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次回忆，都是对爱的延续。', '2025-02-14 20:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对爱的延续%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的回忆在发光。', '2025-02-15 12:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%回忆在发光%');

-- guardian004 的小纸条
SELECT u.id, 1, 1, 0, 0, '早安，守护者！今天也要守护美好。', '2025-02-10 08:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%守护者%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的守护，是这个世界的力量。', '2025-02-11 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%这个世界的力量%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！守护的使命在召唤。', '2025-02-12 12:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%守护的使命%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次守护，都是对生命的承诺。', '2025-02-13 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对生命的承诺%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，守护的力量永存。', '2025-02-14 22:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%守护的力量%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的忠诚，是这个世界的光。', '2025-02-15 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%这个世界的光%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！新的一天，新的守护。', '2025-02-16 08:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%新的守护%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '守护是最美的誓言。', '2025-02-17 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最美的誓言%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的守护让世界更安全。', '2025-02-18 12:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%让世界更安全%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次守护，都是对爱的诠释。', '2025-02-19 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对爱的诠释%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '晚安，明天继续守护。', '2025-02-20 22:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%继续守护%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '你的守护，是治愈的良药。', '2025-02-21 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%治愈的良药%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '早安！今天也要守护如初。', '2025-02-22 08:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%守护如初%')
UNION ALL SELECT u.id, 1, 1, 1, 20, '每一次守护，都是对生命的礼赞。', '2025-02-23 20:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%对生命的礼赞%')
UNION ALL SELECT u.id, 1, 1, 0, 0, '午安！你的光芒正在守护。', '2025-02-24 12:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%光芒正在守护%');

-- =========================
-- 为其他用户生成聊天记录
-- =========================
INSERT INTO `messages`(`user_id`,`sender`,`message_type`,`is_locked`,`unlock_coins`,`content`,`created_at`)
-- memory003 的聊天记录
SELECT u.id, 0, 0, 0, 0, '你好，回忆录', '2025-02-01 12:00:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='你好，回忆录')
UNION ALL SELECT u.id, 1, 0, 0, 0, '你好！我是你的专属回忆录，让我们一起记录生活中的每一个美好瞬间吧！', '2025-02-01 12:01:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%专属回忆录%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '我想记录一些美好的回忆', '2025-02-01 12:02:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='我想记录一些美好的回忆')
UNION ALL SELECT u.id, 1, 0, 0, 0, '太好了！回忆是最珍贵的财富。无论是快乐的时光还是感动的瞬间，都值得被记录下来。', '2025-02-01 12:03:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%最珍贵的财富%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '谢谢你，我会好好记录的', '2025-02-01 12:04:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='谢谢你，我会好好记录的')
UNION ALL SELECT u.id, 1, 0, 0, 0, '记住，每一个回忆都是独一无二的。我会一直陪着你，记录每一个值得珍藏的瞬间。', '2025-02-01 12:05:00' FROM users u WHERE u.username='memory003' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%独一无二的%');

-- guardian004 的聊天记录
SELECT u.id, 0, 0, 0, 0, '你好，守护者', '2025-02-10 13:00:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='你好，守护者')
UNION ALL SELECT u.id, 1, 0, 0, 0, '你好！我是你的专属守护者，我会一直守护着你，保护你的安全和幸福！', '2025-02-10 13:01:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%专属守护者%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '我需要一些守护的力量', '2025-02-10 13:02:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='我需要一些守护的力量')
UNION ALL SELECT u.id, 1, 0, 0, 0, '守护的力量就在你心中。我会一直陪着你，给你勇气和力量，守护你的每一个梦想。', '2025-02-10 13:03:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%守护的力量%')
UNION ALL SELECT u.id, 0, 0, 0, 0, '谢谢你，我感觉更有力量了', '2025-02-10 13:04:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content='谢谢你，我感觉更有力量了')
UNION ALL SELECT u.id, 1, 0, 0, 0, '记住，你永远不是一个人在战斗。我会一直守护着你，直到永远。', '2025-02-10 13:05:00' FROM users u WHERE u.username='guardian004' AND NOT EXISTS(SELECT 1 FROM messages WHERE user_id=u.id AND content LIKE '%永远不是一个人%');

-- =========================
-- 生成更多帖子（批量生成）
-- =========================
INSERT INTO `posts`(`user_id`,`category_id`,`title`,`content`,`type`,`image_urls`,`video_url`,`cover_url`,`locate`,`tags`,`liked_count`,`comment_count`,`is_private`,`created_at`,`updated_at`)
-- 为其他用户生成更多帖子
SELECT u.id, c.id, CONCAT('温暖时光第', n.n, '天'), CONCAT('第', n.n, '天的温暖时光，每一天都充满爱。'), 0,
       CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/warm_time_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/warm_time_', LPAD(n.n, 3, '0'), '_cover.jpg'), '广州','温暖,时光', FLOOR(40 + RAND() * 120), FLOOR(3 + RAND() * 15), 0, DATE_ADD('2025-02-15 14:30:00', INTERVAL n.n DAY), DATE_ADD('2025-02-15 14:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='warm005' AND c.name='治愈日常' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('温暖时光第', n.n, '天'))

UNION ALL SELECT u.id, c.id, CONCAT('快乐奔跑第', n.n, '次'), CONCAT('第', n.n, '次快乐奔跑，每一次都充满活力。'), 0,
       CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/happy_run_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/02/', LPAD(n.n, 2, '0'), '/happy_run_', LPAD(n.n, 3, '0'), '_cover.jpg'), '成都','快乐,奔跑', FLOOR(70 + RAND() * 200), FLOOR(10 + RAND() * 30), 0, DATE_ADD('2025-02-20 15:30:00', INTERVAL n.n DAY), DATE_ADD('2025-02-20 15:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='happy006' AND c.name='日常' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('快乐奔跑第', n.n, '次'))

UNION ALL SELECT u.id, c.id, CONCAT('温柔陪伴第', n.n, '天'), CONCAT('第', n.n, '天的温柔陪伴，每一天都充满爱。'), 0,
       CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/gentle_company_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/gentle_company_', LPAD(n.n, 3, '0'), '_cover.jpg'), '武汉','温柔,陪伴', FLOOR(80 + RAND() * 150), FLOOR(12 + RAND() * 25), 0, DATE_ADD('2025-03-01 16:30:00', INTERVAL n.n DAY), DATE_ADD('2025-03-01 16:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='gentle007' AND c.name='情感' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('温柔陪伴第', n.n, '天'))

UNION ALL SELECT u.id, c.id, CONCAT('忠诚守护第', n.n, '天'), CONCAT('第', n.n, '天的忠诚守护，每一天都充满责任。'), 0,
       CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/loyal_guard_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/loyal_guard_', LPAD(n.n, 3, '0'), '_cover.jpg'), '南京','忠诚,守护', FLOOR(60 + RAND() * 130), FLOOR(8 + RAND() * 20), 0, DATE_ADD('2025-03-05 17:30:00', INTERVAL n.n DAY), DATE_ADD('2025-03-05 17:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='loyal008' AND c.name='情感互助' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('忠诚守护第', n.n, '天'))

UNION ALL SELECT u.id, c.id, CONCAT('智慧治愈第', n.n, '天'), CONCAT('第', n.n, '天的智慧治愈，每一天都充满智慧。'), 0,
       CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/wise_heal_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/wise_heal_', LPAD(n.n, 3, '0'), '_cover.jpg'), '西安','智慧,治愈', FLOOR(50 + RAND() * 140), FLOOR(6 + RAND() * 18), 0, DATE_ADD('2025-03-10 18:30:00', INTERVAL n.n DAY), DATE_ADD('2025-03-10 18:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='healcat009' AND c.name='知识' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('智慧治愈第', n.n, '天'))

UNION ALL SELECT u.id, c.id, CONCAT('聪明训练第', n.n, '天'), CONCAT('第', n.n, '天的聪明训练，每一天都充满挑战。'), 0,
       CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/smart_train_', LPAD(n.n, 3, '0'), '.jpg'), NULL, CONCAT('/static/image/2025/03/', LPAD(n.n, 2, '0'), '/smart_train_', LPAD(n.n, 3, '0'), '_cover.jpg'), '重庆','聪明,训练', FLOOR(90 + RAND() * 180), FLOOR(15 + RAND() * 35), 0, DATE_ADD('2025-03-15 19:30:00', INTERVAL n.n DAY), DATE_ADD('2025-03-15 19:30:00', INTERVAL n.n DAY)
FROM users u, categories c, (SELECT 1 n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5 UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10) n
WHERE u.username='companion010' AND c.name='训练' AND NOT EXISTS(SELECT 1 FROM posts WHERE title=CONCAT('聪明训练第', n.n, '天'));

-- =========================
-- 生成更多评论和点赞
-- =========================
INSERT INTO `comments`(`post_id`,`user_id`,`content`,`liked_count`,`created_at`)
SELECT p.id, u.id, '太棒了！', FLOOR(RAND() * 12), DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM posts p, users u
WHERE u.username IN ('memory003','guardian004','warm005','happy006','gentle007','loyal008','healcat009','companion010')
AND NOT EXISTS(SELECT 1 FROM comments c WHERE c.post_id = p.id AND c.user_id = u.id AND c.content = '太棒了！')
LIMIT 100;

INSERT INTO `comments`(`post_id`,`user_id`,`content`,`liked_count`,`created_at`)
SELECT p.id, u.id, '好喜欢！', FLOOR(RAND() * 10), DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM posts p, users u
WHERE u.username IN ('angel001','healer002','memory003','guardian004','warm005','happy006','gentle007','loyal008','healcat009','companion010')
AND NOT EXISTS(SELECT 1 FROM comments c WHERE c.post_id = p.id AND c.user_id = u.id AND c.content = '好喜欢！')
LIMIT 100;

-- 生成更多点赞
INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 0, p.id, DATE_ADD(p.created_at, INTERVAL FLOOR(RAND() * 24) HOUR)
FROM users u, posts p
WHERE NOT EXISTS(SELECT 1 FROM likes l WHERE l.user_id = u.id AND l.target_type = 0 AND l.target_id = p.id)
LIMIT 300;

INSERT INTO `likes`(`user_id`,`target_type`,`target_id`,`created_at`)
SELECT u.id, 1, c.id, DATE_ADD(c.created_at, INTERVAL FLOOR(RAND() * 12) HOUR)
FROM users u, comments c
WHERE NOT EXISTS(SELECT 1 FROM likes l WHERE l.user_id = u.id AND l.target_type = 1 AND l.target_id = c.id)
LIMIT 200;

-- =========================
-- 生成更多关注关系
-- =========================
INSERT INTO `user_follows`(`follower_id`,`followee_id`,`created_at`)
SELECT u1.id, u2.id, DATE_ADD('2025-01-15 10:00:00', INTERVAL FLOOR(RAND() * 30) DAY)
FROM users u1, users u2
WHERE u1.id <> u2.id
AND NOT EXISTS(SELECT 1 FROM user_follows f WHERE f.follower_id = u1.id AND f.followee_id = u2.id)
LIMIT 50;

SET FOREIGN_KEY_CHECKS = 1; 