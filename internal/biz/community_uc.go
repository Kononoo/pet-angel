package biz

import (
	"context"
	"time"
)

// Category 表示帖子分类
// 对应表：categories
// 仅包含社区模块所需字段

type Category struct {
	ID   int64  // 分类ID
	Name string // 分类名称
}

// CommunityPost 帖子聚合视图
// 混合了 posts 表、用户信息与当前查看者态（是否点赞）

type CommunityPost struct {
	ID           int64     // 帖子ID
	User         UserBrief // 作者信息
	CategoryID   int64     // 分类ID
	Title        string    // 标题
	Content      string    // 文本内容
	PostType     int32     // 0图文 1视频
	ImageUrls    []string  // 图片URL列表
	VideoUrl     string    // 视频URL
	CoverUrl     string    // 封面URL
	Locate       string    // 位置
	Tags         []string  // 标签
	LikedCount   int32     // 点赞数
	CommentCount int32     // 评论数
	CreatedAt    time.Time // 创建时间
	IsLiked      bool      // 当前查看者是否已点赞
	IsPrivate    bool      // 是否私密
}

// CommunityComment 评论聚合视图

type CommunityComment struct {
	ID         int64     // 评论ID
	User       UserBrief // 评论者
	Content    string    // 内容
	LikedCount int32     // 点赞数
	CreatedAt  time.Time // 创建时间
	IsLiked    bool      // 当前查看者是否点赞
}

// CommunityRepo 社区仓储接口（由 data 层实现，优先使用 GORM）

type CommunityRepo interface {
	ListCategories(ctx context.Context) ([]*Category, error)

	ListPosts(ctx context.Context, viewerID int64, categoryID int64, postType int32, sort string, page, pageSize int32) (int32, []*CommunityPost, error)
	GetPostDetail(ctx context.Context, viewerID, postID int64) (*CommunityPost, error)
	CreatePost(ctx context.Context, userID int64, p *CommunityPost) (int64, error)
	LikePost(ctx context.Context, userID, postID int64) error
	UnlikePost(ctx context.Context, userID, postID int64) error

	ListComments(ctx context.Context, viewerID, postID int64, page, pageSize int32) (int32, []*CommunityComment, error)
	CreateComment(ctx context.Context, userID, postID int64, content string) (int64, error)
	LikeComment(ctx context.Context, userID, commentID int64) error
	UnlikeComment(ctx context.Context, userID, commentID int64) error
}

// CommunityUsecase 社区用例

type CommunityUsecase struct {
	repo CommunityRepo
}

func NewCommunityUsecase(repo CommunityRepo) *CommunityUsecase { return &CommunityUsecase{repo: repo} }

func (uc *CommunityUsecase) GetCategories(ctx context.Context) ([]*Category, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *CommunityUsecase) GetPostList(ctx context.Context, viewerID int64, categoryID int64, postType int32, sort string, page, pageSize int32) (int32, []*CommunityPost, error) {
	return uc.repo.ListPosts(ctx, viewerID, categoryID, postType, sort, page, pageSize)
}

func (uc *CommunityUsecase) GetPostDetail(ctx context.Context, viewerID, postID int64) (*CommunityPost, error) {
	return uc.repo.GetPostDetail(ctx, viewerID, postID)
}

func (uc *CommunityUsecase) CreatePost(ctx context.Context, userID int64, p *CommunityPost) (int64, error) {
	return uc.repo.CreatePost(ctx, userID, p)
}

func (uc *CommunityUsecase) LikePost(ctx context.Context, userID, postID int64) error {
	return uc.repo.LikePost(ctx, userID, postID)
}

func (uc *CommunityUsecase) UnlikePost(ctx context.Context, userID, postID int64) error {
	return uc.repo.UnlikePost(ctx, userID, postID)
}

func (uc *CommunityUsecase) GetCommentList(ctx context.Context, viewerID, postID int64, page, pageSize int32) (int32, []*CommunityComment, error) {
	return uc.repo.ListComments(ctx, viewerID, postID, page, pageSize)
}

func (uc *CommunityUsecase) CreateComment(ctx context.Context, userID, postID int64, content string) (int64, error) {
	return uc.repo.CreateComment(ctx, userID, postID, content)
}

func (uc *CommunityUsecase) LikeComment(ctx context.Context, userID, commentID int64) error {
	return uc.repo.LikeComment(ctx, userID, commentID)
}

func (uc *CommunityUsecase) UnlikeComment(ctx context.Context, userID, commentID int64) error {
	return uc.repo.UnlikeComment(ctx, userID, commentID)
}
