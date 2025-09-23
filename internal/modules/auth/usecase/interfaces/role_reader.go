package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// RoleReader representa operações de leitura de roles
type RoleReader interface {
	FindByName(name string) (*domain.Role, error)
	FindByID(id uint) (*domain.Role, error)
	ListAll() ([]domain.Role, error)
}
