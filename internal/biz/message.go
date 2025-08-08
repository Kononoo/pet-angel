package biz

import (
	"context"
	"time"
)

// Message 小纸条实体
type Message struct {
	MessageID     int64     `json:"message_id"`
	Content       string    `json:"content"`
	MessageType   string    `json:"message_type"` // free, paid
	UnlockCoins   int32     `json:"unlock_coins"`
	PetID         int64     `json:"pet_id"`
	IsUnlocked    bool      `json:"is_unlocked"`
	UnlockTime    time.Time `json:"unlock_time"`
	IsAIGenerated bool      `json:"is_ai_generated"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserCoins 用户金币信息
type UserCoins struct {
	UserID      int64     `json:"user_id"`
	CoinBalance int32     `json:"coin_balance"`
	TotalEarned int32     `json:"total_earned"`
	TotalSpent  int32     `json:"total_spent"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UnlockRecord 解锁记录
type UnlockRecord struct {
	RecordID   int64     `json:"record_id"`
	UserID     int64     `json:"user_id"`
	MessageID  int64     `json:"message_id"`
	CoinsSpent int32     `json:"coins_spent"`
	UnlockTime time.Time `json:"unlock_time"`
	Message    *Message  `json:"message"`
}

// MessageRepo 小纸条仓储接口
type MessageRepo interface {
	GetMessageList(ctx context.Context, messageType string, page, pageSize int32) ([]*Message, int32, error)
	GetMessageByID(ctx context.Context, messageID int64) (*Message, error)
	CreateMessage(ctx context.Context, message *Message) error
	GetUserUnlockRecord(ctx context.Context, userID, messageID int64) (*UnlockRecord, error)
	CreateUnlockRecord(ctx context.Context, record *UnlockRecord) error
}

// UserCoinsRepo 用户金币仓储接口
type UserCoinsRepo interface {
	GetUserCoins(ctx context.Context, userID int64) (*UserCoins, error)
	UpdateUserCoins(ctx context.Context, userID int64, delta int32) error
	GetUnlockHistory(ctx context.Context, userID int64, page, pageSize int32) ([]*UnlockRecord, int32, error)
}

// MessagePetRepo 宠物仓储接口（用于检查宠物状态）
type MessagePetRepo interface {
	GetPetByID(ctx context.Context, petID int64) (*Pet, error)
}

// MessageUsecase 小纸条用例
type MessageUsecase struct {
	messageRepo   MessageRepo
	userCoinsRepo UserCoinsRepo
	petRepo       MessagePetRepo
}

// NewMessageUsecase 创建小纸条用例
func NewMessageUsecase(messageRepo MessageRepo, userCoinsRepo UserCoinsRepo, petRepo MessagePetRepo) *MessageUsecase {
	return &MessageUsecase{
		messageRepo:   messageRepo,
		userCoinsRepo: userCoinsRepo,
		petRepo:       petRepo,
	}
}

// GetMessageList 获取小纸条列表
func (uc *MessageUsecase) GetMessageList(ctx context.Context, messageType string, page, pageSize int32, userID int64) ([]*Message, int32, error) {
	messages, total, err := uc.messageRepo.GetMessageList(ctx, messageType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 检查用户解锁状态
	for _, message := range messages {
		record, err := uc.messageRepo.GetUserUnlockRecord(ctx, userID, message.MessageID)
		if err == nil && record != nil {
			message.IsUnlocked = true
			message.UnlockTime = record.UnlockTime
		}
	}

	return messages, total, nil
}

// GetMessageDetail 获取小纸条详情
func (uc *MessageUsecase) GetMessageDetail(ctx context.Context, messageID int64, userID int64) (*Message, error) {
	message, err := uc.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// 检查用户解锁状态
	record, err := uc.messageRepo.GetUserUnlockRecord(ctx, userID, messageID)
	if err == nil && record != nil {
		message.IsUnlocked = true
		message.UnlockTime = record.UnlockTime
	}

	return message, nil
}

// UnlockMessage 解锁小纸条
func (uc *MessageUsecase) UnlockMessage(ctx context.Context, userID, messageID int64) error {
	// 获取小纸条信息
	message, err := uc.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return ErrMessageNotFound
	}

	// 检查是否已经解锁
	record, err := uc.messageRepo.GetUserUnlockRecord(ctx, userID, messageID)
	if err == nil && record != nil {
		return nil // 已经解锁
	}

	// 检查宠物状态（如果是付费小纸条且关联特定宠物）
	if message.PetID > 0 && message.MessageType == "paid" {
		pet, err := uc.petRepo.GetPetByID(ctx, message.PetID)
		if err != nil {
			return ErrPetNotFound
		}

		// 检查宠物是否已去世（某些小纸条只在宠物去世后解锁）
		if pet.PassedAwayDate.IsZero() {
			return ErrMessageLocked
		}
	}

	// 获取用户金币信息
	userCoins, err := uc.userCoinsRepo.GetUserCoins(ctx, userID)
	if err != nil {
		return err
	}

	// 检查金币是否足够
	if userCoins.CoinBalance < message.UnlockCoins {
		return ErrInsufficientCoins
	}

	// 扣除金币
	err = uc.userCoinsRepo.UpdateUserCoins(ctx, userID, -message.UnlockCoins)
	if err != nil {
		return err
	}

	// 创建解锁记录
	unlockRecord := &UnlockRecord{
		UserID:     userID,
		MessageID:  messageID,
		CoinsSpent: message.UnlockCoins,
		UnlockTime: time.Now(),
		Message:    message,
	}

	err = uc.messageRepo.CreateUnlockRecord(ctx, unlockRecord)
	if err != nil {
		return err
	}

	return nil
}

// GenerateMessage 生成小纸条
func (uc *MessageUsecase) GenerateMessage(ctx context.Context, petID int64, prompt string) (*Message, error) {
	// 这里可以调用AI服务生成小纸条内容
	// 简化处理，使用固定内容
	content := "亲爱的主人，我想你了！今天过得怎么样？"

	message := &Message{
		Content:       content,
		MessageType:   "free", // 默认免费
		UnlockCoins:   0,
		PetID:         petID,
		IsAIGenerated: true,
		CreatedAt:     time.Now(),
	}

	err := uc.messageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetUserCoins 获取用户金币信息
func (uc *MessageUsecase) GetUserCoins(ctx context.Context, userID int64) (*UserCoins, error) {
	return uc.userCoinsRepo.GetUserCoins(ctx, userID)
}

// GetUnlockHistory 获取解锁记录
func (uc *MessageUsecase) GetUnlockHistory(ctx context.Context, userID int64, page, pageSize int32) ([]*UnlockRecord, int32, error) {
	return uc.userCoinsRepo.GetUnlockHistory(ctx, userID, page, pageSize)
}
