package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// RoleWriter representa operações de escrita de roles
type RoleWriter interface {
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	DeleteByName(name string) error
}
