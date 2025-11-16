package permission

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type GetPermissionByIDUseCase struct {
	PermissionRepo interfaces.Permission
}

// Supondo que PermissionRepository tenha FindByID
func (u *GetPermissionByIDUseCase) Execute(id uint) (*domain.Permission, error) {
	return u.PermissionRepo.FindByID(id)
}
