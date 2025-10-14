package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type ListPermissionsUseCase struct {
	PermissionRepo interfaces.Permission
}

func (u *ListPermissionsUseCase) Execute() ([]domain.Permission, error) {
	return u.PermissionRepo.ListAll()
}
