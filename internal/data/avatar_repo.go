package data

import (
	"context"
	"errors"
	"time"

	aiclient "pet-angel/internal/ai"
	"pet-angel/internal/biz"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM 模型定义（表结构与 1-init-tables.sql 对齐）

// PetModelDO 映射 pet_models 表
type PetModelDO struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`     // 模型ID
	Name      string    `gorm:"column:name;type:varchar(100);not null"` // 名称
	Path      string    `gorm:"column:path;type:varchar(255);not null"` // 资源URL
	Type      int32     `gorm:"column:type;not null"`                   // 0猫 1狗
	IsDefault bool      `gorm:"column:is_default;not null"`             // 是否默认
	SortOrder int32     `gorm:"column:sort_order;not null"`             // 排序
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`       // 创建时间
}

func (PetModelDO) TableName() string { return "pet_models" }

// ItemDO 映射 items 表
type ItemDO struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`     // 道具ID
	Name        string    `gorm:"column:name;type:varchar(100);not null"` // 名称
	Description string    `gorm:"column:description;type:varchar(255)"`   // 描述
	IconPath    string    `gorm:"column:icon_path;type:varchar(255)"`     // 图标URL
	CoinCost    int32     `gorm:"column:coin_cost;not null"`              // 消耗金币
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`       // 创建时间
}

func (ItemDO) TableName() string { return "items" }

// MessageDO 映射 messages 表
type MessageDO struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"` // 消息ID
	UserID      int64     `gorm:"column:user_id;not null"`            // 用户ID
	Sender      int32     `gorm:"column:sender;not null"`             // 0用户 1AI
	MessageType int32     `gorm:"column:message_type;not null"`       // 0聊天 1小纸条
	IsLocked    bool      `gorm:"column:is_locked;not null"`          // 锁定
	UnlockCoins int32     `gorm:"column:unlock_coins;not null"`       // 解锁金币
	Content     string    `gorm:"column:content;type:text;not null"`  // 内容
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`   // 创建时间
}

func (MessageDO) TableName() string { return "messages" }

// AvatarRepo 使用 GORM 的数据实现
type AvatarRepo struct{ data *Data }

func NewAvatarRepo(d *Data) *AvatarRepo { return &AvatarRepo{data: d} }

// ListPetModels 获取模型列表
func (r *AvatarRepo) ListPetModels(ctx context.Context) ([]*biz.PetModel, error) {
	var rows []PetModelDO
	tx := r.data.Gorm.WithContext(ctx).
		Order("type asc, is_default desc, sort_order asc, id asc").
		Find(&rows)
	if tx.Error != nil {
		return nil, tx.Error
	}
	out := make([]*biz.PetModel, 0, len(rows))
	for _, v := range rows {
		vv := v
		out = append(out, &biz.PetModel{ID: vv.ID, Name: vv.Name, Path: vv.Path, ModelType: vv.Type, IsDefault: vv.IsDefault, SortOrder: vv.SortOrder, CreatedAt: vv.CreatedAt})
	}
	return out, nil
}

// PetModelExists 判断模型是否存在
func (r *AvatarRepo) PetModelExists(ctx context.Context, modelID int64) (bool, error) {
	var count int64
	tx := r.data.Gorm.WithContext(ctx).Model(&PetModelDO{}).Where("id=?", modelID).Count(&count)
	return count > 0, tx.Error
}

// GetModelPath 获取模型路径
func (r *AvatarRepo) GetModelPath(ctx context.Context, modelID int64) (string, error) {
	var row PetModelDO
	if err := r.data.Gorm.WithContext(ctx).Select("path").First(&row, modelID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", biz.ErrAvatarNotFound
		}
		return "", err
	}
	return row.Path, nil
}

// SetUserModel 更新用户当前模型（同时刷新 model_url）
func (r *AvatarRepo) SetUserModel(ctx context.Context, userID, modelID int64) error {
	path, err := r.GetModelPath(ctx, modelID)
	if err != nil {
		return err
	}
	tx := r.data.Gorm.WithContext(ctx).
		Exec("UPDATE users SET model_id=?, model_url=?, updated_at=NOW() WHERE id=?", modelID, path, userID)
	return tx.Error
}

// ListItems 获取道具列表
func (r *AvatarRepo) ListItems(ctx context.Context) ([]*biz.Item, error) {
	var rows []ItemDO
	if err := r.data.Gorm.WithContext(ctx).Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*biz.Item, 0, len(rows))
	for _, v := range rows {
		vv := v
		out = append(out, &biz.Item{ID: vv.ID, Name: vv.Name, Description: vv.Description, IconPath: vv.IconPath, CoinCost: vv.CoinCost, CreatedAt: vv.CreatedAt})
	}
	return out, nil
}

// UseItem 扣金币并返回剩余金币（事务）
func (r *AvatarRepo) UseItem(ctx context.Context, userID, itemID int64) (int32, error) {
	var remaining int32
	err := r.data.Gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 读取道具
		var it ItemDO
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&it, itemID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrPropNotFound
			}
			return err
		}
		// 读取用户金币并加锁
		type userRow struct{ Coins int32 }
		var ur userRow
		if err := tx.Table("users").Clauses(clause.Locking{Strength: "UPDATE"}).Select("coins").Where("id=?", userID).Take(&ur).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrUserNotFound
			}
			return err
		}
		if ur.Coins < it.CoinCost {
			return biz.ErrInsufficientCoins
		}
		remaining = ur.Coins - it.CoinCost
		if err := tx.Exec("UPDATE users SET coins=?, updated_at=NOW() WHERE id=?", remaining, userID).Error; err != nil {
			return err
		}
		return nil
	})
	return remaining, err
}

// CreateChat 写入一条用户消息
func (r *AvatarRepo) CreateChat(ctx context.Context, userID int64, content string) (*biz.ChatMsg, error) {
	row := &MessageDO{UserID: userID, Sender: 0, MessageType: 0, IsLocked: false, UnlockCoins: 0, Content: content}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}
	return &biz.ChatMsg{
		ID:          row.ID,
		UserID:      row.UserID,
		Sender:      row.Sender,
		MessageType: row.MessageType,
		IsLocked:    row.IsLocked,
		UnlockCoins: row.UnlockCoins,
		Content:     row.Content,
		CreatedAt:   row.CreatedAt,
	}, nil
}

// CreateAIChat 调用 AI 并将回复落库为一条来自 AI 的消息
func (r *AvatarRepo) CreateAIChat(ctx context.Context, userID int64, content string) (*biz.ChatMsg, error) {
	// 读取用户/宠物关键信息作为 system prompt 的一部分（可选，不阻塞）
	var profile struct{ Nickname, PetName, Kind string }
	_ = r.data.Gorm.WithContext(ctx).Table("users").
		Select("nickname,pet_name,kind").Where("id=?", userID).Scan(&profile).Error

	system := "你是一个治愈系的宠物数字伙伴，以第一人称‘我’的口吻，温柔简短地回复。"
	if profile.PetName != "" {
		system += " 我的名字是" + profile.PetName + "。"
	}
	if profile.Kind != "" {
		system += " 我是一只" + profile.Kind + "。"
	}
	c := aiclient.Default()
	if c == nil {
		aiclient.SetClient(aiclient.NewClient(aiclient.Config{}))
		c = aiclient.Default()
	}
	reply, err := c.Chat(ctx, system, content)
	if err != nil {
		// 失败时也写一条兜底回复，保证对话连续
		reply = "我听到了你的心情，我会一直陪着你~"
	}
	row := &MessageDO{UserID: userID, Sender: 1, MessageType: 0, IsLocked: false, UnlockCoins: 0, Content: reply}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}
	return &biz.ChatMsg{
		ID:          row.ID,
		UserID:      row.UserID,
		Sender:      row.Sender,
		MessageType: row.MessageType,
		IsLocked:    row.IsLocked,
		UnlockCoins: row.UnlockCoins,
		Content:     row.Content,
		CreatedAt:   row.CreatedAt,
	}, nil
}

// CreateAIMessage 直接写入一条来自 AI 的消息
func (r *AvatarRepo) CreateAIMessage(ctx context.Context, userID int64, content string) (*biz.ChatMsg, error) {
	row := &MessageDO{UserID: userID, Sender: 1, MessageType: 0, IsLocked: false, UnlockCoins: 0, Content: content}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}
	return &biz.ChatMsg{ID: row.ID, UserID: row.UserID, Sender: row.Sender, MessageType: row.MessageType, IsLocked: row.IsLocked, UnlockCoins: row.UnlockCoins, Content: row.Content, CreatedAt: row.CreatedAt}, nil
}
