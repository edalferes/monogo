package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type UpdatePermissionInput struct {
	ID   uint
	Name string
}

type UpdatePermissionUseCase struct {
	Permission interfaces.Permission
}

func (u *UpdatePermissionUseCase) Execute(input UpdatePermissionInput) error {
	perm, err := u.Permission.FindByID(input.ID)
	if err != nil {
		return err
	}
	if input.Name != "" {
		perm.Name = input.Name
	}
	return u.Permission.Update(perm)
}
