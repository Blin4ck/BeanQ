package repository

import (
	"coffe/internal/common"
	"context"

	"github.com/google/uuid"
)

// UserRepository определяет методы для работы с пользователями и ролями.
type UserRepository interface {
	// Основные операции с пользователями
	CreateUser(ctx context.Context, user *common.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*common.User, error)
	GetUserByEmail(ctx context.Context, email string) (*common.User, error)
	UpdateUser(ctx context.Context, user *common.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error

	// Подтверждение по email реализовать бизнес логику
	SendEmailVerification(ctx context.Context, userID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error

	// Операции с ролями
	GetRoleByName(ctx context.Context, name string) (*common.Role, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, newRoleID int) error
	GetUsersByRole(ctx context.Context, roleName string) ([]*common.User, error)
	GetAllUsers(ctx context.Context) ([]*common.User, error)

	// Вспомогательные методы
	CheckUserExists(ctx context.Context, email string) (bool, error)
	Count(ctx context.Context) (int64, error)
}
