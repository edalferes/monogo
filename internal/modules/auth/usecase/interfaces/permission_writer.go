package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// PermissionWriter represents permission write operations
type PermissionWriter interface {
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	DeleteByName(name string) error
}
