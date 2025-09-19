package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type ListPermissionsUseCase struct {
	PermissionRepo repository.PermissionRepository
}

func (u *ListPermissionsUseCase) Execute() ([]domain.Permission, error) {
	return u.PermissionRepo.ListAll()
}
