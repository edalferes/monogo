package permission

import "github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"

type DeletePermissionUseCase struct {
	PermissionRepo interfaces.Permission
}

func (u *DeletePermissionUseCase) Execute(name string) error {
	return u.PermissionRepo.DeleteByName(name)
}
