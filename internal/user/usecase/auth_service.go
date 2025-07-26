package usecase

import (
	"coffe/internal/auth"
	"coffe/internal/common"
	"coffe/internal/user/entity"
	"time"

	"context"
	"errors"

	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository определяет методы для работы с пользователями и ролями.
type UserRepository interface {
	CreateUser(ctx context.Context, user *common.User) error
	GetUserByEmail(ctx context.Context, email string) (*common.User, error)
	GetRoleByName(ctx context.Context, name string) (*common.Role, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*common.User, error)
}

// TokenRepository определяет методы для хранения и получения refresh-токенов.
type TokenRepository interface {
	SetToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetToken(ctx context.Context, userID string) (string, error)
	DeleteToken(ctx context.Context, userID string) error
}

// AuthService реализует логику аутентификации пользователей.
type AuthService struct {
	userRepo   UserRepository
	jwtService *auth.JWTService
	tokenRepo  TokenRepository
}

// NewAuthService создает новый экземпляр AuthService.
func NewAuthService(userRepo UserRepository, jwtService *auth.JWTService, tokenRepo TokenRepository) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		tokenRepo:  tokenRepo,
	}
}

// Login выполняет аутентификацию пользователя и возвращает access-токен, refresh-токен и userID.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, uuid.UUID, error) {
	if email == "" {
		return "", "", uuid.Nil, errors.New("email не может быть пустым")
	}
	if password == "" {
		return "", "", uuid.Nil, errors.New("пароль не может быть пустым")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", uuid.Nil, errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", "", uuid.Nil, errors.New("неверный пароль")
		}
		return "", "", uuid.Nil, err
	}

	accessToken, err := s.jwtService.GenerateJWTToken(user.ID.String(), user.Role.Name, user.Role.ID.String())
	if err != nil {
		return "", "", uuid.Nil, errors.New("ошибка генерации токена")
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", uuid.Nil, errors.New("ошибка генерации refresh токена")
	}

	return accessToken, refreshToken, user.ID, nil
}

// AdminLogin выполняет аутентификацию администратора и возвращает оба токена.
func (s *AuthService) AdminLogin(ctx context.Context, email, password string) (string, string, error) {
	if email == "" {
		return "", "", errors.New("email не может быть пустым")
	}
	if password == "" {
		return "", "", errors.New("пароль не может быть пустым")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("пользователь не найден")
	}

	// Проверяем права администратора
	if user.Role == nil || (user.Role.Name != "admin" && user.Role.Name != "manager") {
		return "", "", errors.New("нет прав для входа")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", "", errors.New("неверный пароль")
		}
		return "", "", err
	}

	accessToken, err := s.jwtService.GenerateJWTToken(user.ID.String(), user.Role.Name, user.Role.ID.String())
	if err != nil {
		return "", "", errors.New("ошибка генерации токена")
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", errors.New("ошибка генерации refresh токена")
	}

	return accessToken, refreshToken, nil
}

// Register регистрирует нового пользователя и возвращает токены.
func (s *AuthService) Register(ctx context.Context, user *common.User) (string, string, uuid.UUID, error) {
	// Валидация пользователя
	if err := entity.ValidateUser(user); err != nil {
		return "", "", uuid.Nil, err
	}

	// Проверяем, что пользователь с таким email не существует
	exists, err := s.userRepo.CheckUserExists(ctx, user.Email)
	if err != nil {
		return "", "", uuid.Nil, err
	}
	if exists {
		return "", "", uuid.Nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return "", "", uuid.Nil, errors.New("ошибка хеширования пароля")
	}
	user.Password = hashedPassword

	// Устанавливаем роль по умолчанию, если не указана
	if user.Role == nil || user.Role.Name == "" {
		role, err := s.userRepo.GetRoleByName(ctx, "client")
		if err != nil {
			return "", "", uuid.Nil, errors.New("роль client не найдена")
		}
		user.Role = role
	}

	// Создаем пользователя в базе данных
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return "", "", uuid.Nil, errors.New("ошибка создания пользователя")
	}

	// Генерируем токены
	accessToken, err := s.jwtService.GenerateJWTToken(user.ID.String(), user.Role.Name, user.Role.ID.String())
	if err != nil {
		return "", "", uuid.Nil, errors.New("ошибка генерации токена")
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", uuid.Nil, errors.New("ошибка генерации refresh токена")
	}

	return accessToken, refreshToken, user.ID, nil
}

// generateSecureToken генерирует криптостойкий случайный токен длиной n байт.
func generateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateRefreshToken создает и сохраняет refresh-токен для пользователя.
func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := generateSecureToken(32)
	if err != nil {
		return "", err
	}
	ttl := 7 * 24 * time.Hour
	if err := s.tokenRepo.SetToken(ctx, userID.String(), token, ttl); err != nil {
		return "", err
	}
	return token, nil
}

// RefreshTokens проверяет refresh-токен и выдает новые токены.
func (s *AuthService) RefreshTokens(ctx context.Context, userID uuid.UUID, refreshToken string) (string, string, error) {
	storedToken, err := s.tokenRepo.GetToken(ctx, userID.String())
	if err != nil || storedToken != refreshToken {
		return "", "", errors.New("ошибочны востановочный токен")
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.jwtService.GenerateJWTToken(userID.String(), user.Role.Name, user.Role.ID.String())
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := s.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}

// Logout удаляет refresh-токен пользователя.
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRepo.DeleteToken(ctx, userID.String())
}

//придумать как реализовать то чтобы у всех ролей были свои точки доступа, учитывая, что у админа есть все права и при добавлении новой роли можно было бы выбирать что она может, а что нет.
