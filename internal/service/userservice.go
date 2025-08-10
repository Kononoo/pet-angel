package service

import (
	"context"

	pb "pet-angel/api/user/v1"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// UserService 提供用户关系与主页相关接口的服务适配层
// 责任：
// - 读取 JWT 解析当前用户 ID
// - 参数校验与类型适配
// - 调用 usecase 并将领域对象映射为 proto 返回体

type UserService struct {
	pb.UnimplementedUserServiceServer
	uc        *biz.UserUsecase
	jwtSecret string
	logger    *log.Helper
}

// NewUserService 创建 UserService
func NewUserService(uc *biz.UserUsecase, cfg *conf.Auth, l log.Logger) *UserService {
	secret := ""
	if cfg != nil {
		secret = cfg.JwtSecret
	}
	return &UserService{uc: uc, jwtSecret: secret, logger: log.NewHelper(l)}
}

// userID 从请求上下文解析当前登录用户 ID（demo模式放开校验）
func (s *UserService) userID(ctx context.Context) (int64, error) {
	// demo模式：如果没有token或token无效，返回默认用户ID
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return 1, nil // 默认用户ID
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return 1, nil // 默认用户ID
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return 1, nil // 默认用户ID
	}
	return claims.UserID, nil
}

// FollowUser 关注用户
func (s *UserService) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("follow user: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.Follow(ctx, uid, req.GetTargetUserId()); err != nil {
		s.logger.WithContext(ctx).Errorf("follow user: usecase error: %v", err)
		return nil, err
	}
	return &pb.FollowUserReply{Success: true}, nil
}

// UnfollowUser 取消关注
func (s *UserService) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.UnfollowUserReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("unfollow user: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.Unfollow(ctx, uid, req.GetTargetUserId()); err != nil {
		s.logger.WithContext(ctx).Errorf("unfollow user: usecase error: %v", err)
		return nil, err
	}
	return &pb.UnfollowUserReply{Success: true}, nil
}

// GetUserProfile 获取用户主页信息
func (s *UserService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileReply, error) {
	viewer, _ := s.userID(ctx)
	u, posts, isFollowed, err := s.uc.GetProfile(ctx, viewer, req.GetUserId())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get profile failed: %v", err)
		return nil, err
	}
	out := &pb.GetUserProfileReply{
		Avatar:      u.Avatar,
		Nickname:    u.Nickname,
		PetName:     u.PetName,
		PetSex:      u.PetSex,
		Kind:        u.Kind,
		Weight:      u.Weight,
		Hobby:       u.Hobby,
		Description: u.Description,
		IsFollowed:  isFollowed,
	}
	for _, p := range posts {
		out.Posts = append(out.Posts, &pb.PostBrief{
			Id:         p.Id,
			Title:      p.Title,
			PostType:   p.PostType,
			CoverUrl:   p.CoverUrl,
			LikedCount: p.LikedCount,
			CreatedAt:  p.CreatedAt,
		})
	}
	return out, nil
}

// GetFollowList 获取关注列表
func (s *UserService) GetFollowList(ctx context.Context, req *pb.GetFollowListRequest) (*pb.GetFollowListReply, error) {
	list, total, err := s.uc.GetFollowList(ctx, req.GetUserId(), req.GetPage(), req.GetPageSize())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get follow list failed: %v", err)
		return nil, err
	}
	out := &pb.GetFollowListReply{Total: total}
	for _, it := range list {
		out.List = append(out.List, &pb.UserBrief{Id: it.Id, Nickname: it.Nickname, Avatar: it.Avatar})
	}
	return out, nil
}

// GetLikeList 获取点赞过的帖子列表
func (s *UserService) GetLikeList(ctx context.Context, req *pb.GetLikeListRequest) (*pb.GetLikeListReply, error) {
	list, total, err := s.uc.GetLikeList(ctx, req.GetUserId(), req.GetPage(), req.GetPageSize())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get like list failed: %v", err)
		return nil, err
	}
	out := &pb.GetLikeListReply{Total: total}
	for _, p := range list {
		out.List = append(out.List, &pb.PostBrief{Id: p.Id, Title: p.Title, PostType: p.PostType, CoverUrl: p.CoverUrl, LikedCount: p.LikedCount, CreatedAt: p.CreatedAt})
	}
	return out, nil
}
