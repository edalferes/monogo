package permission

import (
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type CheckPermissionUseCase struct {
	UserRepo       interfaces.User
	PermissionRepo interfaces.Permission
}

// Execute checks if a user has a specific permission through their roles
func (u *CheckPermissionUseCase) Execute(userID uint, permissionName string) (bool, error) {
	// Find the user with their roles and permissions
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	// Check if user has the permission through any of their roles
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			if permission.Name == permissionName {
				return true, nil
			}
		}
	}

	return false, nil
}
