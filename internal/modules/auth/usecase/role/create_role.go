package role

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type CreateRoleUseCase struct {
	RoleRepo interfaces.Role
}

func (u *CreateRoleUseCase) Execute(name string, permissionIDs []uint) error {
	perms := make([]domain.Permission, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		perms = append(perms, domain.Permission{ID: pid})
	}
	role := &domain.Role{Name: name, Permissions: perms}
	return u.RoleRepo.Create(role)
}
