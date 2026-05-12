package middleware

import (
	"api-gateway/pkg"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware проверяет access токен и сохраняет user_id в контексте
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		// Парсим токен
		claims, err := pkg.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Проверяем, что это access токен (не refresh)
		if claims.Type != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token type",
			})
			return
		}

		// Сохраняем user_id в контексте для последующих обработчиков
		c.Set("user_id", claims.UUID)
		c.Next()
	}
}

// AdminMiddleware проверяет, что пользователь имеет роль admin
// Требует, чтобы JWTAuthMiddleware был выполнен до него
func AdminMiddleware(userServiceClient interface {
	GetUserRole(ctx context.Context, userID string) (string, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем user_id из контекста (установлен JWTAuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// Получаем роль пользователя из user-service
		role, err := userServiceClient.GetUserRole(c.Request.Context(), userID.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get user role",
			})
			return
		}

		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			return
		}

		c.Next()
	}
}
