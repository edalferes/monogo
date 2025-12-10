package auth

import "github.com/edalferes/monetics/internal/modules/auth/domain"

func Entities() []interface{} {
	return []interface{}{
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.AuditLog{},
	}
}
