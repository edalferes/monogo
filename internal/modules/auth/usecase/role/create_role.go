package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type CreateRoleUseCase struct {
	RoleRepo repository.RoleRepository
}

func (u *CreateRoleUseCase) Execute(name string, permissionIDs []uint) error {
	var perms []domain.Permission
	for _, pid := range permissionIDs {
		perms = append(perms, domain.Permission{ID: pid})
	}
	role := &domain.Role{Name: name, Permissions: perms}
	return u.RoleRepo.Create(role)
}
