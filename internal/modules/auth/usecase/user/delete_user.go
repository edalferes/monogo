package user

import "github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"

type DeleteUserUseCase struct {
	UserRepo interfaces.UserWriter
}

func (u *DeleteUserUseCase) Execute(id uint) error {
	return u.UserRepo.Delete(id)
}
