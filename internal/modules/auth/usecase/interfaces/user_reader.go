// Package interfaces defines the contracts for repository operations used by use cases.
//
// These interfaces implement the Interface Segregation Principle (ISP) by separating
// read and write operations. This allows use cases to depend only on the specific
// operations they need, improving testability and reducing coupling.
package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

// UserReader represents read-only user repository operations.
//
// This interface defines all query operations needed for user-related use cases.
// It follows the Interface Segregation Principle by only exposing read operations,
// ensuring that use cases that only need to read user data cannot accidentally
// modify the data.
//
// Implementation notes:
//   - All methods should preload associated data (Roles, Permissions) when needed
//   - Return nil user and error when not found
//   - Use appropriate database indexes for performance
//   - Consider pagination for ListAll in production systems
//
// Example usage in use cases:
//
//	type LoginUseCase struct {
//		UserRepo UserReader // Only needs read access
//	}
//
//	func (uc *LoginUseCase) Execute(username string) error {
//		user, err := uc.UserRepo.FindByUsername(username)
//		if err != nil || user == nil {
//			return ErrUserNotFound
//		}
//		// Process user...
//	}
type UserReader interface {
	// FindByUsername retrieves a user by their unique username.
	// Returns the user with preloaded roles, or nil if not found.
	// This method is primarily used for authentication.
	//
	// Parameters:
	//   - username: the unique username to search for
	//
	// Returns:
	//   - *domain.User: the found user with roles preloaded, nil if not found
	//   - error: database error if query fails
	FindByUsername(username string) (*domain.User, error)

	// FindByID retrieves a user by their unique ID.
	// Returns the user with preloaded roles, or nil if not found.
	// This method is used for profile operations and authorization checks.
	//
	// Parameters:
	//   - id: the unique user ID to search for
	//
	// Returns:
	//   - *domain.User: the found user with roles preloaded, nil if not found
	//   - error: database error if query fails
	FindByID(id uint) (*domain.User, error)

	// ListAll retrieves all users in the system.
	// Returns users with preloaded roles for administration purposes.
	// Consider implementing pagination for large datasets.
	//
	// Returns:
	//   - []domain.User: slice of all users with roles preloaded
	//   - error: database error if query fails
	//
	// Note: This method may return large datasets. Consider adding
	// pagination parameters in production systems:
	//   ListAllPaginated(offset, limit int) ([]domain.User, error)
	ListAll() ([]domain.User, error)
}
