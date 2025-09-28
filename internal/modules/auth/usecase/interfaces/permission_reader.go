package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// PermissionReader represents permission read operations
type PermissionReader interface {
	FindByName(name string) (*domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	ListAll() ([]domain.Permission, error)
}
