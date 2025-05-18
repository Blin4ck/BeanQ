package handler

import (
	"test/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalPrice int
		var orderItems []models.OrderItem

		userID, _ := c.Get("user")
		client := userID.(models.Client)

		var input struct {
			Items []struct {
				PositionID uint `json:"position_id" binding:"required"`
				Quantity   uint `json:"quantity" binding:"required,min=1"`
			} `json:"items" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"status":  "400",
				"error":   "Неверные данные",
				"details": err.Error(),
			})
			return
		}

		for _, item := range input.Items {
			var position models.Position
			if err := db.First(&position, item.PositionID).Error; err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"status":      "400",
					"error":       "Позиция не найдена",
					"position_id": item.PositionID,
				})
				return
			}

			itemPrice := position.Price * int(item.Quantity)
			totalPrice += itemPrice

			orderItems = append(orderItems, models.OrderItem{
				PositionID: item.PositionID,
				Position:   position,
				Quantity:   int(item.Quantity),
				Price:      itemPrice,
			})
		}

		order := models.Order{
			ClientID:   client.ID,
			Price:      totalPrice,
			OrderItems: orderItems,
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			for i := range orderItems {
				orderItems[i].OrderID = order.ID
			}

			if err := tx.CreateInBatches(orderItems, len(orderItems)).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"status":  "500",
				"error":   "Не удалось создать заказ",
				"details": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "200",
			"message": "Заказ успешно создан",
			"order": gin.H{
				"id":         order.ID,
				"created_at": order.CreatedAt,
				"price":      order.Price,
				"items":      order.OrderItems,
			},
		})
	}
}
