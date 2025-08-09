package biz

import (
	"context"
	"time"

	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"golang.org/x/crypto/bcrypt"
)

// AuthRepo 抽象
// 注意：与 data 层实现解耦
type AuthRepo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, userID int64) (*User, error)
	Create(ctx context.Context, user *User) (int64, error)
	UpdateInfo(ctx context.Context, user *User) error
	UpdateCoins(ctx context.Context, userID int64, delta int32) error
}

// AuthUsecase 用例
type AuthUsecase struct {
	repo AuthRepo
	cfg  *conf.Auth
}

func NewAuthUsecase(repo AuthRepo, cfg *conf.Auth) *AuthUsecase {
	return &AuthUsecase{repo: repo, cfg: cfg}
}

// Login 用户名+密码登录
func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (user *User, token string, expiresIn int32, err error) {
	u, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		// 用户不存在则创建
		nu := &User{
			Username:  username,
			Password:  password,
			Nickname:  username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		id, cErr := uc.repo.Create(ctx, nu)
		if cErr != nil {
			return nil, "", 0, cErr
		}
		nu.Id = id
		u = nu
	}
	// 校验密码（支持 bcrypt 哈希）
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		// 兼容极端情况：历史明文存储
		if u.Password != password {
			return nil, "", 0, ErrInvalidPassword
		}
	}
	// 签发 JWT
	var ttl time.Duration = time.Hour * 72
	if uc.cfg != nil && uc.cfg.JwtTtl != nil {
		ttl = uc.cfg.JwtTtl.AsDuration()
	}
	secret := ""
	if uc.cfg != nil {
		secret = uc.cfg.JwtSecret
	}
	jwtStr, exp, jerr := jwtutil.Sign(secret, u.Id, ttl)
	if jerr != nil {
		return nil, "", 0, jerr
	}
	return u, jwtStr, int32(time.Until(exp).Seconds()), nil
}

func (uc *AuthUsecase) GetUserInfo(ctx context.Context, userID int64) (*User, error) {
	return uc.repo.GetByID(ctx, userID)
}

func (uc *AuthUsecase) UpdateUserInfo(ctx context.Context, user *User) error {
	return uc.repo.UpdateInfo(ctx, user)
}
