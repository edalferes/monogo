package permission

import "github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"

type DeletePermissionUseCase struct {
	PermissionRepo interfaces.Permission
}

func (u *DeletePermissionUseCase) Execute(name string) error {
	return u.PermissionRepo.DeleteByName(name)
}
