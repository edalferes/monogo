package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type GetPermissionByIDUseCase struct {
	PermissionRepo repository.PermissionRepository
}

// Supondo que PermissionRepository tenha FindByID
func (u *GetPermissionByIDUseCase) Execute(id uint) (*domain.Permission, error) {
	return u.PermissionRepo.FindByID(id)
}
