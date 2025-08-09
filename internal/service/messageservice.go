package service

import (
	"context"
	"net/http"
	"time"

	msgv1 "pet-angel/api/message/v1"
	"pet-angel/internal/biz"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// MessageService 提供消息/小纸条相关接口

type MessageService struct {
	msgv1.UnimplementedMessageServiceServer
	uc     *biz.MessageUsecase
	logger *log.Helper
}

func NewMessageService(uc *biz.MessageUsecase, l log.Logger) *MessageService {
	return &MessageService{uc: uc, logger: log.NewHelper(l)}
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
	claims, err := jwtutil.ParseSecretFromCtx(ctx, tok)
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

// ParseSecretFromCtx 辅助：若没有全局 secret，则用默认实现
// 我们实现一个简单包装，在 jwt 包中提供

func init() {
	// 确保 time 包被引用（用于格式化）
	_ = time.Now()
}
