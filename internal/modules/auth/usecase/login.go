package usecase

import (
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
)

type LoginUseCase struct {
	UserRepo        repository.UserRepository
	PasswordService service.PasswordService
	JWTService      service.JWTService
}

func (u *LoginUseCase) Execute(username, password string) (string, error) {
	user, err := u.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.ErrInvalidCredentials
	}
	if err := u.PasswordService.Compare(user.Password, password); err != nil {
		return "", errors.ErrInvalidCredentials
	}
	// Extrai roles do usuário
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}
	token, err := u.JWTService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return "", errors.ErrInvalidData
	}
	return token, nil
}
