package dto

// RoleDTO represents the payload for creating a role
type RoleDTO struct {
	Name string `json:"name" validate:"required"`
}

// PermissionDTO represents the payload for creating a permission
type PermissionDTO struct {
	Name string `json:"name" validate:"required"`
}

// LoginDTO represents the payload for user login
type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterDTO represents the payload for user registration
type RegisterDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
