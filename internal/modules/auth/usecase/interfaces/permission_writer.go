package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// PermissionWriter representa operações de escrita de permissions
type PermissionWriter interface {
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	DeleteByName(name string) error
}
