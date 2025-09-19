package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type GetRoleByIDUseCase struct {
	RoleRepo repository.RoleRepository
}

func (u *GetRoleByIDUseCase) Execute(id uint) (*domain.Role, error) {
	return u.RoleRepo.FindByID(id)
}
