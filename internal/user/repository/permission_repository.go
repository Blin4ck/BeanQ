package repository

import (
	"coffe/internal/user/entity"
	"context"

	"github.com/google/uuid"
)

type PermissionRepository interface {
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
	CreateMany(ctx context.Context, permissions []entity.Permission) error
	AssignToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
}
