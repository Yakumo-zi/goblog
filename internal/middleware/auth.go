package middleware

import (
	"strings"

	"goblog/internal/domain"
	"goblog/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware JWT认证中间件
type AuthMiddleware struct {
	authService domain.AuthService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从Header获取Token
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "缺少Authorization header")
			}

			// 检查Bearer前缀
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return response.Unauthorized(c, "无效的token格式")
			}

			// 验证token
			username, err := m.authService.ValidateToken(c.Request().Context(), authHeader)
			if err != nil {
				return response.Unauthorized(c, "token验证失败")
			}

			// 将用户信息存储到context中
			c.Set("username", username)
			return next(c)
		}
	}
}
