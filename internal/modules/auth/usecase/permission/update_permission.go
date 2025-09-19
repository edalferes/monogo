package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type UpdatePermissionInput struct {
	ID   uint
	Name string
}

type UpdatePermissionUseCase struct {
	PermissionRepo repository.PermissionRepository
}

func (u *UpdatePermissionUseCase) Execute(input UpdatePermissionInput) error {
	perm, err := u.PermissionRepo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if input.Name != "" {
		perm.Name = input.Name
	}
	return u.PermissionRepo.Update(perm)
}
