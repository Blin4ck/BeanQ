package repository

import (
	"coffe/internal/common"
	"context"

	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository определяет методы для работы с пользователями и ролями.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *common.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*common.User, error) {
	var user common.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*common.User, error) {
	var user common.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *common.User) error {
	if user.ID == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	result := r.db.WithContext(ctx).Delete(&common.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	if id == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	result := r.db.WithContext(ctx).Model(&common.User{}).Where("id = ?", id).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

// Подтверждение по email реализовать бизнес логику
func (r *UserRepository) SendEmailVerification(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	// TODO: Implement email sending logic
	return nil
}

func (r *UserRepository) VerifyEmail(ctx context.Context, token string) error {
	// TODO: Implement email verification logic
	return nil
}

// Операции с ролями
func (r *UserRepository) GetRoleByName(ctx context.Context, name string) (*common.Role, error) {
	var role common.Role
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *UserRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, newRoleID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	result := r.db.WithContext(ctx).Model(&common.User{}).Where("id = ?", userID).Update("role_id", newRoleID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, roleName string) ([]*common.User, error) {
	var users []*common.User
	result := r.db.WithContext(ctx).Joins("Role").Where("roles.name = ?", roleName).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*common.User, error) {
	var users []*common.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Вспомогательные методы
func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&common.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&common.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
