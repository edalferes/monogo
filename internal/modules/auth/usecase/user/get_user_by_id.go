package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type GetUserByIDUseCase struct {
	UserRepo repository.UserRepository
}

func (u *GetUserByIDUseCase) Execute(id uint) (*domain.User, error) {
	return u.UserRepo.FindByID(id)
}
