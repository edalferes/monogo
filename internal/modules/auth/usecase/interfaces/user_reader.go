package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// UserReader represents user read operations
type UserReader interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	ListAll() ([]domain.User, error)
}
