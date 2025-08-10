package biz

import "errors"

// 用户相关错误
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrCannotFollowSelf  = errors.New("cannot follow self")
)

// 宠物相关错误
var (
	ErrPetNotFound = errors.New("pet not found")
)

// 帖子相关错误
var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostDeleted  = errors.New("post has been deleted")
)

// 评论相关错误
var (
	ErrCommentNotFound = errors.New("comment not found")
)

// 虚拟形象相关错误
var (
	ErrAvatarNotFound = errors.New("avatar not found")
	ErrAvatarLocked   = errors.New("avatar is locked")
)

// 道具相关错误
var (
	ErrPropNotFound      = errors.New("prop not found")
	ErrInsufficientCoins = errors.New("insufficient coins")
	ErrPropNotOwned      = errors.New("prop not owned")
)

// 小纸条相关错误
var (
	ErrMessageNotFound = errors.New("message not found")
	ErrMessageLocked   = errors.New("message is locked")
)

// 通用错误
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrDatabaseError    = errors.New("database error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNotImplemented   = errors.New("method not implemented")
)
