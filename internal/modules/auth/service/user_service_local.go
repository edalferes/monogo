package service

import (
	"context"

	"github.com/edalferes/monetics/internal/contracts"
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository"
)

// UserServiceLocal implements UserService interface using local repository
// This is used when the auth module runs in the same process
type UserServiceLocal struct {
	userRepo repository.UserRepository
}

// NewUserServiceLocal creates a new local user service
func NewUserServiceLocal(userRepo repository.UserRepository) contracts.UserService {
	return &UserServiceLocal{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves user information from local repository
func (s *UserServiceLocal) GetUserByID(ctx context.Context, userID uint) (*contracts.UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to DTO
	return &contracts.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		IsActive: true, // Assuming all users are active by default
	}, nil
}

// ValidateUserExists checks if a user exists in local repository
func (s *UserServiceLocal) ValidateUserExists(ctx context.Context, userID uint) (bool, error) {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetUserPermissions retrieves user permissions from local repository
func (s *UserServiceLocal) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	return permissions, nil
}
