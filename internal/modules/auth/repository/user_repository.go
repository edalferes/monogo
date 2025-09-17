package repository

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
	ListAll() ([]domain.User, error)
}
