package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type ChangePasswordUseCase struct {
	UserRepo        interfaces.User
	PasswordService service.PasswordService
}

// Execute changes the user's password after verifying the current password
func (u *ChangePasswordUseCase) Execute(userID uint, currentPassword, newPassword string) error {
	// Find the user
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Verify current password
	if err := u.PasswordService.Compare(user.Password, currentPassword); err != nil {
		return errors.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := u.PasswordService.Hash(newPassword)
	if err != nil {
		return errors.ErrInvalidData
	}

	// Update user password
	user.Password = hashedPassword
	return u.UserRepo.Update(user)
}
