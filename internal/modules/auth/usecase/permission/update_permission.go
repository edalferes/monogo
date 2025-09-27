package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type UpdatePermissionInput struct {
	ID   uint
	Name string
}

type UpdatePermissionUseCase struct {
	PermissionReader interfaces.PermissionReader
	PermissionWriter interfaces.PermissionWriter
}

func (u *UpdatePermissionUseCase) Execute(input UpdatePermissionInput) error {
	perm, err := u.PermissionReader.FindByID(input.ID)
	if err != nil {
		return err
	}
	if input.Name != "" {
		perm.Name = input.Name
	}
	return u.PermissionWriter.Update(perm)
}
