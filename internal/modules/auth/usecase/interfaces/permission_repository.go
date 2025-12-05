package interfaces

import "github.com/edalferes/monetics/internal/modules/auth/domain"

type PermissionRepository interface {
	FindByName(name string) (*domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	Create(permission *domain.Permission) error
	ListAll() ([]domain.Permission, error)
	DeleteByName(name string) error
	Update(permission *domain.Permission) error
}
