package biz

import (
	"context"
	"time"
)

// Message 业务实体（复用消息表结构）
type Message struct {
	ID          int64     // 消息ID
	UserID      int64     // 用户ID
	Sender      int32     // 0用户 1AI
	MessageType int32     // 0聊天 1小纸条
	IsLocked    bool      // 是否锁定
	UnlockCoins int32     // 解锁所需金币
	Content     string    // 内容
	CreatedAt   time.Time // 创建时间
}

// MessageRepo 数据仓储接口
// ListMessages: 返回总数与列表（倒序分页）
// UnlockMessage: 事务性解锁，扣金币，返回剩余金币与最新消息
// GetMessageByID: 获取单条（用于服务端再取）

type MessageRepo interface {
	ListMessages(ctx context.Context, userID int64, onlyNotes bool, page, pageSize int32) (total int32, list []*Message, err error)
	UnlockMessage(ctx context.Context, userID, messageID int64) (remainingCoins int32, msg *Message, err error)
	GetMessageByID(ctx context.Context, userID, messageID int64) (*Message, error)
}

// MessageUsecase 用例

type MessageUsecase struct {
	repo MessageRepo
}

func NewMessageUsecase(repo MessageRepo) *MessageUsecase { return &MessageUsecase{repo: repo} }

// GetList 获取消息列表
func (uc *MessageUsecase) GetList(ctx context.Context, userID int64, onlyNotes bool, page, pageSize int32) (int32, []*Message, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return uc.repo.ListMessages(ctx, userID, onlyNotes, page, pageSize)
}

// Unlock 解锁小纸条
func (uc *MessageUsecase) Unlock(ctx context.Context, userID, messageID int64) (int32, *Message, error) {
	remaining, msg, err := uc.repo.UnlockMessage(ctx, userID, messageID)
	if err != nil {
		return 0, nil, err
	}
	return remaining, msg, nil
}
