package repository

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type PermissionRepository interface {
	FindByName(name string) (*domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	Create(permission *domain.Permission) error
	ListAll() ([]domain.Permission, error)
	DeleteByName(name string) error
	Update(permission *domain.Permission) error
}
