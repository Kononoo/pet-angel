package service

import (
	"context"
	"strings"

	pb "pet-angel/api/community/v1"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// CommunityService 提供社区相关接口的服务适配层
// - 公开接口：分类、列表、详情
// - 需鉴权接口：发帖、点赞/取消点赞、评论/点赞评论

type CommunityService struct {
	pb.UnimplementedCommunityServiceServer
	uc        *biz.CommunityUsecase
	jwtSecret string
	logger    *log.Helper
}

func NewCommunityService(uc *biz.CommunityUsecase, cfg *conf.Auth, l log.Logger) *CommunityService {
	secret := ""
	if cfg != nil {
		secret = cfg.JwtSecret
	}
	return &CommunityService{uc: uc, jwtSecret: secret, logger: log.NewHelper(l)}
}

// userID 解析登录用户ID（demo模式放开校验）
func (s *CommunityService) userID(ctx context.Context) (int64, error) {
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

// GetCategories 分类列表（公开）
func (s *CommunityService) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesReply, error) {
	list, err := s.uc.GetCategories(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get categories failed: %v", err)
		return nil, err
	}
	out := &pb.GetCategoriesReply{}
	for _, c := range list {
		out.Categories = append(out.Categories, &pb.Category{Id: c.ID, Name: c.Name})
	}
	return out, nil
}

// GetPostList 帖子列表（公开）
func (s *CommunityService) GetPostList(ctx context.Context, req *pb.GetPostListRequest) (*pb.GetPostListReply, error) {
	viewerID := int64(0)
	if _, ok := transport.FromServerContext(ctx); ok {
		id, err := s.userID(ctx)
		if err == nil {
			viewerID = id
		}
	}

	// 设置默认分页参数
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total, list, err := s.uc.GetPostList(ctx, viewerID, req.GetCategoryId(), req.GetPostType(), strings.ToLower(req.GetSort()), page, pageSize)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get posts failed: %v", err)
		return nil, err
	}
	out := &pb.GetPostListReply{Total: total}
	for _, p := range list {
		out.List = append(out.List, &pb.Post{
			Id:           p.ID,
			User:         &pb.UserBrief{Id: p.User.Id, Nickname: p.User.Nickname, Avatar: p.User.Avatar},
			CategoryId:   p.CategoryID,
			Title:        p.Title,
			Content:      p.Content,
			PostType:     p.PostType,
			ImageUrls:    p.ImageUrls,
			VideoUrl:     p.VideoUrl,
			CoverUrl:     p.CoverUrl,
			Locate:       p.Locate,
			Tags:         p.Tags,
			LikedCount:   p.LikedCount,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
			IsLiked:      p.IsLiked,
			IsPrivate:    p.IsPrivate,
		})
	}
	return out, nil
}

// GetPostDetail 帖子详情（公开）
func (s *CommunityService) GetPostDetail(ctx context.Context, req *pb.GetPostDetailRequest) (*pb.GetPostDetailReply, error) {
	viewerID := int64(0)
	if _, ok := transport.FromServerContext(ctx); ok {
		id, err := s.userID(ctx)
		if err == nil {
			viewerID = id
		}
	}
	p, err := s.uc.GetPostDetail(ctx, viewerID, req.GetPostId())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get post detail failed: %v", err)
		return nil, err
	}
	return &pb.GetPostDetailReply{Post: &pb.Post{
		Id:           p.ID,
		User:         &pb.UserBrief{Id: p.User.Id, Nickname: p.User.Nickname, Avatar: p.User.Avatar},
		CategoryId:   p.CategoryID,
		Title:        p.Title,
		Content:      p.Content,
		PostType:     p.PostType,
		ImageUrls:    p.ImageUrls,
		VideoUrl:     p.VideoUrl,
		CoverUrl:     p.CoverUrl,
		Locate:       p.Locate,
		Tags:         p.Tags,
		LikedCount:   p.LikedCount,
		CommentCount: p.CommentCount,
		CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		IsLiked:      p.IsLiked,
		IsPrivate:    p.IsPrivate,
	}}, nil
}

// CreatePost 发帖（鉴权）
func (s *CommunityService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("create post: auth failed: %v", err)
		return nil, err
	}
	id, err := s.uc.CreatePost(ctx, uid, &biz.CommunityPost{
		CategoryID: req.GetCategoryId(),
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
		PostType:   req.GetPostType(),
		ImageUrls:  req.GetImageUrls(),
		VideoUrl:   req.GetVideoUrl(),
		CoverUrl:   req.GetCoverUrl(),
		Locate:     req.GetLocate(),
		Tags:       req.GetTags(),
		IsPrivate:  req.GetIsPrivate(),
	})
	if err != nil {
		s.logger.WithContext(ctx).Errorf("create post: usecase error: %v", err)
		return nil, err
	}
	return &pb.CreatePostReply{Id: id}, nil
}

// LikePost 点赞（鉴权）
func (s *CommunityService) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("like post: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.LikePost(ctx, uid, req.GetPostId()); err != nil {
		s.logger.WithContext(ctx).Errorf("like post: usecase error: %v", err)
		return nil, err
	}
	return &pb.LikePostReply{Success: true}, nil
}

// UnlikePost 取消点赞（鉴权）
func (s *CommunityService) UnlikePost(ctx context.Context, req *pb.UnlikePostRequest) (*pb.UnlikePostReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("unlike post: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.UnlikePost(ctx, uid, req.GetPostId()); err != nil {
		s.logger.WithContext(ctx).Errorf("unlike post: usecase error: %v", err)
		return nil, err
	}
	return &pb.UnlikePostReply{Success: true}, nil
}

// GetCommentList 评论列表（公开）
func (s *CommunityService) GetCommentList(ctx context.Context, req *pb.GetCommentListRequest) (*pb.GetCommentListReply, error) {
	viewerID := int64(0)
	if _, ok := transport.FromServerContext(ctx); ok {
		id, err := s.userID(ctx)
		if err == nil {
			viewerID = id
		}
	}

	// 设置默认分页参数
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	total, list, err := s.uc.GetCommentList(ctx, viewerID, req.GetPostId(), page, pageSize)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("get comment list failed: %v", err)
		return nil, err
	}
	out := &pb.GetCommentListReply{Total: total}
	for _, c := range list {
		out.List = append(out.List, &pb.Comment{
			Id:         c.ID,
			User:       &pb.UserBrief{Id: c.User.Id, Nickname: c.User.Nickname, Avatar: c.User.Avatar},
			Content:    c.Content,
			LikedCount: c.LikedCount,
			CreatedAt:  c.CreatedAt.Format("2006-01-02 15:04:05"),
			IsLiked:    c.IsLiked,
		})
	}
	return out, nil
}

// CreateComment 发表评论（鉴权）
func (s *CommunityService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("create comment: auth failed: %v", err)
		return nil, err
	}
	id, err := s.uc.CreateComment(ctx, uid, req.GetPostId(), req.GetContent())
	if err != nil {
		s.logger.WithContext(ctx).Errorf("create comment: usecase error: %v", err)
		return nil, err
	}
	return &pb.CreateCommentReply{Id: id}, nil
}

// LikeComment 点赞评论（鉴权）
func (s *CommunityService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("like comment: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.LikeComment(ctx, uid, req.GetCommentId()); err != nil {
		s.logger.WithContext(ctx).Errorf("like comment: usecase error: %v", err)
		return nil, err
	}
	return &pb.LikeCommentReply{Success: true}, nil
}

// UnlikeComment 取消点赞评论（鉴权）
func (s *CommunityService) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.UnlikeCommentReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		s.logger.WithContext(ctx).Errorf("unlike comment: auth failed: %v", err)
		return nil, err
	}
	if err := s.uc.UnlikeComment(ctx, uid, req.GetCommentId()); err != nil {
		s.logger.WithContext(ctx).Errorf("unlike comment: usecase error: %v", err)
		return nil, err
	}
	return &pb.UnlikeCommentReply{Success: true}, nil
}
