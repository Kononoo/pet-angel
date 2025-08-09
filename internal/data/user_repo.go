package data

import (
	"context"

	"pet-angel/internal/biz"
)

// GORM 模型：用户表（只含 user 模块需要字段）
// 对应表：users

type UserModel struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`           // 主键
	Username    string `gorm:"column:username;uniqueIndex;size:64;not null"` // 登录名
	Password    string `gorm:"column:password;size:255;not null"`            // 密码哈希
	Nickname    string `gorm:"column:nickname;size:50"`                      // 昵称
	Avatar      string `gorm:"column:avatar;size:255"`                       // 头像
	ModelID     int64  `gorm:"column:model_id"`                              // 当前模型
	ModelURL    string `gorm:"column:model_url;size:255"`                    // 当前模型URL
	PetName     string `gorm:"column:pet_name;size:50"`                      // 宠物名
	PetAvatar   string `gorm:"column:pet_avatar;size:255"`                   // 宠物头像
	PetSex      int32  `gorm:"column:pet_sex"`                               // 宠物性别
	Kind        string `gorm:"column:kind;size:50"`                          // 宠物种类
	Weight      int32  `gorm:"column:weight"`                                // 体重
	Hobby       string `gorm:"column:hobby;size:255"`                        // 爱好
	Description string `gorm:"column:description;type:text"`                 // 简介
	Coins       int32  `gorm:"column:coins"`                                 // 金币
}

func (UserModel) TableName() string { return "users" }

// GORM 模型：关注关系
// 对应表：user_follows

type FollowModel struct {
	FollowerID int64 `gorm:"column:follower_id;primaryKey"` // 关注者
	FolloweeID int64 `gorm:"column:followee_id;primaryKey"` // 被关注者
}

func (FollowModel) TableName() string { return "user_follows" }

// GORM 模型：帖子简要（user 模块内部使用）
// 对应表：posts

type UserPostModel struct {
	ID         int64  `gorm:"column:id;primaryKey"`
	UserID     int64  `gorm:"column:user_id"`
	Title      string `gorm:"column:title"`
	Type       int32  `gorm:"column:type"`
	CoverUrl   string `gorm:"column:cover_url"`
	LikedCount int32  `gorm:"column:liked_count"`
	CreatedAt  string `gorm:"column:created_at"`
}

func (UserPostModel) TableName() string { return "posts" }

// UserRepoImpl 实现（GORM 优先，Gorm 为空时返回空集）

type UserRepoImpl struct{ data *Data }

func NewUserRepo(d *Data) *UserRepoImpl {
	return &UserRepoImpl{data: d}
}

func (r *UserRepoImpl) Follow(ctx context.Context, followerID, followeeID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	return r.data.Gorm.WithContext(ctx).Clauses(
	// upsert 由复合主键保证幂等
	).Create(&FollowModel{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}).Error
}

func (r *UserRepoImpl) Unfollow(ctx context.Context, followerID, followeeID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	return r.data.Gorm.WithContext(ctx).
		Where("follower_id=? AND followee_id=?", followerID, followeeID).
		Delete(&FollowModel{}).Error
}

func (r *UserRepoImpl) IsFollow(ctx context.Context, followerID, followeeID int64) (bool, error) {
	if r.data.Gorm == nil {
		return false, nil
	}
	var cnt int64
	if err := r.data.Gorm.WithContext(ctx).
		Model(&FollowModel{}).
		Where("follower_id=? AND followee_id=?", followerID, followeeID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *UserRepoImpl) GetUserByID(ctx context.Context, userID int64) (*biz.User, error) {
	if r.data.Gorm == nil {
		return nil, nil
	}
	var m UserModel
	if err := r.data.Gorm.WithContext(ctx).First(&m, userID).Error; err != nil {
		return nil, err
	}
	return &biz.User{
		Id:          m.ID,
		Username:    m.Username,
		Password:    m.Password,
		Nickname:    m.Nickname,
		Avatar:      m.Avatar,
		ModelID:     m.ModelID,
		ModelURL:    m.ModelURL,
		PetName:     m.PetName,
		PetAvatar:   m.PetAvatar,
		PetSex:      m.PetSex,
		Kind:        m.Kind,
		Weight:      m.Weight,
		Hobby:       m.Hobby,
		Description: m.Description,
		Coins:       m.Coins,
	}, nil
}

func (r *UserRepoImpl) GetUserPostsBrief(ctx context.Context, userID int64, limit int32) ([]*biz.PostBrief, error) {
	if r.data.Gorm == nil {
		return []*biz.PostBrief{}, nil
	}
	var rows []UserPostModel
	if err := r.data.Gorm.WithContext(ctx).
		Select("id,title,type,cover_url,liked_count,DATE_FORMAT(created_at,'%Y-%m-%d %H:%i:%s') as created_at").
		Where("user_id=?", userID).
		Order("id DESC").
		Limit(int(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*biz.PostBrief, 0, len(rows))
	for _, p := range rows {
		pp := p
		out = append(out, &biz.PostBrief{
			Id:         pp.ID,
			Title:      pp.Title,
			PostType:   pp.Type,
			CoverUrl:   pp.CoverUrl,
			LikedCount: pp.LikedCount,
			CreatedAt:  pp.CreatedAt,
		})
	}
	return out, nil
}

func (r *UserRepoImpl) GetFollowList(ctx context.Context, userID int64, page, pageSize int32) ([]*biz.UserBrief, int32, error) {
	if r.data.Gorm == nil {
		return []*biz.UserBrief{}, 0, nil
	}
	var total int64
	if err := r.data.Gorm.WithContext(ctx).Model(&FollowModel{}).Where("follower_id=?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var rows []struct {
		ID       int64
		Nickname string
		Avatar   string
	}
	if err := r.data.Gorm.WithContext(ctx).
		Table("user_follows f").
		Select("u.id,u.nickname,u.avatar").
		Joins("JOIN users u ON f.followee_id=u.id").
		Where("f.follower_id=?", userID).
		Order("f.created_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*biz.UserBrief, 0, len(rows))
	for _, r := range rows {
		rr := r
		out = append(out, &biz.UserBrief{
			Id:       rr.ID,
			Nickname: rr.Nickname,
			Avatar:   rr.Avatar,
		})
	}
	return out, int32(total), nil
}

func (r *UserRepoImpl) GetLikeList(ctx context.Context, userID int64, page, pageSize int32) ([]*biz.PostBrief, int32, error) {
	if r.data.Gorm == nil {
		return []*biz.PostBrief{}, 0, nil
	}
	var total int64
	if err := r.data.Gorm.WithContext(ctx).Table("likes").Where("user_id=? AND target_type=0", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	var rows []UserPostModel
	if err := r.data.Gorm.WithContext(ctx).
		Table("likes l").
		Select("p.id,p.title,p.type,p.cover_url,p.liked_count,DATE_FORMAT(p.created_at,'%Y-%m-%d %H:%i:%s') as created_at").
		Joins("JOIN posts p ON l.target_id=p.id").
		Where("l.user_id=? AND l.target_type=0", userID).
		Order("l.created_at DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*biz.PostBrief, 0, len(rows))
	for _, p := range rows {
		pp := p
		out = append(out, &biz.PostBrief{
			Id:         pp.ID,
			Title:      pp.Title,
			PostType:   pp.Type,
			CoverUrl:   pp.CoverUrl,
			LikedCount: pp.LikedCount,
			CreatedAt:  pp.CreatedAt,
		})
	}
	return out, int32(total), nil
}
