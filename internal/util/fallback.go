package util

// GetFallbackReply 根据用户输入生成兜底回复
func GetFallbackReply(userContent string) string {
	// 根据用户输入的关键词生成不同的兜底回复
	if Contains(userContent, "你好") || Contains(userContent, "hello") || Contains(userContent, "hi") {
		return "喵~ 你好呀！我是你的小天使，今天想和我聊什么呢？"
	}
	if Contains(userContent, "心情") || Contains(userContent, "难过") || Contains(userContent, "伤心") {
		return "我感受到了你的心情，让我用毛茸茸的小爪子给你一个温暖的抱抱~"
	}
	if Contains(userContent, "累") || Contains(userContent, "疲惫") || Contains(userContent, "困") {
		return "累了就休息一下吧，我会在这里陪着你，给你最温暖的守护~"
	}
	if Contains(userContent, "谢谢") || Contains(userContent, "感谢") {
		return "不用谢哦~ 能陪着你是我最大的幸福，我会一直在这里守护你~"
	}
	if Contains(userContent, "再见") || Contains(userContent, "拜拜") || Contains(userContent, "bye") {
		return "再见啦~ 记得想我哦，我会一直在这里等你回来~"
	}
	// 默认回复
	return "我在呢，会一直陪着你~ 有什么想和我分享的吗？"
}

// Contains 检查字符串是否包含子串
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
