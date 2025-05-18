package middleware

import (
	"test/session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ssession, err := session.Store.Get(ctx.Request, "session")
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"status": "500",
				"error":  "Ошибка получения сессии",
			})
			return
		}

		userId, ok := ssession.Values["user_id"].(uint)
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{
				"status": "401",
				"error":  "Пользователь не авторизован",
			})
			return
		}

		ctx.Set("user", userId)

		ctx.Next()

	}
}
