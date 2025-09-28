// Package user contains use cases related to user operations such as authentication,
// registration, and profile management.
package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

// LoginWithAuditUseCase handles user authentication with comprehensive audit logging.
//
// This use case implements secure user login functionality that:
//   - Validates user credentials against the database
//   - Uses bcrypt for secure password comparison
//   - Generates JWT tokens for successful authentication
//   - Logs all authentication attempts for security monitoring
//   - Returns appropriate errors for different failure scenarios
//
// Security features:
//   - Password comparison using constant-time bcrypt
//   - Audit logging for both successful and failed attempts
//   - Generic error messages to prevent username enumeration
//   - IP address tracking for suspicious activity detection
//
// Dependencies:
//   - UserRepo: for user data retrieval (read-only operations)
//   - PasswordService: for secure password comparison
//   - JWTService: for token generation
//   - AuditService: for security event logging
//
// Example usage:
//
//	loginUC := &LoginWithAuditUseCase{
//		UserRepo:        userRepo,
//		PasswordService: passwordSvc,
//		JWTService:      jwtSvc,
//		AuditService:    auditSvc,
//	}
//
//	token, err := loginUC.Execute("admin", "password123", "192.168.1.100")
//	if err != nil {
//		// Handle authentication failure
//		return err
//	}
//	// Use token for subsequent requests
type LoginWithAuditUseCase struct {
	UserRepo        interfaces.UserReader   // Repository for user data access
	PasswordService service.PasswordService // Service for password operations
	JWTService      service.JWTService      // Service for JWT token management
	AuditService    service.AuditService    // Service for audit logging
}

// Execute performs user authentication with comprehensive audit logging.
//
// This method handles the complete authentication flow:
//  1. Looks up user by username
//  2. Compares provided password with stored hash
//  3. Generates JWT token with user roles
//  4. Logs the authentication attempt (success or failure)
//
// Parameters:
//   - username: the username provided by the client
//   - password: the plain text password provided by the client
//   - ip: the client's IP address for audit logging
//
// Returns:
//   - string: JWT token if authentication succeeds
//   - error: authentication error (see errors package for types)
//
// Possible errors:
//   - ErrInvalidCredentials: username not found or password mismatch
//   - ErrInvalidData: token generation failed
//
// Security considerations:
//   - Both "user not found" and "wrong password" return the same error to prevent username enumeration
//   - All attempts are logged with appropriate detail level
//   - IP address is captured for suspicious activity detection
//   - Failed attempts include reason in audit details for administrator review
//
// Example:
//
//	token, err := usecase.Execute("alice", "secret123", "192.168.1.100")
//	if err == errors.ErrInvalidCredentials {
//		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
//	}
//	// Return token to client
func (u *LoginWithAuditUseCase) Execute(username, password, ip string) (string, error) {
	user, err := u.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		u.AuditService.Log(nil, username, "login_failed", "fail", ip, "user not found or error")
		return "", errors.ErrInvalidCredentials
	}
	if err := u.PasswordService.Compare(user.Password, password); err != nil {
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "wrong password")
		return "", errors.ErrInvalidCredentials
	}
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}
	token, err := u.JWTService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "token error")
		return "", errors.ErrInvalidData
	}
	u.AuditService.Log(&user.ID, username, "login_success", "ok", ip, "")
	return token, nil
}
