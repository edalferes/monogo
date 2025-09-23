package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// PermissionReader representa operações de leitura de permissions
type PermissionReader interface {
	FindByName(name string) (*domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	ListAll() ([]domain.Permission, error)
}
