package router

import (
	"test/handler"
	"test/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	r := gin.Default()

	r.GET("/", handler.Ping)
	r.POST("/register", handler.Register(db))
	r.POST("/login", handler.Login(db))

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(db))
	{
		authorized.GET("/orders", handler.GetOrders(db))
		authorized.POST("/order", handler.CreateOrder(db))
		authorized.GET("/menu", handler.GetMenu(db))
	}
	return r
}
