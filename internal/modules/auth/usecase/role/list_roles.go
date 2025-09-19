package role

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type ListRolesUseCase struct {
	RoleRepo repository.RoleRepository
}

func (u *ListRolesUseCase) Execute() ([]domain.Role, error) {
	return u.RoleRepo.ListAll()
}
