package biz

import (
	"context"
	"time"
)

// Post 帖子实体
type Post struct {
	PostID       int64     `json:"post_id"`
	UserID       int64     `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	Videos       []string  `json:"videos"`
	TagIDs       []int64   `json:"tag_ids"`
	LikeCount    int32     `json:"like_count"`
	CommentCount int32     `json:"comment_count"`
	ViewCount    int32     `json:"view_count"`
	IsLiked      bool      `json:"is_liked"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Author       *User     `json:"author"`
}

// Comment 评论实体
type Comment struct {
	CommentID int64     `json:"comment_id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	LikeCount int32     `json:"like_count"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author"`
}

// Tag 标签实体
type Tag struct {
	TagID     int64  `json:"tag_id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	PostCount int32  `json:"post_count"`
	IsDefault bool   `json:"is_default"`
}

// PostRepo 帖子仓储接口
type PostRepo interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPostByID(ctx context.Context, postID int64) (*Post, error)
	GetPostList(ctx context.Context, tagIDs []int64, sortBy string, page, pageSize int32) ([]*Post, int32, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, postID int64) error
	GetUserPosts(ctx context.Context, userID int64, page, pageSize int32) ([]*Post, int32, error)
	GetLikedPosts(ctx context.Context, userID int64, page, pageSize int32) ([]*Post, int32, error)
	LikePost(ctx context.Context, userID, postID int64) error
	UnlikePost(ctx context.Context, userID, postID int64) error
	IsLiked(ctx context.Context, userID, postID int64) (bool, error)
	IncrementViewCount(ctx context.Context, postID int64) error
}

// CommentRepo 评论仓储接口
type CommentRepo interface {
	CreateComment(ctx context.Context, comment *Comment) error
	GetCommentList(ctx context.Context, postID int64, page, pageSize int32) ([]*Comment, int32, error)
	DeleteComment(ctx context.Context, commentID int64) error
	LikeComment(ctx context.Context, userID, commentID int64) error
	UnlikeComment(ctx context.Context, userID, commentID int64) error
	IsCommentLiked(ctx context.Context, userID, commentID int64) (bool, error)
}

// TagRepo 标签仓储接口
type TagRepo interface {
	GetTagList(ctx context.Context) ([]*Tag, error)
	GetTagByID(ctx context.Context, tagID int64) (*Tag, error)
	UpdateTagPostCount(ctx context.Context, tagID int64, delta int32) error
}

// CommunityUsecase 社区用例
type CommunityUsecase struct {
	postRepo    PostRepo
	commentRepo CommentRepo
	tagRepo     TagRepo
	userRepo    UserRepo
}

// NewCommunityUsecase 创建社区用例
func NewCommunityUsecase(postRepo PostRepo, commentRepo CommentRepo, tagRepo TagRepo, userRepo UserRepo) *CommunityUsecase {
	return &CommunityUsecase{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
	}
}

// CreatePost 创建帖子
func (uc *CommunityUsecase) CreatePost(ctx context.Context, userID int64, title, content string, images, videos []string, tagIDs []int64, isDraft bool) (*Post, error) {
	status := "normal"
	if isDraft {
		status = "draft"
	}

	post := &Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Images:    images,
		Videos:    videos,
		TagIDs:    tagIDs,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.postRepo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	// 更新标签帖子数量
	for _, tagID := range tagIDs {
		uc.tagRepo.UpdateTagPostCount(ctx, tagID, 1)
	}

	return post, nil
}

// GetPostList 获取帖子列表
func (uc *CommunityUsecase) GetPostList(ctx context.Context, tagIDs []int64, sortBy string, page, pageSize int32, currentUserID int64) ([]*Post, int32, error) {
	posts, total, err := uc.postRepo.GetPostList(ctx, tagIDs, sortBy, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 填充作者信息和点赞状态
	for _, post := range posts {
		// 获取作者信息
		author, err := uc.userRepo.GetUserByID(ctx, post.UserID)
		if err == nil {
			post.Author = author
		}

		// 获取点赞状态
		if currentUserID > 0 {
			isLiked, _ := uc.postRepo.IsLiked(ctx, currentUserID, post.PostID)
			post.IsLiked = isLiked
		}
	}

	return posts, total, nil
}

// GetPostDetail 获取帖子详情
func (uc *CommunityUsecase) GetPostDetail(ctx context.Context, postID, currentUserID int64) (*Post, error) {
	post, err := uc.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.Status == "deleted" {
		return nil, ErrPostDeleted
	}

	// 获取作者信息
	author, err := uc.userRepo.GetUserByID(ctx, post.UserID)
	if err == nil {
		post.Author = author
	}

	// 获取点赞状态
	if currentUserID > 0 {
		isLiked, _ := uc.postRepo.IsLiked(ctx, currentUserID, post.PostID)
		post.IsLiked = isLiked
	}

	// 增加浏览数
	uc.postRepo.IncrementViewCount(ctx, postID)

	return post, nil
}

// LikePost 点赞帖子
func (uc *CommunityUsecase) LikePost(ctx context.Context, userID, postID int64) error {
	return uc.postRepo.LikePost(ctx, userID, postID)
}

// UnlikePost 取消点赞
func (uc *CommunityUsecase) UnlikePost(ctx context.Context, userID, postID int64) error {
	return uc.postRepo.UnlikePost(ctx, userID, postID)
}

// CreateComment 创建评论
func (uc *CommunityUsecase) CreateComment(ctx context.Context, userID, postID int64, content string) (*Comment, error) {
	comment := &Comment{
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	err := uc.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentList 获取评论列表
func (uc *CommunityUsecase) GetCommentList(ctx context.Context, postID int64, page, pageSize int32, currentUserID int64) ([]*Comment, int32, error) {
	comments, total, err := uc.commentRepo.GetCommentList(ctx, postID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 填充作者信息和点赞状态
	for _, comment := range comments {
		// 获取作者信息
		author, err := uc.userRepo.GetUserByID(ctx, comment.UserID)
		if err == nil {
			comment.Author = author
		}

		// 获取点赞状态
		if currentUserID > 0 {
			isLiked, _ := uc.commentRepo.IsCommentLiked(ctx, currentUserID, comment.CommentID)
			comment.IsLiked = isLiked
		}
	}

	return comments, total, nil
}

// GetTagList 获取标签列表
func (uc *CommunityUsecase) GetTagList(ctx context.Context) ([]*Tag, error) {
	return uc.tagRepo.GetTagList(ctx)
}

// GetUserPosts 获取用户帖子列表
func (uc *CommunityUsecase) GetUserPosts(ctx context.Context, userID int64, page, pageSize int32) ([]*Post, int32, error) {
	return uc.postRepo.GetUserPosts(ctx, userID, page, pageSize)
}

// GetLikedPosts 获取点赞的帖子列表
func (uc *CommunityUsecase) GetLikedPosts(ctx context.Context, userID int64, page, pageSize int32) ([]*Post, int32, error) {
	return uc.postRepo.GetLikedPosts(ctx, userID, page, pageSize)
}
