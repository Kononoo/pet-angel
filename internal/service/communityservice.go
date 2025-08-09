package service

import (
	"context"
	"net/http"
	"strings"

	pb "pet-angel/api/community/v1"
	"pet-angel/internal/biz"
	"pet-angel/internal/conf"
	jwtutil "pet-angel/internal/util/jwt"

	"github.com/go-kratos/kratos/v2/transport"
)

// CommunityService 提供社区相关接口的服务适配层
// - 公开接口：分类、列表、详情
// - 需鉴权接口：发帖、点赞/取消点赞、评论/点赞评论

type CommunityService struct {
	pb.UnimplementedCommunityServiceServer
	uc        *biz.CommunityUsecase
	jwtSecret string
}

func NewCommunityService(uc *biz.CommunityUsecase, cfg *conf.Auth) *CommunityService {
	secret := ""
	if cfg != nil {
		secret = cfg.JwtSecret
	}
	return &CommunityService{uc: uc, jwtSecret: secret}
}

// userID 解析登录用户ID
func (s *CommunityService) userID(ctx context.Context) (int64, error) {
	ts, ok := transport.FromServerContext(ctx)
	if !ok {
		return 0, http.ErrNoCookie
	}
	tok, err := jwtutil.FromAuthHeader(ts.RequestHeader().Get("Authorization"))
	if err != nil {
		return 0, err
	}
	claims, err := jwtutil.Parse(s.jwtSecret, tok)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetCategories 分类列表（公开）
func (s *CommunityService) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesReply, error) {
	list, err := s.uc.GetCategories(ctx)
	if err != nil {
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
	total, list, err := s.uc.GetPostList(ctx, viewerID, req.GetCategoryId(), req.GetPostType(), strings.ToLower(req.GetSort()), req.GetPage(), req.GetPageSize())
	if err != nil {
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
		return nil, err
	}
	return &pb.CreatePostReply{Id: id}, nil
}

// LikePost 点赞（鉴权）
func (s *CommunityService) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.uc.LikePost(ctx, uid, req.GetPostId()); err != nil {
		return nil, err
	}
	return &pb.LikePostReply{Success: true}, nil
}

// UnlikePost 取消点赞（鉴权）
func (s *CommunityService) UnlikePost(ctx context.Context, req *pb.UnlikePostRequest) (*pb.UnlikePostReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.uc.UnlikePost(ctx, uid, req.GetPostId()); err != nil {
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
	total, list, err := s.uc.GetCommentList(ctx, viewerID, req.GetPostId(), req.GetPage(), req.GetPageSize())
	if err != nil {
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
		return nil, err
	}
	id, err := s.uc.CreateComment(ctx, uid, req.GetPostId(), req.GetContent())
	if err != nil {
		return nil, err
	}
	return &pb.CreateCommentReply{Id: id}, nil
}

// LikeComment 点赞评论（鉴权）
func (s *CommunityService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.uc.LikeComment(ctx, uid, req.GetCommentId()); err != nil {
		return nil, err
	}
	return &pb.LikeCommentReply{Success: true}, nil
}

// UnlikeComment 取消点赞评论（鉴权）
func (s *CommunityService) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.UnlikeCommentReply, error) {
	uid, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.uc.UnlikeComment(ctx, uid, req.GetCommentId()); err != nil {
		return nil, err
	}
	return &pb.UnlikeCommentReply{Success: true}, nil
}
