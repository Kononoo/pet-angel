package service

import (
	"context"
	"net/http"

	authv1 "pet-angel/api/auth/v1"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/transport"
)

// AuthService 认证服务
// 提供登录、重新登录校验、获取与更新用户信息的接口

type AuthService struct {
	authv1.UnimplementedAuthServiceServer
	uc        *biz.AuthUsecase
	jwtSecret string
}

// NewAuthService 构造函数，注入用例与 JWT 配置
func NewAuthService(uc *biz.AuthUsecase, authCfg *conf.Auth) *AuthService {
	secret := ""
	if authCfg != nil {
		secret = authCfg.JwtSecret
	}
	return &AuthService{uc: uc, jwtSecret: secret}
}

// Login 用户名+密码登录，返回 JWT
func (s *AuthService) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginReply, error) {
	u, token, exp, err := s.uc.Login(ctx, in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return &authv1.LoginReply{UserId: u.Id, Token: token, ExpiresIn: exp}, nil
}

// Relogin 校验当前请求头中的 JWT 是否有效
func (s *AuthService) Relogin(ctx context.Context, in *authv1.ReloginRequest) (*authv1.ReloginReply, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return &authv1.ReloginReply{Expire: true}, nil
	}
	authHeader := ts.RequestHeader().Get("Authorization")
	tok, err := jwtutil.FromAuthHeader(authHeader)
	if err != nil {
		return &authv1.ReloginReply{Expire: true}, nil
	}
	_, perr := jwtutil.Parse(s.jwtSecret, tok)
	return &authv1.ReloginReply{Expire: perr != nil}, nil
}

// GetUserInfo 解析 JWT 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, in *authv1.GetUserInfoRequest) (*authv1.GetUserInfoReply, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return nil, err
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return nil, err
	}
	u, err := s.uc.GetUserInfo(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	return &authv1.GetUserInfoReply{
		UserId:      u.Id,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		ModelId:     u.ModelID,
		PetName:     u.PetName,
		PetAvatar:   u.PetAvatar,
		PetSex:      u.PetSex,
		Kind:        u.Kind,
		Weight:      int32(u.Weight),
		Hobby:       u.Hobby,
		Description: u.Description,
		Coins:       u.Coins,
		CreatedAt:   u.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateUserInfo 更新当前用户资料
func (s *AuthService) UpdateUserInfo(ctx context.Context, in *authv1.UpdateUserInfoRequest) (*authv1.UpdateUserInfoReply, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, http.ErrNoCookie
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return nil, err
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return nil, err
	}
	user := &biz.User{
		Id:          claims.UserID,
		Nickname:    in.GetNickname(),
		Avatar:      in.GetAvatar(),
		ModelID:     in.GetModelId(),
		PetName:     in.GetPetName(),
		PetAvatar:   in.GetPetAvatar(),
		PetSex:      in.GetPetSex(),
		Kind:        in.GetKind(),
		Weight:      in.GetWeight(),
		Hobby:       in.GetHobby(),
		Description: in.GetDescription(),
	}
	if err := s.uc.UpdateUserInfo(ctx, user); err != nil {
		return nil, err
	}
	return &authv1.UpdateUserInfoReply{Success: true}, nil
}
