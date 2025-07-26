package repositories

import (
	"coffe/internal/user/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	var permissions []entity.Permission

	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error

	return permissions, err
}

func (r *PermissionRepository) CreateMany(ctx context.Context, permissions []entity.Permission) error {
	if len(permissions) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&permissions).Error
}

func (r *PermissionRepository) AssignToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	var rolePermissions []entity.RolePermission
	for _, pid := range permissionIDs {
		rolePermissions = append(rolePermissions, entity.RolePermission{
			RoleID:       roleID,
			PermissionID: pid,
		})
	}

	return r.db.WithContext(ctx).Create(&rolePermissions).Error
}
