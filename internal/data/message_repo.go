package data

import (
	"context"
	"errors"

	"pet-angel/internal/biz"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserUnlockRecordDO 映射 user_unlock_records 表
// 记录用户解锁小纸条的扣费流水

type UserUnlockRecordDO struct {
	ID        int64 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64 `gorm:"column:user_id;not null"`
	MessageID int64 `gorm:"column:message_id;not null"`
	Coins     int32 `gorm:"column:coins_spent;not null"`
}

func (UserUnlockRecordDO) TableName() string { return "user_unlock_records" }

// MessageRepoImpl 实现 biz.MessageRepo

type MessageRepoImpl struct{ data *Data }

func NewMessageRepo(d *Data) *MessageRepoImpl { return &MessageRepoImpl{data: d} }

// ListMessages 倒序分页
func (r *MessageRepoImpl) ListMessages(ctx context.Context, userID int64, onlyNotes bool, page, pageSize int32) (int32, []*biz.Message, error) {
	if r.data.Gorm == nil {
		return 0, nil, errors.New("gorm not initialized")
	}
	q := r.data.Gorm.WithContext(ctx).
		Model(&MessageDO{}).
		Where("user_id=?", userID)
	if onlyNotes {
		q = q.Where("message_type=?", 1)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []MessageDO
	offset := int((page - 1) * pageSize)
	if err := q.
		Order("created_at desc, id desc").
		Limit(int(pageSize)).
		Offset(offset).
		Find(&rows).Error; err != nil {
		return 0, nil, err
	}
	list := make([]*biz.Message, 0, len(rows))
	for _, v := range rows {
		vv := v
		list = append(list, &biz.Message{
			ID:          vv.ID,
			UserID:      vv.UserID,
			Sender:      vv.Sender,
			MessageType: vv.MessageType,
			IsLocked:    vv.IsLocked,
			UnlockCoins: vv.UnlockCoins,
			Content:     vv.Content,
			CreatedAt:   vv.CreatedAt,
		})
	}
	return int32(total), list, nil
}

// GetMessageByID 查询单条
func (r *MessageRepoImpl) GetMessageByID(ctx context.Context, userID, messageID int64) (*biz.Message, error) {
	var m MessageDO
	tx := r.data.Gorm.WithContext(ctx).
		Where("id=? AND user_id=?", messageID, userID).
		Take(&m)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrMessageNotFound
		}
		return nil, tx.Error
	}
	return &biz.Message{
		ID:          m.ID,
		UserID:      m.UserID,
		Sender:      m.Sender,
		MessageType: m.MessageType,
		IsLocked:    m.IsLocked,
		UnlockCoins: m.UnlockCoins,
		Content:     m.Content,
		CreatedAt:   m.CreatedAt,
	}, nil
}

// UnlockMessage 事务扣金币并解锁
func (r *MessageRepoImpl) UnlockMessage(ctx context.Context, userID, messageID int64) (int32, *biz.Message, error) {
	var remaining int32
	var out *biz.Message
	err := r.data.Gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 锁定消息记录
		var m MessageDO
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id=? AND user_id=? AND message_type=?", messageID, userID, 1).
			Take(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return biz.ErrMessageNotFound
			}
			return err
		}
		if !m.IsLocked {
			// 已解锁，直接返回当前信息与余额
			type urow struct{ Coins int32 }
			var ur urow
			if err := tx.
				Table("users").
				Select("coins").
				Where("id=?", userID).
				Take(&ur).Error; err != nil {
				return err
			}
			remaining = ur.Coins
			out = &biz.Message{ID: m.ID, UserID: m.UserID, Sender: m.Sender, MessageType: m.MessageType, IsLocked: m.IsLocked, UnlockCoins: m.UnlockCoins, Content: m.Content, CreatedAt: m.CreatedAt}
			return nil
		}
		// 扣金币
		type urow struct{ Coins int32 }
		var ur urow
		if err := tx.
			Table("users").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("coins").
			Where("id=?", userID).
			Take(&ur).Error; err != nil {
			return err
		}
		if ur.Coins < m.UnlockCoins {
			return biz.ErrInsufficientCoins
		}
		remaining = ur.Coins - m.UnlockCoins
		if err := tx.Exec("UPDATE users SET coins=?, updated_at=NOW() WHERE id=?", remaining, userID).Error; err != nil {
			return err
		}
		// 解锁消息
		if err := tx.
			Model(&MessageDO{}).
			Where("id=?", m.ID).
			Update("is_locked", false).Error; err != nil {
			return err
		}
		// 记录解锁流水
		rec := &UserUnlockRecordDO{UserID: userID, MessageID: m.ID, Coins: m.UnlockCoins}
		if err := tx.Create(rec).Error; err != nil {
			return err
		}
		// 返回最新消息
		out = &biz.Message{ID: m.ID, UserID: m.UserID, Sender: m.Sender, MessageType: m.MessageType, IsLocked: false, UnlockCoins: m.UnlockCoins, Content: m.Content, CreatedAt: m.CreatedAt}
		return nil
	})
	if err != nil {
		return 0, nil, err
	}
	return remaining, out, nil
}

// CreateLockedNote 生成一条锁定的小纸条（AI 个性化内容）
func (r *MessageRepoImpl) CreateLockedNote(ctx context.Context, userID int64, unlockCoins int32, content string) (int64, error) {
	row := &MessageDO{UserID: userID, Sender: 1, MessageType: 1, IsLocked: true, UnlockCoins: unlockCoins, Content: content}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return 0, err
	}
	return row.ID, nil
}
