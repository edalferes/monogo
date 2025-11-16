package role

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type GetRoleByIDUseCase struct {
	RoleRepo interfaces.Role
}

func (u *GetRoleByIDUseCase) Execute(id uint) (*domain.Role, error) {
	return u.RoleRepo.FindByID(id)
}
