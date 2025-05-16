package handler

import (
	"log"
	"net/http"
	"test/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PostAuthentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Не удалось обработать запрос"})
			return
		}
		if input.Email == "" {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Email не может быть пустым"})
			return
		}
		if input.Password == "" {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Пароль не может быть пустым"})
			return
		}

		var client models.Client
		if result := db.Where("email = ?", input.Email).First(&client); result.Error != nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Клиент с таким email не найден"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(client.Pasword), []byte(input.Password)); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"status": "400", "error": "Неверный пароль"})
			return
		}

		session, err := Store.Get(c.Request, "session")
		if err != nil {
			log.Printf("Ошибка получения сессии: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка аутентификации"})
			return
		}
		session.Values["user_id"] = client.ID
		session.Options = &sessions.Options{
			HttpOnly: true,
			MaxAge:   86400 * 7,
			Secure:   false,
		}

		if err := session.Save(c.Request, c.Writer); err != nil {
			log.Printf("Ошибка сохранения сессии: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания сессии"})
			return
		}

		c.JSON(200, gin.H{"status": "200", "message": "Успешная аутентификация"})

	}
}
