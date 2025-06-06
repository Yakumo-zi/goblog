package service

import (
	"context"
	"strings"
	"time"

	"goblog/internal/config"
	"goblog/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务实现
type AuthService struct {
	config *config.Config
}

// Claims JWT载荷
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewAuthService 创建认证服务
func NewAuthService(cfg *config.Config) domain.AuthService {
	return &AuthService{config: cfg}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// 验证用户名
	if req.Username != s.config.Admin.Username {
		return nil, domain.ErrUnauthorized
	}

	// 验证密码 - 在实际项目中应该从数据库获取用户信息
	expectedPasswordHash := "$2a$10$gYH0.79TG8qC6hC4zVamLOuJt/77ZiHEqoeinlos3.jz6u2DDqg42"
	err := bcrypt.CompareHashAndPassword([]byte(expectedPasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	// 生成JWT token
	token, err := s.GenerateToken(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token:    token,
		Username: req.Username,
	}, nil
}

// ValidateToken 验证token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	// 移除Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", domain.ErrUnauthorized
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", domain.ErrUnauthorized
}

// GenerateToken 生成JWT token
func (s *AuthService) GenerateToken(ctx context.Context, username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}
