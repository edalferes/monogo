package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/service"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type RegisterUseCase struct {
	User            interfaces.User
	Role            interfaces.Role
	PasswordService service.PasswordService
}

func (u *RegisterUseCase) Execute(username, password string) error {
	if user, _ := u.User.FindByUsername(username); user != nil {
		return errors.ErrUserAlreadyExists
	}
	hash, err := u.PasswordService.Hash(password)
	if err != nil {
		return err
	}
	role, err := u.Role.FindByName("user")
	if err != nil {
		return err
	}
	user := &domain.User{
		Username: username,
		Password: hash,
		Roles:    []domain.Role{*role},
	}
	return u.User.Create(user)
}
