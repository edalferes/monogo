package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type UpdateUserByAdminUseCase struct {
	UserRepo        interfaces.User
	PasswordService interfaces.PasswordService
}

// Execute updates a user's basic information (admin action)
func (u *UpdateUserByAdminUseCase) Execute(userID uint, username, password string) error {
	// Find the user
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Update username if provided
	if username != "" {
		user.Username = username
	}

	// Update password if provided
	if password != "" {
		hashedPassword, err := u.PasswordService.Hash(password)
		if err != nil {
			return errors.ErrInvalidData
		}
		user.Password = hashedPassword
	}

	// Update the user
	return u.UserRepo.Update(user)
}
