package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type RegisterUseCase struct {
	UserReader      interfaces.UserReader
	UserWriter      interfaces.UserWriter
	RoleReader      interfaces.RoleReader
	PasswordService service.PasswordService
}

func (u *RegisterUseCase) Execute(username, password string) error {
	if user, _ := u.UserReader.FindByUsername(username); user != nil {
		return errors.ErrUserAlreadyExists
	}
	hash, err := u.PasswordService.Hash(password)
	if err != nil {
		return err
	}
	role, err := u.RoleReader.FindByName("user")
	if err != nil {
		return err
	}
	user := &domain.User{
		Username: username,
		Password: hash,
		Roles:    []domain.Role{*role},
	}
	return u.UserWriter.Create(user)
}
