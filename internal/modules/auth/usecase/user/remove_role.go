package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type RemoveRoleUseCase struct {
	UserRepo interfaces.User
	RoleRepo interfaces.Role
}

// Execute removes a role from a user by role name
func (u *RemoveRoleUseCase) Execute(userID uint, roleName string) error {
	// Find the user
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Find the role to ensure it exists
	role, err := u.RoleRepo.FindByName(roleName)
	if err != nil {
		return err // Role not found or other error
	}

	// Remove the role from user's roles
	var updatedRoles []domain.Role
	roleRemoved := false
	for _, existingRole := range user.Roles {
		if existingRole.ID != role.ID {
			updatedRoles = append(updatedRoles, existingRole)
		} else {
			roleRemoved = true
		}
	}

	// If role wasn't found in user's roles, return without error
	if !roleRemoved {
		return nil
	}

	// Update user with new roles list
	user.Roles = updatedRoles

	// Update the user
	return u.UserRepo.Update(user)
}
