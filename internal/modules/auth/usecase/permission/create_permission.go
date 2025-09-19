package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type CreatePermissionUseCase struct {
	PermissionRepo repository.PermissionRepository
}

func (u *CreatePermissionUseCase) Execute(name string) error {
	perm := &domain.Permission{Name: name}
	return u.PermissionRepo.Create(perm)
}
