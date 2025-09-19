package user

import "github.com/edalferes/monogo/internal/modules/auth/repository"

type DeleteUserUseCase struct {
	UserRepo repository.UserRepository
}

func (u *DeleteUserUseCase) Execute(id uint) error {
	return u.UserRepo.Delete(id)
}
