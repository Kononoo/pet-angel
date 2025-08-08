package biz

import (
	"context"
	"time"
)

// User 用户实体
type User struct {
	UserID      int64     `json:"user_id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Gender      string    `json:"gender"`
	Region      string    `json:"region"`
	Partner     string    `json:"partner"`
	CoinBalance int32     `json:"coin_balance"`
	TotalEarned int32     `json:"total_earned"`
	TotalSpent  int32     `json:"total_spent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Pet 宠物实体
type Pet struct {
	PetID           int64     `json:"pet_id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	Species         string    `json:"species"`
	Gender          string    `json:"gender"`
	Weight          float64   `json:"weight"`
	Hobbies         []string  `json:"hobbies"`
	Birthday        time.Time `json:"birthday"`
	AdoptionDate    time.Time `json:"adoption_date"`
	PassedAwayDate  time.Time `json:"passed_away_date"`
	MemorialWords   string    `json:"memorial_words"`
	AvatarID        int64     `json:"avatar_id"`
	BackgroundImage string    `json:"background_image"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserRepo 用户仓储接口
type UserRepo interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	FollowUser(ctx context.Context, userID, targetUserID int64) error
	UnfollowUser(ctx context.Context, userID, targetUserID int64) error
	GetFollowingList(ctx context.Context, userID int64, page, pageSize int32) ([]*User, int32, error)
	GetFollowersList(ctx context.Context, userID int64, page, pageSize int32) ([]*User, int32, error)
	IsFollowing(ctx context.Context, userID, targetUserID int64) (bool, error)
	GetUserStats(ctx context.Context, userID int64) (followingCount, followersCount, postsCount int32, err error)
}

// PetRepo 宠物仓储接口
type PetRepo interface {
	CreatePet(ctx context.Context, pet *Pet) error
	GetPetByID(ctx context.Context, petID int64) (*Pet, error)
	GetPetsByUserID(ctx context.Context, userID int64) ([]*Pet, error)
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, petID int64) error
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo    UserRepo
	petRepo PetRepo
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, petRepo PetRepo) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		petRepo: petRepo,
	}
}

// Register 用户注册
func (uc *UserUsecase) Register(ctx context.Context, username, password, nickname, email string) (*User, error) {
	// 检查用户名是否已存在
	existingUser, err := uc.repo.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 创建新用户
	user := &User{
		Username:  username,
		Password:  password, // 实际应用中需要加密
		Nickname:  nickname,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (uc *UserUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	user, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 验证密码（实际应用中需要加密比较）
	if user.Password != password {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// GetUserInfo 获取用户信息
func (uc *UserUsecase) GetUserInfo(ctx context.Context, userID int64) (*User, []*Pet, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	pets, err := uc.petRepo.GetPetsByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return user, pets, nil
}

// UpdateUserInfo 更新用户信息
func (uc *UserUsecase) UpdateUserInfo(ctx context.Context, userID int64, nickname, avatar, gender, region, partner string) (*User, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}
	if region != "" {
		user.Region = region
	}
	if partner != "" {
		user.Partner = partner
	}

	user.UpdatedAt = time.Now()

	err = uc.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FollowUser 关注用户
func (uc *UserUsecase) FollowUser(ctx context.Context, userID, targetUserID int64) error {
	if userID == targetUserID {
		return ErrCannotFollowSelf
	}

	// 检查目标用户是否存在
	_, err := uc.repo.GetUserByID(ctx, targetUserID)
	if err != nil {
		return ErrUserNotFound
	}

	return uc.repo.FollowUser(ctx, userID, targetUserID)
}

// UnfollowUser 取消关注
func (uc *UserUsecase) UnfollowUser(ctx context.Context, userID, targetUserID int64) error {
	return uc.repo.UnfollowUser(ctx, userID, targetUserID)
}

// GetFollowingList 获取关注列表
func (uc *UserUsecase) GetFollowingList(ctx context.Context, userID int64, page, pageSize int32) ([]*User, int32, error) {
	return uc.repo.GetFollowingList(ctx, userID, page, pageSize)
}

// GetFollowersList 获取粉丝列表
func (uc *UserUsecase) GetFollowersList(ctx context.Context, userID int64, page, pageSize int32) ([]*User, int32, error) {
	return uc.repo.GetFollowersList(ctx, userID, page, pageSize)
}

// CreatePet 创建宠物
func (uc *UserUsecase) CreatePet(ctx context.Context, userID int64, name, species, gender string, weight float64, hobbies []string) (*Pet, error) {
	pet := &Pet{
		UserID:    userID,
		Name:      name,
		Species:   species,
		Gender:    gender,
		Weight:    weight,
		Hobbies:   hobbies,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.petRepo.CreatePet(ctx, pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

// GetUserPets 获取用户宠物列表
func (uc *UserUsecase) GetUserPets(ctx context.Context, userID int64) ([]*Pet, error) {
	return uc.petRepo.GetPetsByUserID(ctx, userID)
}

// UpdatePet 更新宠物信息
func (uc *UserUsecase) UpdatePet(ctx context.Context, petID int64, updates map[string]interface{}) (*Pet, error) {
	pet, err := uc.petRepo.GetPetByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok {
		pet.Name = name
	}
	if species, ok := updates["species"].(string); ok {
		pet.Species = species
	}
	if gender, ok := updates["gender"].(string); ok {
		pet.Gender = gender
	}
	if weight, ok := updates["weight"].(float64); ok {
		pet.Weight = weight
	}
	if hobbies, ok := updates["hobbies"].([]string); ok {
		pet.Hobbies = hobbies
	}
	if memorialWords, ok := updates["memorial_words"].(string); ok {
		pet.MemorialWords = memorialWords
	}
	if avatarID, ok := updates["avatar_id"].(int64); ok {
		pet.AvatarID = avatarID
	}
	if backgroundImage, ok := updates["background_image"].(string); ok {
		pet.BackgroundImage = backgroundImage
	}

	pet.UpdatedAt = time.Now()

	err = uc.petRepo.UpdatePet(ctx, pet)
	if err != nil {
		return nil, err
	}

	return pet, nil
}
