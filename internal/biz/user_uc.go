package biz

import (
	"context"
	"time"
)

// User 表示 users 表对应的业务实体
// 注意：Password 应存储 bcrypt 哈希；Weight 单位为 kg 的整数

type User struct {
	Id          int64     // users.id 主键
	Username    string    // 登录名（唯一）
	Password    string    // 密码哈希
	Nickname    string    // 昵称
	Avatar      string    // 头像URL
	ModelID     int64     // 当前模型ID
	ModelURL    string    // 当前模型URL（与 pet_models.path 一致）
	PetName     string    // 宠物名称
	PetAvatar   string    // 宠物头像
	PetSex      int32     // 宠物性别 0未知/1男/2女
	Kind        string    // 宠物种类
	Weight      int32     // 体重（kg）
	Hobby       string    // 爱好
	Description string    // 简介
	Coins       int32     // 金币余额
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}

// UserBrief 用户简要信息

type UserBrief struct {
	Id       int64  // 用户ID
	Nickname string // 昵称
	Avatar   string // 头像URL
}

// PostBrief 帖子简要（用于主页与点赞列表）

type PostBrief struct {
	Id         int64  // 帖子ID
	Title      string // 标题
	PostType   int32  // 0图文 1视频
	CoverUrl   string // 封面URL
	LikedCount int32  // 点赞数
	CreatedAt  string // 创建时间
}

// UserRepo 用户关系与主页相关仓储

type UserRepo interface {
	Follow(ctx context.Context, followerID, followeeID int64) error
	Unfollow(ctx context.Context, followerID, followeeID int64) error
	IsFollow(ctx context.Context, followerID, followeeID int64) (bool, error)
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	GetUserPostsBrief(ctx context.Context, userID int64, limit int32) ([]*PostBrief, error)
	GetFollowList(ctx context.Context, userID int64, page, pageSize int32) ([]*UserBrief, int32, error)
	GetLikeList(ctx context.Context, userID int64, page, pageSize int32) ([]*PostBrief, int32, error)

	GetModelPath(ctx context.Context, modelID int64) (string, error)
}

// UserUsecase 用户关系与主页用例

type UserUsecase struct{ repo UserRepo }

func NewUserUsecase(repo UserRepo) *UserUsecase { return &UserUsecase{repo: repo} }

func (u *UserUsecase) Follow(ctx context.Context, followerID, followeeID int64) error {
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}
	return u.repo.Follow(ctx, followerID, followeeID)
}

func (u *UserUsecase) Unfollow(ctx context.Context, followerID, followeeID int64) error {
	return u.repo.Unfollow(ctx, followerID, followeeID)
}

func (u *UserUsecase) GetProfile(ctx context.Context, viewerID, userID int64) (*User, []*PostBrief, bool, error) {
	usr, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, false, err
	}
	// 以 pet_models.path 为准，刷新 ModelURL
	if usr != nil && usr.ModelID > 0 {
		if path, perr := u.repo.GetModelPath(ctx, usr.ModelID); perr == nil && path != "" {
			usr.ModelURL = path
		}
	}
	posts, err := u.repo.GetUserPostsBrief(ctx, userID, 10)
	if err != nil {
		return usr, nil, false, err
	}
	isFollowed := false
	if viewerID > 0 {
		isFollowed, _ = u.repo.IsFollow(ctx, viewerID, userID)
	}
	return usr, posts, isFollowed, nil
}

func (u *UserUsecase) GetFollowList(ctx context.Context, userID int64, page, pageSize int32) ([]*UserBrief, int32, error) {
	return u.repo.GetFollowList(ctx, userID, page, pageSize)
}

func (u *UserUsecase) GetLikeList(ctx context.Context, userID int64, page, pageSize int32) ([]*PostBrief, int32, error) {
	return u.repo.GetLikeList(ctx, userID, page, pageSize)
}
