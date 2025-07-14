package usecase

import (
	"coffe/internal/auth"
	"coffe/internal/common"
	"coffe/internal/user/entity"

	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository определяет методы для работы с пользователями и ролями.
type UserRepository interface {
	CreateUser(ctx context.Context, user *common.User) error
	GetUserByEmail(ctx context.Context, email string) (*common.User, error)
	GetRoleByName(ctx context.Context, name string) (*common.Role, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
}

// AuthService реализует логику аутентификации пользователей.
type AuthService struct {
	userRepo   UserRepository
	jwtService *auth.JWTService
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(userRepo UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login выполняет аутентификацию пользователя и возвращает JWT-токен.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email не может быть пустым")
	}
	if password == "" {
		return "", errors.New("пароль не может быть пустым")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", errors.New("неверный пароль")
		}
		return "", err
	}

	token, err := s.jwtService.GenerateJWTToken(user.ID.String())
	if err != nil {
		return "", errors.New("ошибка генерации токена")
	}

	return token, nil
}

// AdminLogin выполняет аутентификацию администратора
func (s *AuthService) AdminLogin(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email не может быть пустым")
	}
	if password == "" {
		return "", errors.New("пароль не может быть пустым")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("пользователь не найден")
	}

	// Проверяем права администратора
	if user.Role == nil || (user.Role.Name != "owner" && user.Role.Name != "manager") {
		return "", errors.New("нет прав для входа")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", errors.New("неверный пароль")
		}
		return "", err
	}

	token, err := s.jwtService.GenerateJWTToken(user.ID.String())
	if err != nil {
		return "", errors.New("ошибка генерации токена")
	}

	return token, nil
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, user *common.User) (uuid.UUID, error) {
	// Валидация пользователя
	if err := entity.ValidateUser(user); err != nil {
		return uuid.Nil, err
	}

	// Проверяем, что пользователь с таким email не существует
	exists, err := s.userRepo.CheckUserExists(ctx, user.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return uuid.Nil, errors.New("ошибка хеширования пароля")
	}
	user.Password = hashedPassword

	// Устанавливаем роль по умолчанию, если не указана
	if user.Role == nil || user.Role.Name == "" {
		role, err := s.userRepo.GetRoleByName(ctx, "client")
		if err != nil {
			return uuid.Nil, errors.New("роль client не найдена")
		}
		user.Role = role
	}

	// Создаем пользователя в базе данных
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return uuid.Nil, errors.New("ошибка создания пользователя")
	}

	return user.ID, nil
}
