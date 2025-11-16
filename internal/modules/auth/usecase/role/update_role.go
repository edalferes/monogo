package role

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type UpdateRoleInput struct {
	ID            uint
	Name          string
	PermissionIDs []uint
}

type UpdateRoleUseCase struct {
	Role interfaces.Role
}

func (u *UpdateRoleUseCase) Execute(input UpdateRoleInput) error {
	role, err := u.Role.FindByID(input.ID)
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
	return u.Role.Update(role)
}
