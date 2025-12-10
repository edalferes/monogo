package domain

// Role represents a role/function in the system for role-based access control (RBAC).
//
// Roles group related permissions together and are assigned to users.
// This allows for flexible authorization where users can have multiple roles,
// and each role can have multiple permissions.
//
// Common roles might include:
//   - "admin": full system access
//   - "user": basic user operations
//   - "moderator": content moderation rights
//   - "viewer": read-only access
//
// Relationships:
//   - Role has many Permissions (many-to-many through role_permissions table)
//   - Role belongs to many Users (many-to-many through user_roles table)
//
// Persistence considerations:
//   - ID should be mapped to roles.id (primary key)
//   - Name should be mapped to roles.name (unique constraint)
//   - Permissions relationship should be loaded when needed
//
// Example:
//
//	adminRole := &Role{
//		Name: "admin",
//		Permissions: []Permission{
//			{Name: "users:read"},
//			{Name: "users:write"},
//			{Name: "roles:read"},
//			{Name: "roles:write"},
//		},
//	}
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE"`
}

func (Role) TableName() string {
	return "roles"
}
