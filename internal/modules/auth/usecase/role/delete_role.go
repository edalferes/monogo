package role

import "github.com/edalferes/monogo/internal/modules/auth/repository"

type DeleteRoleUseCase struct {
	RoleRepo repository.RoleRepository
}

func (u *DeleteRoleUseCase) Execute(name string) error {
	return u.RoleRepo.DeleteByName(name)
}
