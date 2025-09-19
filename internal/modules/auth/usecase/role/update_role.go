package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type UpdateRoleInput struct {
	ID            uint
	Name          string
	PermissionIDs []uint
}

type UpdateRoleUseCase struct {
	RoleRepo repository.RoleRepository
}

func (u *UpdateRoleUseCase) Execute(input UpdateRoleInput) error {
	role, err := u.RoleRepo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if input.Name != "" {
		role.Name = input.Name
	}
	if input.PermissionIDs != nil {
		var perms []domain.Permission
		for _, pid := range input.PermissionIDs {
			perms = append(perms, domain.Permission{ID: pid})
		}
		role.Permissions = perms
	}
	return u.RoleRepo.Update(role)
}
