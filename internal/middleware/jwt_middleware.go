package middleware

import (
	"coffe/internal/auth"
	"coffe/internal/common"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserContextKey string

type JWTMiddleware struct {
	jwtService *auth.JWTService
	userRepo   UserRepository
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*common.User, error)
}

func NewJWTMiddleware(jwtService *auth.JWTService, userRepo UserRepository) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Получить токен из заголовка
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует заголовок авторизации"})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка авторизации"})
			return
		}

		token := parts[1]

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token не найден"})
			return
		}

		// 2. Провалидировать токен

		claims, err := m.jwtService.ParseJWTToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ошибка чтения Tokena"})
			return
		}
		// 3. Получить userID из claims

		userID, err := uuid.Parse(claims.UserID)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный идентификатор пользователя в токене"})
			return
		}
		// 4. Получить пользователя из базы
		user, err := m.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			return
		}
		ctx.Set("user", user)
		ctx.Set("user_id", user.ID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

func (m *JWTMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInterface, exists := ctx.Get("user")

		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден в контексте"})
			return
		}

		user, ok := userInterface.(*common.User)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип пользователя в контексте"})
			return
		}

		if user.Role == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "У пользователя не найдена роль"})
			return
		}

		hasRole := false

		for _, role := range roles {
			if user.Role.Name == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			return
		}

		ctx.Next()

	}
}

// RequireRole проверяет, что пользователь обладает одной из указанных ролей.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInterface, exists := ctx.Get("user")
		if !exists {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Пользователь не найден"})
			return
		}
		user, ok := userInterface.(*common.User)
		if !ok || user.Role == nil {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Нет роли"})
			return
		}
		for _, role := range roles {
			if user.Role.Name == role {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(403, gin.H{"error": "Недостаточно прав"})
	}
}
