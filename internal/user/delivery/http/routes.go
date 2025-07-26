package http

import (
	"coffe/internal/middleware"
	"coffe/internal/user/usecase"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes настраивает все маршруты для пользовательского модуля
func SetupUserRoutes(router *gin.RouterGroup, handler *UserHandler, jwtMiddleware *middleware.JWTMiddleware, permissionUC usecase.PermissionUsecase) {
	// Публичные маршруты аутентификации
	setupAuthRoutes(router, handler)

	// Защищенные пользовательские маршруты
	setupProtectedUserRoutes(router, handler, jwtMiddleware)

	// Админские маршруты
	setupAdminRoutes(router, handler, jwtMiddleware, permissionUC)
}

// setupAuthRoutes настраивает публичные маршруты аутентификации
func setupAuthRoutes(router *gin.RouterGroup, handler *UserHandler) {
	auth := router.Group("/auth")
	{
		// Регистрация и вход
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/admin/login", handler.AdminLogin)

		// Управление токенами
		auth.POST("/refresh", handler.RefreshToken)
	}
}

// setupProtectedUserRoutes настраивает защищенные пользовательские маршруты
func setupProtectedUserRoutes(router *gin.RouterGroup, handler *UserHandler, jwtMiddleware *middleware.JWTMiddleware) {
	protected := router.Group("/users")
	protected.Use(jwtMiddleware.Authenticate())
	{
		// Управление профилем
		protected.GET("/profile", handler.GetProfile)
		protected.PUT("/profile", handler.UpdateProfile)

		// Безопасность
		protected.POST("/change-password", handler.ChangePassword)
		protected.POST("/logout", handler.Logout)
	}
}

// setupAdminRoutes настраивает админские маршруты
// Эти маршруты требуют JWT токен + роль администратора + проверку разрешений
func setupAdminRoutes(router *gin.RouterGroup, handler *UserHandler, jwtMiddleware *middleware.JWTMiddleware, permissionUC usecase.PermissionUsecase) {
	admin := router.Group("/admin/users")
	admin.Use(jwtMiddleware.Authenticate())
	admin.Use(jwtMiddleware.RequireRole("admin"))
	{
		admin.GET("", middleware.PermissionMiddleware(permissionUC, "read_user"), handler.GetAllUsers)
		admin.GET("/:id", middleware.PermissionMiddleware(permissionUC, "read_user"), handler.GetUserByID)
		admin.GET("/role/:role", middleware.PermissionMiddleware(permissionUC, "read_user"), handler.GetUsersByRole)
	}
}

func GetRoutesList() map[string][]RouteInfo {
	return map[string][]RouteInfo{
		"auth": {
			{Method: "POST", Path: "/auth/register", Description: "Регистрация пользователя", Auth: false},
			{Method: "POST", Path: "/auth/login", Description: "Вход пользователя", Auth: false},
			{Method: "POST", Path: "/auth/admin/login", Description: "Вход администратора", Auth: false},
			{Method: "POST", Path: "/auth/refresh", Description: "Обновление токенов", Auth: false},
		},
		"users": {
			{Method: "GET", Path: "/users/profile", Description: "Получение профиля", Auth: true},
			{Method: "PUT", Path: "/users/profile", Description: "Обновление профиля", Auth: true},
			{Method: "POST", Path: "/users/change-password", Description: "Смена пароля", Auth: true},
			{Method: "POST", Path: "/users/logout", Description: "Выход из системы", Auth: true},
		},
		"admin": {
			{Method: "GET", Path: "/admin/users", Description: "Список всех пользователей", Auth: true, AdminOnly: true},
			{Method: "GET", Path: "/admin/users/:id", Description: "Пользователь по ID", Auth: true, AdminOnly: true},
			{Method: "GET", Path: "/admin/users/role/:role", Description: "Пользователи по роли", Auth: true, AdminOnly: true},
		},
	}
}

// RouteInfo содержит информацию о маршруте для документации
type RouteInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Auth        bool   `json:"auth_required"`
	AdminOnly   bool   `json:"admin_only,omitempty"`
}
