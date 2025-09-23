package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// UserWriter representa operações de escrita de usuários
type UserWriter interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}
