package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var globalSecret string

// SetGlobalSecret 设置全局JWT密钥（由服务初始化时设置）
func SetGlobalSecret(secret string) { globalSecret = secret }

// Claims 自定义Claims
type Claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

// Sign 生成 JWT 字符串
func Sign(secret string, userID int64, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(ttl)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	return s, exp, err
}

// Parse 使用指定secret解析
func Parse(secret, tokenString string) (*Claims, error) {
	tk, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tk.Claims.(*Claims); ok && tk.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseSecretFromCtx 使用全局密钥解析
func ParseSecretFromCtx(_ interface{}, tokenString string) (*Claims, error) {
	if globalSecret == "" {
		return nil, errors.New("jwt secret not initialized")
	}
	return Parse(globalSecret, tokenString)
}

// FromAuthHeader 提取 Bearer Token
func FromAuthHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("empty authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}
