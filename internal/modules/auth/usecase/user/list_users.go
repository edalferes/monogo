package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type ListUsersUseCase struct {
	UserRepo repository.UserRepository
}

func (u *ListUsersUseCase) Execute() ([]domain.User, error) {
	return u.UserRepo.ListAll()
}
