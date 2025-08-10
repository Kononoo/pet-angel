package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pet-angel/internal/biz"

	"golang.org/x/crypto/bcrypt"
)

// AuthRepo 实现用户认证与资料的数据访问。
// 当 Data.DB 不为空时使用 MySQL；否则使用内存映射（便于本地快速运行）。

type AuthRepo struct {
	data *Data
}

func NewAuthRepo(d *Data) *AuthRepo { return &AuthRepo{data: d} }

func dtoToBiz(u *UserDTO) *biz.User {
	if u == nil {
		return nil
	}
	return &biz.User{
		Id:          u.ID,
		Username:    u.Username,
		Password:    u.Password,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		ModelID:     u.ModelID,
		ModelURL:    u.ModelURL,
		PetName:     u.PetName,
		PetAvatar:   u.PetAvatar,
		PetSex:      u.PetSex,
		Kind:        u.Kind,
		Weight:      u.Weight,
		Hobby:       u.Hobby,
		Description: u.Description,
		Coins:       u.Coins,
		CreatedAt:   time.Unix(u.CreatedAt, 0),
	}
}

// --- MySQL helpers ---

func (r *AuthRepo) getByUsernameSQL(ctx context.Context, username string) (*biz.User, error) {
	row := r.data.DB.QueryRowContext(
		ctx,
		`SELECT id,username,password,nickname,avatar,model_id,model_url,pet_name,pet_avatar,pet_sex,kind,weight,hobby,description,coins,created_at
		 FROM users WHERE username=?`,
		username,
	)
	var u biz.User
	var createdAt time.Time
	err := row.Scan(
		&u.Id, &u.Username, &u.Password, &u.Nickname, &u.Avatar, &u.ModelID, &u.ModelURL,
		&u.PetName, &u.PetAvatar, &u.PetSex, &u.Kind, &u.Weight, &u.Hobby, &u.Description, &u.Coins, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	u.CreatedAt = createdAt
	return &u, nil
}

func (r *AuthRepo) getByIDSQL(ctx context.Context, userID int64) (*biz.User, error) {
	row := r.data.DB.QueryRowContext(
		ctx,
		`SELECT id,username,password,nickname,avatar,model_id,model_url,pet_name,pet_avatar,pet_sex,kind,weight,hobby,description,coins,created_at
		 FROM users WHERE id=?`,
		userID,
	)
	var u biz.User
	var createdAt time.Time
	err := row.Scan(
		&u.Id, &u.Username, &u.Password, &u.Nickname, &u.Avatar, &u.ModelID, &u.ModelURL,
		&u.PetName, &u.PetAvatar, &u.PetSex, &u.Kind, &u.Weight, &u.Hobby, &u.Description, &u.Coins, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	u.CreatedAt = createdAt
	return &u, nil
}

func (r *AuthRepo) createSQL(ctx context.Context, user *biz.User) (int64, error) {
	// 密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	// 默认模型与URL
	if user.ModelURL == "" {
		user.ModelURL = "/models/Dog_1.glb"
	}
	res, err := r.data.DB.ExecContext(
		ctx,
		`INSERT INTO users(
		 username,password,nickname,avatar,model_id,model_url,pet_name,pet_avatar,pet_sex,kind,weight,hobby,description,coins,created_at,updated_at
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW(),NOW())`,
		user.Username, string(hash), user.Nickname, user.Avatar, user.ModelID, user.ModelURL, user.PetName, user.PetAvatar, user.PetSex, user.Kind, user.Weight, user.Hobby, user.Description, 0,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepo) updateInfoSQL(ctx context.Context, user *biz.User) error {
	_, err := r.data.DB.ExecContext(
		ctx,
		`UPDATE users
		 SET nickname=IF(?<>'',?,nickname),
		     avatar=IF(?<>'',?,avatar),
		     model_id=IF(?<>0,?,model_id),
		     model_url=IF(?<>'',?,model_url),
		     pet_name=IF(?<>'',?,pet_name),
		     pet_avatar=IF(?<>'',?,pet_avatar),
		     pet_sex=IF(?<>0,?,pet_sex),
		     kind=IF(?<>'',?,kind),
		     weight=IF(?<>0,?,weight),
		     hobby=IF(?<>'',?,hobby),
		     description=IF(?<>'',?,description),
		     coins=IF(?<>0,?,coins),
		     updated_at=NOW()
		 WHERE id=?`,
		user.Nickname, user.Nickname,
		user.Avatar, user.Avatar,
		user.ModelID, user.ModelID,
		user.ModelURL, user.ModelURL,
		user.PetName, user.PetName,
		user.PetAvatar, user.PetAvatar,
		user.PetSex, user.PetSex,
		user.Kind, user.Kind,
		user.Weight, user.Weight,
		user.Hobby, user.Hobby,
		user.Description, user.Description,
		user.Coins, user.Coins,
		user.Id,
	)
	return err
}

func (r *AuthRepo) getModelPathSQL(ctx context.Context, modelID int64) (string, error) {
	row := r.data.DB.QueryRowContext(ctx, `SELECT path FROM pet_models WHERE id=?`, modelID)
	var path string
	if err := row.Scan(&path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", biz.ErrAvatarNotFound
		}
		return "", err
	}
	return path, nil
}

// --- interface methods ---

func (r *AuthRepo) GetByUsername(ctx context.Context, username string) (*biz.User, error) {
	if r.data.DB != nil {
		return r.getByUsernameSQL(ctx, username)
	}
	// in-memory fallback
	r.data.mu.RLock()
	defer r.data.mu.RUnlock()
	u, ok := r.data.userByUsername[username]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	return dtoToBiz(u), nil
}

func (r *AuthRepo) GetByID(ctx context.Context, userID int64) (*biz.User, error) {
	if r.data.DB != nil {
		return r.getByIDSQL(ctx, userID)
	}
	r.data.mu.RLock()
	defer r.data.mu.RUnlock()
	u, ok := r.data.userByID[userID]
	if !ok {
		return nil, biz.ErrUserNotFound
	}
	return dtoToBiz(u), nil
}

func (r *AuthRepo) Create(ctx context.Context, user *biz.User) (int64, error) {
	if r.data.DB != nil {
		return r.createSQL(ctx, user)
	}
	// in-memory fallback
	r.data.mu.Lock()
	defer r.data.mu.Unlock()
	if _, exists := r.data.userByUsername[user.Username]; exists {
		return 0, errors.New("username exists")
	}
	id := r.data.nextUserID
	r.data.nextUserID++
	// 内存模式下也保存 bcrypt 哈希，保持与 MySQL 行为一致
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	if user.ModelURL == "" {
		user.ModelURL = "/models/Dog_1.glb"
	}
	d := &UserDTO{
		ID:          id,
		Username:    user.Username,
		Password:    string(hash),
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		ModelID:     user.ModelID,
		ModelURL:    user.ModelURL,
		PetName:     user.PetName,
		PetAvatar:   user.PetAvatar,
		PetSex:      user.PetSex,
		Kind:        user.Kind,
		Weight:      user.Weight,
		Hobby:       user.Hobby,
		Description: user.Description,
		Coins:       0,
		CreatedAt:   time.Now().Unix(),
	}
	r.data.userByID[id] = d
	r.data.userByUsername[user.Username] = d
	return id, nil
}

func (r *AuthRepo) UpdateInfo(ctx context.Context, user *biz.User) error {
	if r.data.DB != nil {
		return r.updateInfoSQL(ctx, user)
	}
	// in-memory fallback
	r.data.mu.Lock()
	defer r.data.mu.Unlock()
	u, ok := r.data.userByID[user.Id]
	if !ok {
		return biz.ErrUserNotFound
	}
	if user.Nickname != "" {
		u.Nickname = user.Nickname
	}
	if user.Avatar != "" {
		u.Avatar = user.Avatar
	}
	if user.ModelID != 0 {
		u.ModelID = user.ModelID
	}
	if user.ModelURL != "" {
		u.ModelURL = user.ModelURL
	}
	if user.PetName != "" {
		u.PetName = user.PetName
	}
	if user.PetAvatar != "" {
		u.PetAvatar = user.PetAvatar
	}
	if user.PetSex != 0 {
		u.PetSex = user.PetSex
	}
	if user.Kind != "" {
		u.Kind = user.Kind
	}
	if user.Weight != 0 {
		u.Weight = user.Weight
	}
	if user.Hobby != "" {
		u.Hobby = user.Hobby
	}
	if user.Description != "" {
		u.Description = user.Description
	}
	// 允许直接修改金币（demo模式）
	if user.Coins != 0 {
		u.Coins = user.Coins
	}
	return nil
}

func (r *AuthRepo) UpdateCoins(ctx context.Context, userID int64, delta int32) error {
	if r.data.DB != nil {
		_, err := r.data.DB.ExecContext(ctx, `UPDATE users SET coins=GREATEST(0,coins+?), updated_at=NOW() WHERE id=?`, delta, userID)
		return err
	}
	r.data.mu.Lock()
	defer r.data.mu.Unlock()
	u, ok := r.data.userByID[userID]
	if !ok {
		return biz.ErrUserNotFound
	}
	u.Coins += delta
	if u.Coins < 0 {
		u.Coins = 0
	}
	return nil
}

func (r *AuthRepo) GetModelPath(ctx context.Context, modelID int64) (string, error) {
	if r.data.DB != nil {
		return r.getModelPathSQL(ctx, modelID)
	}
	// 内存模式：无模型表，返回默认
	return "/models/Dog_1.glb", nil
}
