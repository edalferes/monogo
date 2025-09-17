package user

import "github.com/edalferes/monogo/internal/modules/user/domain"

func Entities() []interface{} {
	return []interface{}{&domain.User{}}
}
