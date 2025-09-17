package auth

import "github.com/edalferes/monogo/internal/modules/auth/domain"

func Entities() []interface{} {
	return []interface{}{&domain.User{}, &domain.Role{}, &domain.Permission{}}
}
