package usecase

import (
	"coffe/internal/common"
	"coffe/internal/user/repository"

	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrUserAlreadyExists возвращается, если пользователь уже существует.
var ErrUserAlreadyExists = errors.New("пользователь уже существует")

// ErrInvalidCredentials возвращается, если переданы неверные учетные данные.
var ErrInvalidCredentials = errors.New("неверные учетные данные")

// ErrUserNotFound возвращается, если пользователь не найден.
var ErrUserNotFound = errors.New("пользователь не найден")

// RegisterRequest содержит данные для регистрации пользователя.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

// ChangePasswordRequest содержит данные для смены пароля пользователя.
type ChangePasswordRequest struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

// UserResponse содержит публичные данные пользователя для ответа API.
type UserResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Email   string    `json:"email"`
	Role    string    `json:"role"`
}

// LoginResponse содержит данные для ответа при успешной аутентификации.
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserUseCase реализует бизнес-логику для работы с пользователями.
type UserUseCase struct {
	userRepo    repository.UserRepository
	authService *AuthService
}

// NewUserUseCase создает новый usecase для пользователя.
func NewUserUseCase(userRepo repository.UserRepository, authService *AuthService) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		authService: authService,
	}
}

// RegisterUser регистрирует нового пользователя.
func (uc *UserUseCase) RegisterUser(ctx context.Context, req RegisterRequest) error {
	if err := uc.validateRegisterRequest(req); err != nil {
		return err
	}
	user := &common.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: req.Password,
	}
	_, err := uc.authService.Register(ctx, user)
	return err
}

// LoginUser аутентифицирует пользователя и возвращает токен.
func (uc *UserUseCase) LoginUser(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := uc.validateLoginRequest(req); err != nil {
		return nil, err
	}
	token, err := uc.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &LoginResponse{
		User:  uc.toUserResponseFromCommon(user),
		Token: token,
	}, nil
}

// AdminLogin аутентифицирует администратора и возвращает токен.
func (uc *UserUseCase) AdminLogin(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := uc.validateLoginRequest(req); err != nil {
		return nil, err
	}
	token, err := uc.authService.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	user, err := uc.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &LoginResponse{
		User:  uc.toUserResponseFromCommon(user),
		Token: token,
	}, nil
}

// GetUserByID возвращает пользователя по его ID.
func (uc *UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	response := uc.toUserResponseFromCommon(user)
	return &response, nil
}

// UpdateUser обновляет данные пользователя.
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) error {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	user.Name = req.Name
	user.Surname = req.Surname
	user.Email = req.Email
	return uc.userRepo.UpdateUser(ctx, user)
}

// GetUsersByRole возвращает пользователей по роли.
func (uc *UserUseCase) GetUsersByRole(ctx context.Context, roleName string) ([]UserResponse, error) {
	users, err := uc.userRepo.GetUsersByRole(ctx, roleName)
	if err != nil {
		return nil, err
	}
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = uc.toUserResponseFromCommon(user)
	}
	return responses, nil
}

// ChangePassword меняет пароль пользователя с проверкой старого пароля и хешированием нового.
func (uc *UserUseCase) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	user, err := uc.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return ErrUserNotFound
	}
	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("старый пароль неверен")
	}
	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}
	// Обновляем пароль в базе
	if err := uc.userRepo.UpdatePassword(ctx, req.UserID, string(hashedPassword)); err != nil {
		return errors.New("ошибка обновления пароля")
	}
	return nil
}

// validateRegisterRequest валидирует данные для регистрации пользователя.
func (uc *UserUseCase) validateRegisterRequest(req RegisterRequest) error {
	if req.Name == "" {
		return errors.New("имя обязательно")
	}
	if req.Surname == "" {
		return errors.New("фамилия обязательна")
	}
	if req.Email == "" {
		return errors.New("email обязателен")
	}
	if len(req.Password) < 6 {
		return errors.New("пароль должен содержать минимум 6 символов")
	}
	return nil
}

// validateLoginRequest валидирует данные для входа пользователя.
func (uc *UserUseCase) validateLoginRequest(req LoginRequest) error {
	if req.Email == "" {
		return errors.New("email обязателен")
	}
	if req.Password == "" {
		return errors.New("пароль обязателен")
	}
	return nil
}

// toUserResponseFromCommon преобразует common.User в UserResponse для API.
func (uc *UserUseCase) toUserResponseFromCommon(user *common.User) UserResponse {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	return UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Role:    roleName,
	}
}
