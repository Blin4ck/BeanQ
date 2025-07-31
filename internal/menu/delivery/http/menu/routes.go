package http

import (
	"coffe/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes настраивает все маршруты для модуля меню
func SetupMenuRoutes(router *gin.RouterGroup, handler *MenuHandler, middleware *middleware.JWTMiddleware) {
	// Публичные маршруты меню
	setupPublicMenuRoutes(router, handler)

	// Админские маршруты
	setupAdminMenuRoutes(router, handler, middleware)
}

// настраивает публичные маршруты меню
func setupPublicMenuRoutes(router *gin.RouterGroup, handler *MenuHandler) {
	menu := router.Group("/menu")
	{
		// Просмотр меню
		menu.GET("", handler.GetAllMenuItems)
		menu.GET("/:id", handler.GetMenuItemByID)
		menu.GET("/category/:category", handler.GetMenuItemsByCategory)

		// Поиск и фильтрация
		menu.GET("/search", handler.SearchMenuItems)
		menu.GET("/available", handler.GetAvailableItems)
	}
}

// настраивает админские маршруты для управления меню
func setupAdminMenuRoutes(router *gin.RouterGroup, handler *MenuHandler, middleware *middleware.JWTMiddleware) {
	admin := router.Group("/admin/menu")
	admin.Use(middleware.Authenticate())
	admin.Use(middleware.RequireRole("admin"))
	{
		// CRUD операции с меню
		admin.POST("", handler.CreateMenuItem)
		admin.PUT("/:id", handler.UpdateMenuItem)
		admin.DELETE("/:id", handler.DeleteMenuItem)

		// Управление доступностью
		admin.PATCH("/:id/availability", handler.UpdateAvailability)
		admin.PATCH("/:id/activation", handler.MenuItemActivation)

		// Управление категориями
		admin.GET("/categories", handler.GetCategories)
		admin.POST("/categories", handler.CreateCategory)
	}
}
