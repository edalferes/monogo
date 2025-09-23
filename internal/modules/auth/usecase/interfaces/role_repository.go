package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type RoleRepository interface {
	FindByName(name string) (*domain.Role, error)
	FindByID(id uint) (*domain.Role, error)
	Create(role *domain.Role) error
	ListAll() ([]domain.Role, error)
	DeleteByName(name string) error
	Update(role *domain.Role) error
}
