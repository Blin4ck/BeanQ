package handler

import (
	"test/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			FirstName string `json:"first_name" binding:"required"`
			LastName  string `json:"last_name" binding:"required"`
			Email     string `json:"email" binding:"required,email"`
			Pasword   string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Неверные данные"})
			return
		}
		var existingClient models.Client

		if input.FirstName == "" || input.LastName == "" {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Имя и фамилия не могут быть пустыми"})
			return
		}
		if input.Email == "" {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Email не может быть пустым"})
			return
		}
		if input.Pasword == "" {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Пароль не может быть пустым"})
			return
		}
		if err := db.Where("email = ?", input.Email).First(&existingClient).Error; err == nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Клиент с таким email уже существует"})
			return
		}

		if len(input.Pasword) < 6 {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Пароль должен содержать не менее 6 символов"})
			return
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Pasword), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"status": "500", "error": "Ошибка при хешировании пароля"})
			return
		}

		client := models.Client{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Pasword:   string(hashPassword),
		}

		if err := db.Create(&client).Error; err != nil {
			c.AbortWithStatusJSON(500, gin.H{"status": "500", "error": "Ошибка при создании клиента"})
			return
		}

		c.JSON(200, gin.H{
			"status":  "200",
			"message": "Клиент успешно создан",
			"client":  client.ID,
		})

	}
}
