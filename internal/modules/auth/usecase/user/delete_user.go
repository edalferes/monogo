package user

import "github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"

type DeleteUserUseCase struct {
	UserRepo interfaces.User
}

func (u *DeleteUserUseCase) Execute(id uint) error {
	return u.UserRepo.Delete(id)
}
