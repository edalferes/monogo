// Package errors defines authentication and authorization error types.
//
// This package provides sentinel errors used throughout the auth module.
// Using sentinel errors allows for proper error handling and type checking
// in use cases and handlers.
//
// Security considerations:
//   - Error messages are intentionally generic to prevent information disclosure
//   - No specific details about why authentication failed (prevents username enumeration)
//   - Consistent error messages across different failure scenarios
//
// Usage example:
//
//	token, err := loginUseCase.Execute(username, password, ip)
//	if err == errors.ErrInvalidCredentials {
//		return c.JSON(http.StatusUnauthorized, "Authentication failed")
//	}
package errors

import "errors"

var (
	// ErrInvalidCredentials indicates that the provided username/password combination is incorrect.
	// This error is returned for both "user not found" and "wrong password" scenarios
	// to prevent username enumeration attacks.
	//
	// HTTP Status: 401 Unauthorized
	// User Message: "Invalid username or password"
	// Audit Action: "login_failed"
	ErrInvalidCredentials = errors.New("invalid username or password")

	// ErrMissingCredentials indicates that required authentication fields are missing.
	// This error is returned when username or password fields are empty in the request.
	//
	// HTTP Status: 400 Bad Request
	// User Message: "Username and password are required"
	// Audit Action: Not logged (no credentials to audit)
	ErrMissingCredentials = errors.New("username and password are required")

	// ErrInvalidData indicates that the provided data failed validation.
	// This is a generic error for various data validation failures such as:
	//   - Invalid request format
	//   - Failed data binding
	//   - Token generation failures
	//   - Invalid field values
	//
	// HTTP Status: 400 Bad Request
	// User Message: "Invalid data provided"
	// Audit Action: Context-dependent
	ErrInvalidData = errors.New("invalid data")

	// ErrUserAlreadyExists indicates that a user with the same username already exists.
	// This error is returned during user registration when attempting to create
	// a user with a username that is already taken.
	//
	// HTTP Status: 409 Conflict
	// User Message: "User already exists"
	// Audit Action: "user_creation_failed"
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound indicates that the requested user does not exist.
	// This error is used internally but often converted to ErrInvalidCredentials
	// for security reasons.
	//
	// HTTP Status: 404 Not Found (internal) / 401 Unauthorized (external)
	// User Message: "User not found" (internal) / "Invalid credentials" (external)
	// Audit Action: Context-dependent
	ErrUserNotFound = errors.New("user not found")

	// ErrInsufficientPermissions indicates that the user lacks required permissions.
	// This error is returned when a user attempts to access a resource or perform
	// an action that requires permissions they don't have.
	//
	// HTTP Status: 403 Forbidden
	// User Message: "Insufficient permissions"
	// Audit Action: "access_denied"
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// ErrInvalidToken indicates that the provided JWT token is invalid or expired.
	// This includes scenarios such as:
	//   - Malformed tokens
	//   - Expired tokens
	//   - Tokens with invalid signatures
	//   - Tokens for non-existent users
	//
	// HTTP Status: 401 Unauthorized
	// User Message: "Invalid or expired token"
	// Audit Action: "invalid_token_used"
	ErrInvalidToken = errors.New("invalid or expired token")
)
