package role

import "github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"

type DeleteRoleUseCase struct {
	RoleRepo interfaces.Role
}

func (u *DeleteRoleUseCase) Execute(name string) error {
	return u.RoleRepo.DeleteByName(name)
}
