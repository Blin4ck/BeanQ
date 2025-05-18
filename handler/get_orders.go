package handler

import (
	"net/http"
	"test/models"
	"test/session"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		session, err := session.Store.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"status": "500", "error": "Ошибка получения сессии"})
			return
		}
		userID, ok := session.Values["user_id"].(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Вы не авторизованы",
			})
			return
		}

		var client models.Client
		err = db.Preload("Orders.Status").Preload("Orders.OrderItems").First(&client, userID).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка загрузки данных",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "200",
			"message": "Заказы успешно загружены",
			"orders":  client.Orders,
		})

	}

}
