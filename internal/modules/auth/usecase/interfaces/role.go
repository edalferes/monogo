package interfaces

import "github.com/edalferes/monetics/internal/modules/auth/domain"

type Role interface {
	FindByName(name string) (*domain.Role, error)
	FindByID(id uint) (*domain.Role, error)
	ListAll() ([]domain.Role, error)
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	DeleteByName(name string) error
}
