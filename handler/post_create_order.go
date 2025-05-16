package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ssesion, err := Store.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"status": "500", "error": "Ошибка получения сессии"})
			return
		}
		userID, ok := ssesion.Values["user_id"].(uint)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"status": "401", "error": "Пользователь не авторизован"})
			return
		}

		var input struct {
			Item []struct {
				PositionID uint `json:"position_id" binding:"required"`
				Quantity   uint `json:"quantity" binding:"required"`
			}
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Неверные данные"})
			return
		}
		// дописать
	}
}
