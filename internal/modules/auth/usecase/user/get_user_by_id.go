package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type GetUserByIDUseCase struct {
	UserRepo interfaces.UserReader
}

func (u *GetUserByIDUseCase) Execute(id uint) (*domain.User, error) {
	return u.UserRepo.FindByID(id)
}
