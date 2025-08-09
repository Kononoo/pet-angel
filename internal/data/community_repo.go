package data

import (
	"context"
	"strings"
	"time"

	"pet-angel/internal/biz"

	"gorm.io/gorm/clause"
)

// 表模型定义（仅含社区模块用到的字段）

type CategoryModel struct {
	ID   int64  `gorm:"column:id;primaryKey"` // 分类ID
	Name string `gorm:"column:name"`          // 名称
}

func (CategoryModel) TableName() string { return "categories" }

type PostModel struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	UserID       int64  `gorm:"column:user_id"`
	CategoryID   int64  `gorm:"column:category_id"`
	Title        string `gorm:"column:title"`
	Content      string `gorm:"column:content;type:text"`
	Type         int32  `gorm:"column:type"` // 0图文 1视频
	ImageUrls    string `gorm:"column:image_urls"`
	VideoUrl     string `gorm:"column:video_url"`
	CoverUrl     string `gorm:"column:cover_url"`
	Locate       string `gorm:"column:locate"`
	Tags         string `gorm:"column:tags"`
	LikedCount   int32  `gorm:"column:liked_count"`
	CommentCount int32  `gorm:"column:comment_count"`
	CreatedAt    string `gorm:"column:created_at"` // DATE_FORMAT 生成字符串
	IsPrivate    int32  `gorm:"column:is_private"`
}

func (PostModel) TableName() string { return "posts" }

type CommentModel struct {
	ID         int64  `gorm:"column:id;primaryKey"`
	PostID     int64  `gorm:"column:post_id"`
	UserID     int64  `gorm:"column:user_id"`
	Content    string `gorm:"column:content"`
	LikedCount int32  `gorm:"column:liked_count"`
	CreatedAt  string `gorm:"column:created_at"`
}

func (CommentModel) TableName() string { return "comments" }

type LikeModel struct {
	UserID     int64 `gorm:"column:user_id;primaryKey"`
	TargetType int32 `gorm:"column:target_type;primaryKey"` // 0帖子 1评论
	TargetID   int64 `gorm:"column:target_id;primaryKey"`
}

func (LikeModel) TableName() string { return "likes" }

// CommunityRepoImpl 实现 CommunityRepo

type CommunityRepoImpl struct{ data *Data }

func NewCommunityRepo(d *Data) *CommunityRepoImpl { return &CommunityRepoImpl{data: d} }

// 工具：拆分/合并逗号分隔字符串
func splitCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		pp := strings.TrimSpace(p)
		if pp != "" {
			out = append(out, pp)
		}
	}
	return out
}

func joinCSV(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return strings.Join(arr, ",")
}

func parseDT(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	loc := time.Local
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	return t
}

func (r *CommunityRepoImpl) ListCategories(ctx context.Context) ([]*biz.Category, error) {
	if r.data.Gorm == nil {
		return []*biz.Category{}, nil
	}
	var rows []CategoryModel
	if err := r.data.Gorm.WithContext(ctx).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*biz.Category, 0, len(rows))
	for _, c := range rows {
		cc := c
		out = append(out, &biz.Category{
			ID:   cc.ID,
			Name: cc.Name,
		})
	}
	return out, nil
}

func (r *CommunityRepoImpl) ListPosts(ctx context.Context, viewerID int64, categoryID int64, postType int32, sort string, page, pageSize int32) (int32, []*biz.CommunityPost, error) {
	if r.data.Gorm == nil {
		return 0, []*biz.CommunityPost{}, nil
	}
	// 1) 先取 posts（不做 JOIN）
	q := r.data.Gorm.WithContext(ctx).
		Table("posts p").
		Select("p.id,p.user_id,p.category_id,p.title,p.content,p.type,p.image_urls,p.video_url,p.cover_url,p.locate,p.tags,p.liked_count,p.comment_count,DATE_FORMAT(p.created_at,'%Y-%m-%d %H:%i:%s') as created_at,p.is_private")
	if categoryID > 0 {
		q = q.Where("p.category_id=?", categoryID)
	}
	if postType == 0 || postType == 1 {
		q = q.Where("p.type=?", postType)
	}
	orderBy := "p.id DESC"
	if sort == "liked" {
		orderBy = "p.liked_count DESC,p.id DESC"
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	offset := (page - 1) * pageSize
	var rows []PostModel
	if err := q.
		Order(orderBy).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(&rows).Error; err != nil {
		return 0, nil, err
	}
	// 2) 批量取作者信息（users）
	userIds := make([]int64, 0, len(rows))
	for _, r0 := range rows {
		userIds = append(userIds, r0.UserID)
	}
	userMap := map[int64]struct{ Nickname, Avatar string }{}
	if len(userIds) > 0 {
		var urows []struct {
			ID       int64  `gorm:"column:id"`
			Nickname string `gorm:"column:nickname"`
			Avatar   string `gorm:"column:avatar"`
		}
		_ = r.data.Gorm.WithContext(ctx).
			Table("users").
			Select("id,nickname,avatar").
			Where("id IN ?", userIds).
			Scan(&urows).Error
		for _, u := range urows {
			userMap[u.ID] = struct{ Nickname, Avatar string }{u.Nickname, u.Avatar}
		}
	}
	// 3) liked/followed 状态
	likedSet := map[int64]bool{}
	followSet := map[int64]bool{}
	if viewerID > 0 && len(rows) > 0 {
		postIds := make([]int64, 0, len(rows))
		for _, r0 := range rows {
			postIds = append(postIds, r0.ID)
		}
		var likeRows []struct{ TargetID int64 }
		if err := r.data.Gorm.WithContext(ctx).
			Table("likes").
			Select("target_id").
			Where("user_id=? AND target_type=0", viewerID).
			Where("target_id IN ?", postIds).
			Find(&likeRows).Error; err == nil {
			for _, lr := range likeRows {
				likedSet[lr.TargetID] = true
			}
		}
		var followRows []struct{ FolloweeID int64 }
		if err := r.data.Gorm.WithContext(ctx).
			Table("user_follows").
			Select("followee_id").
			Where("follower_id=?", viewerID).
			Where("followee_id IN ?", userIds).
			Find(&followRows).Error; err == nil {
			for _, fr := range followRows {
				followSet[fr.FolloweeID] = true
			}
		}
	}
	// 4) 组装
	out := make([]*biz.CommunityPost, 0, len(rows))
	for _, p := range rows {
		u := userMap[p.UserID]
		out = append(out, &biz.CommunityPost{
			ID:           p.ID,
			User:         biz.UserBrief{Id: p.UserID, Nickname: u.Nickname, Avatar: u.Avatar},
			CategoryID:   p.CategoryID,
			Title:        p.Title,
			Content:      p.Content,
			PostType:     p.Type,
			ImageUrls:    splitCSV(p.ImageUrls),
			VideoUrl:     p.VideoUrl,
			CoverUrl:     p.CoverUrl,
			Locate:       p.Locate,
			Tags:         splitCSV(p.Tags),
			LikedCount:   p.LikedCount,
			CommentCount: p.CommentCount,
			CreatedAt:    parseDT(p.CreatedAt),
			IsLiked:      likedSet[p.ID],
			IsPrivate:    p.IsPrivate == 1,
		})
	}
	return int32(total), out, nil
}

func (r *CommunityRepoImpl) GetPostDetail(ctx context.Context, viewerID, postID int64) (*biz.CommunityPost, error) {
	if r.data.Gorm == nil {
		return nil, nil
	}
	// 1) 取帖子
	var p PostModel
	if err := r.data.Gorm.WithContext(ctx).
		Table("posts p").
		Select("p.id,p.user_id,p.category_id,p.title,p.content,p.type,p.image_urls,p.video_url,p.cover_url,p.locate,p.tags,p.liked_count,p.comment_count,DATE_FORMAT(p.created_at,'%Y-%m-%d %H:%i:%s') as created_at,p.is_private").
		Where("p.id=?", postID).
		Scan(&p).Error; err != nil {
		return nil, err
	}
	// 2) 取作者信息
	var urow struct {
		Nickname string
		Avatar   string
	}
	_ = r.data.Gorm.WithContext(ctx).
		Table("users").
		Select("nickname,avatar").
		Where("id=?", p.UserID).
		Scan(&urow).Error
	// 3) 是否点赞
	isLiked := false
	if viewerID > 0 {
		var cnt int64
		_ = r.data.Gorm.WithContext(ctx).
			Table("likes").
			Where("user_id=? AND target_type=0 AND target_id=?", viewerID, postID).
			Count(&cnt).Error
		isLiked = cnt > 0
	}
	return &biz.CommunityPost{
		ID:           p.ID,
		User:         biz.UserBrief{Id: p.UserID, Nickname: urow.Nickname, Avatar: urow.Avatar},
		CategoryID:   p.CategoryID,
		Title:        p.Title,
		Content:      p.Content,
		PostType:     p.Type,
		ImageUrls:    splitCSV(p.ImageUrls),
		VideoUrl:     p.VideoUrl,
		CoverUrl:     p.CoverUrl,
		Locate:       p.Locate,
		Tags:         splitCSV(p.Tags),
		LikedCount:   p.LikedCount,
		CommentCount: p.CommentCount,
		CreatedAt:    parseDT(p.CreatedAt),
		IsLiked:      isLiked,
		IsPrivate:    p.IsPrivate == 1,
	}, nil
}

func (r *CommunityRepoImpl) CreatePost(ctx context.Context, userID int64, p *biz.CommunityPost) (int64, error) {
	if r.data.Gorm == nil {
		return 0, nil
	}
	row := &PostModel{
		UserID:     userID,
		CategoryID: p.CategoryID,
		Title:      p.Title,
		Content:    p.Content,
		Type:       p.PostType,
		ImageUrls:  joinCSV(p.ImageUrls),
		VideoUrl:   p.VideoUrl,
		CoverUrl:   p.CoverUrl,
		Locate:     p.Locate,
		Tags:       joinCSV(p.Tags),
		IsPrivate:  0,
	}
	if p.IsPrivate {
		row.IsPrivate = 1
	}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (r *CommunityRepoImpl) LikePost(ctx context.Context, userID, postID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	// 幂等：仅当插入发生时才自增计数
	tx := r.data.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&LikeModel{UserID: userID, TargetType: 0, TargetID: postID})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected > 0 {
		return r.data.Gorm.WithContext(ctx).
			Exec("UPDATE posts SET liked_count=liked_count+1 WHERE id=?", postID).Error
	}
	return nil
}

func (r *CommunityRepoImpl) UnlikePost(ctx context.Context, userID, postID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	del := r.data.Gorm.WithContext(ctx).
		Where("user_id=? AND target_type=0 AND target_id=?", userID, postID).
		Delete(&LikeModel{})
	if del.Error != nil {
		return del.Error
	}
	if del.RowsAffected > 0 {
		return r.data.Gorm.WithContext(ctx).
			Exec("UPDATE posts SET liked_count=GREATEST(liked_count-1,0) WHERE id=?", postID).Error
	}
	return nil
}

func (r *CommunityRepoImpl) ListComments(ctx context.Context, viewerID, postID int64, page, pageSize int32) (int32, []*biz.CommunityComment, error) {
	if r.data.Gorm == nil {
		return 0, []*biz.CommunityComment{}, nil
	}
	var total int64
	if err := r.data.Gorm.WithContext(ctx).
		Model(&CommentModel{}).
		Where("post_id=?", postID).
		Count(&total).Error; err != nil {
		return 0, nil, err
	}
	offset := (page - 1) * pageSize
	var rows []CommentModel
	if err := r.data.Gorm.WithContext(ctx).
		Table("comments c").
		Select("c.id,c.post_id,c.user_id,c.content,c.liked_count,DATE_FORMAT(c.created_at,'%Y-%m-%d %H:%i:%s') as created_at").
		Where("c.post_id=?", postID).
		Order("c.id DESC").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Scan(&rows).Error; err != nil {
		return 0, nil, err
	}
	// 批量取用户昵称、头像
	userIds := make([]int64, 0, len(rows))
	for _, r0 := range rows {
		userIds = append(userIds, r0.UserID)
	}
	userMap := map[int64]struct{ Nickname, Avatar string }{}
	if len(userIds) > 0 {
		var urows []struct {
			ID       int64
			Nickname string
			Avatar   string
		}
		_ = r.data.Gorm.WithContext(ctx).
			Table("users").
			Select("id,nickname,avatar").
			Where("id IN ?", userIds).
			Scan(&urows).Error
		for _, u := range urows {
			userMap[u.ID] = struct{ Nickname, Avatar string }{u.Nickname, u.Avatar}
		}
	}
	likedSet := map[int64]bool{}
	if viewerID > 0 && len(rows) > 0 {
		ids := make([]int64, 0, len(rows))
		for _, r0 := range rows {
			ids = append(ids, r0.ID)
		}
		var likeRows []struct{ TargetID int64 }
		if err := r.data.Gorm.WithContext(ctx).
			Table("likes").
			Select("target_id").
			Where("user_id=? AND target_type=1", viewerID).
			Where("target_id IN ?", ids).
			Find(&likeRows).Error; err == nil {
			for _, lr := range likeRows {
				likedSet[lr.TargetID] = true
			}
		}
	}
	out := make([]*biz.CommunityComment, 0, len(rows))
	for _, c := range rows {
		u := userMap[c.UserID]
		out = append(out, &biz.CommunityComment{
			ID:         c.ID,
			User:       biz.UserBrief{Id: c.UserID, Nickname: u.Nickname, Avatar: u.Avatar},
			Content:    c.Content,
			LikedCount: c.LikedCount,
			CreatedAt:  parseDT(c.CreatedAt),
			IsLiked:    likedSet[c.ID],
		})
	}
	return int32(total), out, nil
}

func (r *CommunityRepoImpl) CreateComment(ctx context.Context, userID, postID int64, content string) (int64, error) {
	if r.data.Gorm == nil {
		return 0, nil
	}
	row := &CommentModel{PostID: postID, UserID: userID, Content: content}
	if err := r.data.Gorm.WithContext(ctx).Create(row).Error; err != nil {
		return 0, err
	}
	_ = r.data.Gorm.WithContext(ctx).
		Exec("UPDATE posts SET comment_count=comment_count+1 WHERE id=?", postID).Error
	return row.ID, nil
}

func (r *CommunityRepoImpl) LikeComment(ctx context.Context, userID, commentID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	// 幂等：仅当插入发生时才自增
	tx := r.data.Gorm.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&LikeModel{UserID: userID, TargetType: 1, TargetID: commentID})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected > 0 {
		return r.data.Gorm.WithContext(ctx).
			Exec("UPDATE comments SET liked_count=liked_count+1 WHERE id=?", commentID).Error
	}
	return nil
}

func (r *CommunityRepoImpl) UnlikeComment(ctx context.Context, userID, commentID int64) error {
	if r.data.Gorm == nil {
		return nil
	}
	del := r.data.Gorm.WithContext(ctx).
		Where("user_id=? AND target_type=1 AND target_id=?", userID, commentID).
		Delete(&LikeModel{})
	if del.Error != nil {
		return del.Error
	}
	if del.RowsAffected > 0 {
		return r.data.Gorm.WithContext(ctx).
			Exec("UPDATE comments SET liked_count=GREATEST(liked_count-1,0) WHERE id=?", commentID).Error
	}
	return nil
}
