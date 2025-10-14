package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type CreatePermissionUseCase struct {
	PermissionRepo interfaces.Permission
}

func (u *CreatePermissionUseCase) Execute(name string) error {
	perm := &domain.Permission{Name: name}
	return u.PermissionRepo.Create(perm)
}
