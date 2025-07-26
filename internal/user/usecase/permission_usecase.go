package usecase

import (
	"coffe/internal/user/entity"
	"coffe/internal/user/repository"
	"context"

	"github.com/google/uuid"
)

type PermissionUsecase interface {
	GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
	InitPredefinedPermissions(ctx context.Context) error
}

type permissionUsecase struct {
	repo repository.PermissionRepository
}

func NewPermissionUsecase(repo repository.PermissionRepository) PermissionUsecase {
	return &permissionUsecase{repo: repo}
}

func (u *permissionUsecase) GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	return u.repo.GetByRoleID(ctx, roleID)
}

func (u *permissionUsecase) InitPredefinedPermissions(ctx context.Context) error {
	return u.repo.CreateMany(ctx, entity.PredefinedPermissions)
}
