package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	msgv1 "pet-angel/api/message/v1"
	"pet-angel/internal/ai"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// MessageService 提供消息/小纸条相关接口

type MessageService struct {
	msgv1.UnimplementedMessageServiceServer
	uc        *biz.MessageUsecase
	jwtSecret string
	logger    *log.Helper
}

func NewMessageService(uc *biz.MessageUsecase, cfg *conf.Auth, l log.Logger) *MessageService {
	secret := ""
	if cfg != nil {
		secret = cfg.JwtSecret
	}
	return &MessageService{uc: uc, jwtSecret: secret, logger: log.NewHelper(l)}
}

// userIDFromCtx 解析 JWT
func (s *MessageService) userIDFromCtx(ctx context.Context) (int64, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return 0, http.ErrNoCookie
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return 0, err
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetMessageList 列表
func (s *MessageService) GetMessageList(ctx context.Context, in *msgv1.GetMessageListRequest) (*msgv1.GetMessageListReply, error) {
	userID, err := s.userIDFromCtx(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get messages: parse auth failed: %v", err)
		return nil, err
	}
	total, list, err := s.uc.GetList(ctx, userID, in.GetOnlyNotes(), in.GetPage(), in.GetPageSize())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get messages: usecase error: %v", err)
		return nil, err
	}
	out := make([]*msgv1.Message, 0, len(list))
	for _, m := range list {
		mm := &msgv1.Message{
			Id:          m.ID,
			Sender:      m.Sender,
			MessageType: m.MessageType,
			IsLocked:    m.IsLocked,
			UnlockCoins: m.UnlockCoins,
			Content:     "",
			CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if !m.IsLocked {
			mm.Content = m.Content
		}
		out = append(out, mm)
	}
	return &msgv1.GetMessageListReply{Total: total, List: out}, nil
}

// UnlockMessage 解锁
func (s *MessageService) UnlockMessage(ctx context.Context, in *msgv1.UnlockMessageRequest) (*msgv1.UnlockMessageReply, error) {
	userID, err := s.userIDFromCtx(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("unlock message: parse auth failed: %v", err)
		return nil, err
	}
	remain, m, err := s.uc.Unlock(ctx, userID, in.GetMessageId())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("unlock message: usecase error: %v", err)
		return nil, err
	}
	reply := &msgv1.UnlockMessageReply{
		Success:        true,
		RemainingCoins: remain,
		Message: &msgv1.Message{
			Id:          m.ID,
			Sender:      m.Sender,
			MessageType: m.MessageType,
			IsLocked:    m.IsLocked,
			UnlockCoins: m.UnlockCoins,
			Content:     m.Content,
			CreatedAt:   m.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	return reply, nil
}

// 预留：生成当日小纸条（仅用于脚本/后台触发，业务端暂不开放 HTTP）
// 说明：这里不暴露接口，后续可在 service/script 或 job 中调用 usecase.CreateDailyNotes

// ParseSecretFromCtx 辅助：若没有全局 secret，则用默认实现
// 我们实现一个简单包装，在 jwt 包中提供

func init() {
	// 确保 time 包被引用（用于格式化）
	_ = time.Now()
}

// GenerateNotesHTTP 暂供演示/内部工具：根据当前登录用户生成今日小纸条（3条免费+1条20金币）
func (s *MessageService) GenerateNotesHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// 鉴权
		authHeader := r.Header.Get("Authorization")
		tok, err := jwtutil.FromAuthHeader(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, err := jwtutil.Parse(s.jwtSecret, tok)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		// 生成 4 条文案（调用统一 AI）
		client := ai.Default()
		if client == nil {
			client = ai.NewClient(ai.Config{})
		}
		prompts := []string{
			"请用20字以内写一条温柔的每日问候。",
			"请用20字以内写一条积极的鼓励语。",
			"请用20字以内写一条与宠物相关的暖心提醒。",
			"请用20字以内写一条晚上温柔的陪伴句子。",
		}
		notes := make([]struct {
			Coins   int32
			Content string
		}, 0, 4)
		for i, p := range prompts {
			txt, _ := client.Chat(r.Context(), "你是治愈系宠物数字伙伴，第一人称‘我’，中文简短。", p)
			coin := int32(0)
			if i == 3 {
				coin = 20
			}
			notes = append(notes, struct {
				Coins   int32
				Content string
			}{Coins: coin, Content: txt})
		}
		if err := s.uc.CreateDailyNotes(r.Context(), userID, notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "count": len(notes)})
	}
}
