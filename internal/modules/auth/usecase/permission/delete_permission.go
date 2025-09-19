package permission

import "github.com/edalferes/monogo/internal/modules/auth/repository"

type DeletePermissionUseCase struct {
	PermissionRepo repository.PermissionRepository
}

func (u *DeletePermissionUseCase) Execute(name string) error {
	return u.PermissionRepo.DeleteByName(name)
}
