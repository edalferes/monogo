package contracts

import "context"

// UserService defines operations available to other modules
// This interface allows modules to depend on abstractions rather than concrete implementations,

type UserService interface {
	// GetUserByID retrieves basic user information
	GetUserByID(ctx context.Context, userID uint) (*UserInfo, error)

	// ValidateUserExists checks if a user exists in the system
	ValidateUserExists(ctx context.Context, userID uint) (bool, error)

	// GetUserPermissions returns all permissions for a given user
	GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
}

// UserInfo represents basic user information shared across modules
// This is a Data Transfer Object (DTO) to decouple domain entities
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
