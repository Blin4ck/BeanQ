package handler

import (
	"test/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddPosition(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name        string `json:"name" binding:"required"`
			Price       int    `json:"price" binding:"required"`
			Description string `json:"description" binding:"required"`
			Category    string `json:"category" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"status": "400", "error": "Неверные данные"})
			return
		}
		if input.Name == "" || input.Price < 0 || input.Description == "" || input.Category == "" {
			c.JSON(400, gin.H{"status": "400", "error": "Все поля должны быть заполнены корректно"})
			return
		}
		position := models.Position{
			Name:        input.Name,
			Price:       input.Price,
			Description: input.Description,
			Category:    input.Category,
		}

		if err := db.Create(&position).Error; err != nil {
			c.JSON(500, gin.H{"status": "500", "error": "Ошибка при добавлении позиции"})
			return
		}
		c.JSON(200, gin.H{"status": "200", "message": "Позиция успешно добавлена"})
	}
}
