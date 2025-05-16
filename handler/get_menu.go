package handler

import (
	"test/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMenu(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menus models.Menu
		db.Preload("Position").Find(&menus)
	}
}
