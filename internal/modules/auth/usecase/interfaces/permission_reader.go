package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type Permission interface {
	FindByName(name string) (*domain.Permission, error)
	FindByID(id uint) (*domain.Permission, error)
	ListAll() ([]domain.Permission, error)
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	DeleteByName(name string) error
}
