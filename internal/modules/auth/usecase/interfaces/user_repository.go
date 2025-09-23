package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
	ListAll() ([]domain.User, error)
}
