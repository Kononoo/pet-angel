package biz

import (
	"context"
	"time"
)

// Avatar 虚拟形象实体
type Avatar struct {
	AvatarID        int64     `json:"avatar_id"`
	Name            string    `json:"name"`
	ResourcePath    string    `json:"resource_path"`
	IdleAnimation   string    `json:"idle_animation"`
	SwitchAnimation string    `json:"switch_animation"`
	SortOrder       int32     `json:"sort_order"`
	IsDefault       bool      `json:"is_default"`
	IsUnlocked      bool      `json:"is_unlocked"`
	CreatedAt       time.Time `json:"created_at"`
}

// UserAvatar 用户当前形象
type UserAvatar struct {
	UserID         int64     `json:"user_id"`
	AvatarID       int64     `json:"avatar_id"`
	LastSwitchTime time.Time `json:"last_switch_time"`
	Avatar         *Avatar   `json:"avatar"`
}

// Prop 道具实体
type Prop struct {
	PropID            int64  `json:"prop_id"`
	Name              string `json:"name"`
	CategoryID        int64  `json:"category_id"`
	IconPath          string `json:"icon_path"`
	CoinCost          int32  `json:"coin_cost"`
	EffectDescription string `json:"effect_description"`
	UserQuantity      int32  `json:"user_quantity"`
}

// PropCategory 道具分类
type PropCategory struct {
	CategoryID int64   `json:"category_id"`
	Name       string  `json:"name"`
	PropIDs    []int64 `json:"prop_ids"`
	SortOrder  int32   `json:"sort_order"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	MessageID   int64     `json:"message_id"`
	UserID      int64     `json:"user_id"`
	Sender      string    `json:"sender"` // user, avatar
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"` // text, prop_use, note
	RelatedID   int64     `json:"related_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// AvatarRepo 虚拟形象仓储接口
type AvatarRepo interface {
	GetAvatarList(ctx context.Context) ([]*Avatar, error)
	GetAvatarByID(ctx context.Context, avatarID int64) (*Avatar, error)
	GetUserCurrentAvatar(ctx context.Context, userID int64) (*UserAvatar, error)
	SwitchAvatar(ctx context.Context, userID, avatarID int64) error
	CreateAvatar(ctx context.Context, avatar *Avatar) error
}

// PropRepo 道具仓储接口
type PropRepo interface {
	GetPropList(ctx context.Context, categoryID int64) ([]*Prop, error)
	GetPropByID(ctx context.Context, propID int64) (*Prop, error)
	GetPropCategories(ctx context.Context) ([]*PropCategory, error)
	GetUserPropQuantity(ctx context.Context, userID, propID int64) (int32, error)
	UseProp(ctx context.Context, userID, propID int64) error
	AddUserProp(ctx context.Context, userID, propID int64, quantity int32) error
}

// ChatMessageRepo 聊天消息仓储接口
type ChatMessageRepo interface {
	GetChatHistory(ctx context.Context, userID int64, page, pageSize int32) ([]*ChatMessage, int32, error)
	CreateMessage(ctx context.Context, message *ChatMessage) error
}

// AvatarUsecase 虚拟形象用例
type AvatarUsecase struct {
	avatarRepo      AvatarRepo
	propRepo        PropRepo
	chatMessageRepo ChatMessageRepo
	userRepo        UserRepo
}

// NewAvatarUsecase 创建虚拟形象用例
func NewAvatarUsecase(avatarRepo AvatarRepo, propRepo PropRepo, chatMessageRepo ChatMessageRepo, userRepo UserRepo) *AvatarUsecase {
	return &AvatarUsecase{
		avatarRepo:      avatarRepo,
		propRepo:        propRepo,
		chatMessageRepo: chatMessageRepo,
		userRepo:        userRepo,
	}
}

// GetAvatarList 获取虚拟形象列表
func (uc *AvatarUsecase) GetAvatarList(ctx context.Context, userID int64) ([]*Avatar, error) {
	avatars, err := uc.avatarRepo.GetAvatarList(ctx)
	if err != nil {
		return nil, err
	}

	// 检查用户是否解锁了每个形象
	for _, avatar := range avatars {
		// 预设形象默认解锁
		if avatar.IsDefault {
			avatar.IsUnlocked = true
		} else {
			// 这里可以根据业务逻辑判断是否解锁
			// 比如检查用户是否拥有该形象
			avatar.IsUnlocked = true // 简化处理
		}
	}

	return avatars, nil
}

// GetCurrentAvatar 获取用户当前形象
func (uc *AvatarUsecase) GetCurrentAvatar(ctx context.Context, userID int64) (*UserAvatar, error) {
	userAvatar, err := uc.avatarRepo.GetUserCurrentAvatar(ctx, userID)
	if err != nil {
		// 如果没有设置当前形象，使用默认形象
		defaultAvatar, err := uc.avatarRepo.GetAvatarByID(ctx, 1) // 假设ID为1是默认形象
		if err != nil {
			return nil, err
		}

		userAvatar = &UserAvatar{
			UserID:         userID,
			AvatarID:       defaultAvatar.AvatarID,
			LastSwitchTime: time.Now(),
			Avatar:         defaultAvatar,
		}

		// 设置默认形象
		uc.avatarRepo.SwitchAvatar(ctx, userID, defaultAvatar.AvatarID)
	}

	return userAvatar, nil
}

// SwitchAvatar 切换虚拟形象
func (uc *AvatarUsecase) SwitchAvatar(ctx context.Context, userID, avatarID int64) (*UserAvatar, error) {
	// 检查形象是否存在
	avatar, err := uc.avatarRepo.GetAvatarByID(ctx, avatarID)
	if err != nil {
		return nil, ErrAvatarNotFound
	}

	// 检查用户是否解锁该形象
	if !avatar.IsDefault {
		// 这里可以添加解锁检查逻辑
		// 暂时简化处理
	}

	// 切换形象
	err = uc.avatarRepo.SwitchAvatar(ctx, userID, avatarID)
	if err != nil {
		return nil, err
	}

	userAvatar := &UserAvatar{
		UserID:         userID,
		AvatarID:       avatarID,
		LastSwitchTime: time.Now(),
		Avatar:         avatar,
	}

	return userAvatar, nil
}

// GetPropList 获取道具列表
func (uc *AvatarUsecase) GetPropList(ctx context.Context, userID, categoryID int64) ([]*Prop, error) {
	props, err := uc.propRepo.GetPropList(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// 获取用户持有数量
	for _, prop := range props {
		quantity, err := uc.propRepo.GetUserPropQuantity(ctx, userID, prop.PropID)
		if err == nil {
			prop.UserQuantity = quantity
		}
	}

	return props, nil
}

// GetPropCategories 获取道具分类
func (uc *AvatarUsecase) GetPropCategories(ctx context.Context) ([]*PropCategory, error) {
	return uc.propRepo.GetPropCategories(ctx)
}

// UseProp 使用道具
func (uc *AvatarUsecase) UseProp(ctx context.Context, userID, propID int64) error {
	// 检查道具是否存在
	prop, err := uc.propRepo.GetPropByID(ctx, propID)
	if err != nil {
		return ErrPropNotFound
	}

	// 检查用户是否拥有该道具
	quantity, err := uc.propRepo.GetUserPropQuantity(ctx, userID, propID)
	if err != nil || quantity <= 0 {
		return ErrPropNotOwned
	}

	// 使用道具
	err = uc.propRepo.UseProp(ctx, userID, propID)
	if err != nil {
		return err
	}

	// 记录聊天消息
	message := &ChatMessage{
		UserID:      userID,
		Sender:      "user",
		Content:     "使用了" + prop.Name,
		MessageType: "prop_use",
		RelatedID:   propID,
		CreatedAt:   time.Now(),
	}

	uc.chatMessageRepo.CreateMessage(ctx, message)

	return nil
}

// GetChatHistory 获取聊天记录
func (uc *AvatarUsecase) GetChatHistory(ctx context.Context, userID int64, page, pageSize int32) ([]*ChatMessage, int32, error) {
	return uc.chatMessageRepo.GetChatHistory(ctx, userID, page, pageSize)
}

// SendMessage 发送消息
func (uc *AvatarUsecase) SendMessage(ctx context.Context, userID int64, content string) (*ChatMessage, error) {
	message := &ChatMessage{
		UserID:      userID,
		Sender:      "user",
		Content:     content,
		MessageType: "text",
		CreatedAt:   time.Now(),
	}

	err := uc.chatMessageRepo.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GenerateAvatar 生成虚拟形象
func (uc *AvatarUsecase) GenerateAvatar(ctx context.Context, userID int64, imageURL, petName, petSpecies string) (*Avatar, error) {
	// 这里可以调用AI服务生成虚拟形象
	// 简化处理，创建一个新的形象记录

	avatar := &Avatar{
		Name:         petName + "的虚拟形象",
		ResourcePath: imageURL, // 这里应该是AI生成的资源路径
		SortOrder:    999,
		IsDefault:    false,
		IsUnlocked:   true,
		CreatedAt:    time.Now(),
	}

	err := uc.avatarRepo.CreateAvatar(ctx, avatar)
	if err != nil {
		return nil, err
	}

	return avatar, nil
}
