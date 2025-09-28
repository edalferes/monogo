package domain

// Permission represents a granular permission in the system for fine-grained access control.
//
// Permissions define specific actions that can be performed within the system.
// They follow a resource:action naming convention for clarity and consistency.
//
// Permission naming patterns:
//   - "users:read" - can view users
//   - "users:write" - can create/update users
//   - "users:delete" - can delete users
//   - "roles:read" - can view roles
//   - "roles:write" - can create/update roles
//   - "admin:access" - can access admin panel
//
// Relationships:
//   - Permission belongs to many Roles (many-to-many through role_permissions table)
//   - Through Roles, Permission is accessible to Users
//
// Persistence considerations:
//   - ID should be mapped to permissions.id (primary key)
//   - Name should be mapped to permissions.name (unique constraint)
//
// Example:
//
//	permissions := []Permission{
//		{Name: "users:read"},
//		{Name: "users:write"},
//		{Name: "users:delete"},
//		{Name: "roles:read"},
//		{Name: "roles:write"},
//	}
type Permission struct {
	ID   uint   `json:"id"`   // Unique identifier
	Name string `json:"name"` // Unique permission name (resource:action format)
}
