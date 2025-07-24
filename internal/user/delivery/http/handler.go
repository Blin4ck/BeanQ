package http

import (
	"coffe/internal/middleware"
	"coffe/internal/user/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	middleware  *middleware.JWTMiddleware
}

func NewUserHandler(userUseCase *usecase.UserUseCase, middleware *middleware.JWTMiddleware) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		middleware:  middleware,
	}
}

// регистрация пользователя
func (h *UserHandler) Register(ctx *gin.Context) {
	var req usecase.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	if err := h.userUseCase.RegisterUser(ctx.Request.Context(), req); err != nil {
		switch err {
		case usecase.ErrUserAlreadyExists:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации пользователя"})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})

}

// вход пользователя
func (h *UserHandler) Login(ctx *gin.Context) {
	var req usecase.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	loginResponse, err := h.userUseCase.LoginUser(ctx.Request.Context(), req)

	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при входе пользователя"})
		}
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

// вход администратора
func (h *UserHandler) AdminLogin(ctx *gin.Context) {
	var req usecase.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	loginResponse, err := h.userUseCase.AdminLogin(ctx.Request.Context(), req)

	if err != nil {
		switch err {
		case usecase.ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при входе пользователя"})
		}
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)

}

// обновление токенов
func (h *UserHandler) RefreshToken(ctx *gin.Context) {

	var req struct {
		RefreshToken string    `json:"refresh_token" validate:"required"`
		UserID       uuid.UUID `json:"user_id" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.userUseCase.RefreshTokens(ctx.Request.Context(), req.UserID, req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении токенов"})
		return
	}

	response := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	ctx.JSON(http.StatusOK, response)

}

// получение профиля пользователя
func (h *UserHandler) GetProfile(ctx *gin.Context) {

	userID, exists := ctx.Get("user_id")

	if !exists {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})

	}

	id, ok := userID.(uuid.UUID)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ID пользователя"})
		return
	}

	userProfile, err := h.userUseCase.GetUserByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения профиля пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, userProfile)
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {

	userID, exists := ctx.Get("user_id")

	if !exists {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})

	}

	id, ok := userID.(uuid.UUID)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ID пользователя"})
		return
	}

	var req usecase.UpdateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	err := h.userUseCase.UpdateUser(ctx, id, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля пользователя"})
		return
	}

	updatedUser, err := h.userUseCase.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения обновленного профиля"})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)

}

// изменение пароля
func (h *UserHandler) ChangePassword(ctx *gin.Context) {

	var req usecase.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	err := h.userUseCase.ChangePassword(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка изменения пароля"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})

}

// выход пользователя
func (h *UserHandler) Logout(ctx *gin.Context) {

	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	id, ok := userID.(uuid.UUID)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ID пользователя"})
		return
	}

	err := h.userUseCase.Logout(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка выхода пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно вышел"})

}

// получение списка всех пользователей (админская функция)
func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.userUseCase.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения списка пользователей"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

// получение пользователя по ID из URL параметра (админская функция)
func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	// Получаем ID из URL параметра
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	userProfile, err := h.userUseCase.GetUserByID(ctx, id)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения профиля пользователя"})
		}
		return
	}

	ctx.JSON(http.StatusOK, userProfile)
}

// получение пользователей по роли (админская функция)
func (h *UserHandler) GetUsersByRole(ctx *gin.Context) {
	role := ctx.Param("role")
	if role == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Роль обязательна"})
		return
	}

	users, err := h.userUseCase.GetUsersByRole(ctx, role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения пользователей по роли"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users, "role": role})
}
