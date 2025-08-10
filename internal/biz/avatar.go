package biz

import (
	"context"
	"time"
)

// PetModel 业务实体（对应表 pet_models）
// 使用整数与字符串，避免 JSON/ENUM；时间统一为 datetime 字符串在 service 层格式化
type PetModel struct {
	ID        int64     // 模型ID
	Name      string    // 模型名称
	Path      string    // 资源URL
	ModelType int32     // 0=猫 1=狗
	IsDefault bool      // 是否默认
	SortOrder int32     // 排序
	CreatedAt time.Time // 创建时间
}

// Item 业务实体（对应表 items）
type Item struct {
	ID          int64     // 道具ID
	Name        string    // 名称
	Description string    // 描述
	IconPath    string    // 图标URL
	CoinCost    int32     // 消耗金币
	CreatedAt   time.Time // 创建时间
}

// Message 业务实体（对应表 messages）
type ChatMsg struct {
	ID          int64     // 消息ID
	UserID      int64     // 用户ID
	Sender      int32     // 0用户 1AI
	MessageType int32     // 0聊天 1小纸条
	IsLocked    bool      // 锁定
	UnlockCoins int32     // 解锁金币
	Content     string    // 内容
	CreatedAt   time.Time // 创建时间
}

// AvatarRepo 数据仓储接口（GORM 实现）
type AvatarRepo interface {
	ListPetModels(ctx context.Context) ([]*PetModel, error)
	PetModelExists(ctx context.Context, modelID int64) (bool, error)
	SetUserModel(ctx context.Context, userID, modelID int64) error

	ListItems(ctx context.Context) ([]*Item, error)
	UseItem(ctx context.Context, userID, itemID int64) (remainingCoins int32, err error)

	CreateChat(ctx context.Context, userID int64, content string) (*ChatMsg, error)
	CreateAIChat(ctx context.Context, userID int64, content string) (*ChatMsg, error)
	// 直接写入一条 AI 消息（用于流式完成后落库）
	CreateAIMessage(ctx context.Context, userID int64, content string) (*ChatMsg, error)
	// 获取最新的AI消息
	GetLatestAIMessage(ctx context.Context, userID int64) (*ChatMsg, error)
}

// AvatarUsecase 业务用例
type AvatarUsecase struct {
	repo AvatarRepo
}

func NewAvatarUsecase(repo AvatarRepo) *AvatarUsecase { return &AvatarUsecase{repo: repo} }

// GetModels 列出所有模型
func (uc *AvatarUsecase) GetModels(ctx context.Context) ([]*PetModel, error) {
	return uc.repo.ListPetModels(ctx)
}

// SetPetModel 设置用户当前模型
func (uc *AvatarUsecase) SetPetModel(ctx context.Context, userID, modelID int64) error {
	exists, err := uc.repo.PetModelExists(ctx, modelID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAvatarNotFound
	}
	return uc.repo.SetUserModel(ctx, userID, modelID)
}

// GetItems 列出所有道具
func (uc *AvatarUsecase) GetItems(ctx context.Context) ([]*Item, error) {
	return uc.repo.ListItems(ctx)
}

// UseItem 使用道具并扣金币
func (uc *AvatarUsecase) UseItem(ctx context.Context, userID, itemID int64) (int32, error) {
	remaining, err := uc.repo.UseItem(ctx, userID, itemID)
	if err != nil {
		return 0, err
	}
	return remaining, nil
}

// Chat 发送聊天消息（同步返回该条消息）
func (uc *AvatarUsecase) Chat(ctx context.Context, userID int64, content string) (*ChatMsg, error) {
	// 1) 先写入用户消息
	userMsg, err := uc.repo.CreateChat(ctx, userID, content)
	if err != nil {
		return nil, err
	}
	// 2) 生成 AI 回复（由 data 层调用 AI 客户端并落库），保持简单直连
	//    返回值可选；业务层只需保证用户消息已记录
	_, _ = uc.repo.CreateAIChat(ctx, userID, content)
	return userMsg, nil
}

// SaveAIMessage 将一段 AI 文本回复直接写库（供流式完成后调用）
func (uc *AvatarUsecase) SaveAIMessage(ctx context.Context, userID int64, content string) (*ChatMsg, error) {
	return uc.repo.CreateAIMessage(ctx, userID, content)
}

// SaveUserMessage 仅保存用户消息（不触发 AI 回复）
func (uc *AvatarUsecase) SaveUserMessage(ctx context.Context, userID int64, content string) (*ChatMsg, error) {
	return uc.repo.CreateChat(ctx, userID, content)
}

// GetLatestAIMessage 获取最新的AI消息
func (uc *AvatarUsecase) GetLatestAIMessage(ctx context.Context, userID int64) (*ChatMsg, error) {
	return uc.repo.GetLatestAIMessage(ctx, userID)
}
