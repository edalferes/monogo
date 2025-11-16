package role

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type ListRolesUseCase struct {
	RoleRepo interfaces.Role
}

func (u *ListRolesUseCase) Execute() ([]domain.Role, error) {
	return u.RoleRepo.ListAll()
}
