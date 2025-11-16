package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type AssignRoleUseCase struct {
	UserRepo interfaces.User
	RoleRepo interfaces.Role
}

// Execute assigns a role to a user by role name
func (u *AssignRoleUseCase) Execute(userID uint, roleName string) error {
	// Find the user
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Find the role
	role, err := u.RoleRepo.FindByName(roleName)
	if err != nil {
		return err // Role not found or other error
	}

	// Check if user already has this role
	for _, existingRole := range user.Roles {
		if existingRole.ID == role.ID {
			return nil // User already has this role, no error
		}
	}

	// Add the role to user's roles
	user.Roles = append(user.Roles, *role)

	// Update the user
	return u.UserRepo.Update(user)
}
