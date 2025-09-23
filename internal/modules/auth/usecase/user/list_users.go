package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type ListUsersUseCase struct {
	UserRepo interfaces.UserReader
}

func (u *ListUsersUseCase) Execute() ([]domain.User, error) {
	return u.UserRepo.ListAll()
}
