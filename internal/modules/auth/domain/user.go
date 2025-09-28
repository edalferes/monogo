// Package domain contains the core business entities for the auth module.
package domain

// User represents a system user with authentication and authorization data.
//
// A User contains the essential information needed for authentication (username, password)
// and authorization (roles). The password field stores the bcrypt hash, never plain text.
//
// Relationships:
//   - User has many Roles (many-to-many through user_roles table)
//   - Through Roles, User has access to Permissions
//
// Persistence considerations:
//   - ID should be mapped to users.id (primary key)
//   - Username should be mapped to users.username (unique constraint)
//   - Password should be mapped to users.password_hash
//   - Roles relationship should be loaded when needed
//
// Example:
//
//	user := &User{
//		Username: "admin",
//		Password: "$2a$10$...", // bcrypt hash
//		Roles: []Role{
//			{Name: "admin"},
//			{Name: "user"},
//		},
//	}
type User struct {
	ID       uint   `json:"id"`       // Unique identifier
	Username string `json:"username"` // Unique username for login
	Password string `json:"-"`        // Bcrypt hash (never exposed in JSON)
	Roles    []Role `json:"roles"`    // Associated roles for authorization
}
