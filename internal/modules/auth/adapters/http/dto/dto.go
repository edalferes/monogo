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

// ChangePasswordDTO represents the payload for changing password
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

// UpdateUserDTO represents the payload for updating user
type UpdateUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
