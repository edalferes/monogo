package usecase

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
)

type RegisterUseCase struct {
	UserRepo        repository.UserRepository
	RoleRepo        repository.RoleRepository
	PasswordService service.PasswordService
}

func (u *RegisterUseCase) Execute(username, password string) error {
	if user, _ := u.UserRepo.FindByUsername(username); user != nil {
		return errors.ErrUserAlreadyExists
	}
	hash, err := u.PasswordService.Hash(password)
	if err != nil {
		return err
	}
	role, err := u.RoleRepo.FindByName("user")
	if err != nil {
		return err
	}
	user := &domain.User{
		Username: username,
		Password: hash,
		Roles:    []domain.Role{*role},
	}
	return u.UserRepo.Create(user)
}
