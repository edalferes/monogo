package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type ListRolesUseCase struct {
	RoleRepo interfaces.RoleReader
}

func (u *ListRolesUseCase) Execute() ([]domain.Role, error) {
	return u.RoleRepo.ListAll()
}
