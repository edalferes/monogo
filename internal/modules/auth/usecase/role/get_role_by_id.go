package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type GetRoleByIDUseCase struct {
	RoleRepo interfaces.RoleReader
}

func (u *GetRoleByIDUseCase) Execute(id uint) (*domain.Role, error) {
	return u.RoleRepo.FindByID(id)
}
