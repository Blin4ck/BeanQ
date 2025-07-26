package repositories

import (
	"coffe/internal/common"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных
func (r *UserRepository) CreateUser(ctx context.Context, user *common.User) error {
	if user == nil {
		return errors.New("пользователь не может быть nil")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByEmail получает пользователя по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*common.User, error) {
	if email == "" {
		return nil, errors.New("email не может быть пустым")
	}
	var user common.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*common.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("ID не может быть пустым")
	}
	var user common.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

// GetRoleByName получает роль по имени
func (r *UserRepository) GetRoleByName(ctx context.Context, name string) (*common.Role, error) {
	if name == "" {
		return nil, errors.New("название роли не может быть пустым")
	}
	var role common.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("роль не найдена")
		}
		return nil, err
	}
	return &role, nil
}

// CheckUserExists проверяет существование пользователя по email
func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email не может быть пустым")
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&common.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteUser удаляет пользователя по ID
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("ID не может быть пустым")
	}
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&common.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}

// UpdateUserRole обновляет роль пользователя
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

// GetUsersByRole получает всех пользователей с определенной ролью
func (r *UserRepository) GetUsersByRole(ctx context.Context, roleName string) ([]*common.User, error) {
	if roleName == "" {
		return nil, errors.New("название роли не может быть пустым")
	}
	var users []*common.User
	if err := r.db.WithContext(ctx).Preload("Role").Joins("JOIN roles ON users.role_id = roles.id").Where("roles.name = ?", roleName).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser обновляет данные пользователя
func (r *UserRepository) UpdateUser(ctx context.Context, user *common.User) error {
	if user == nil {
		return errors.New("пользователь не может быть nil")
	}
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

// UpdatePassword обновляет пароль пользователя по id
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	if id == uuid.Nil {
		return errors.New("ID пользователя не может быть пустым")
	}
	return r.db.WithContext(ctx).
		Model(&common.User{}).
		Where("id = ?", id).
		Update("password", newPassword).Error
}

// GetAllUsers получает всех пользователей из базы данных
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*common.User, error) {
	var users []*common.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
