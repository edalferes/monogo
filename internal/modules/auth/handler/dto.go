package handler

// RoleDTO represents the payload for creating a role
type RoleDTO struct {
	Name string `json:"name" validate:"required"`
}

// PermissionDTO represents the payload for creating a permission
type PermissionDTO struct {
	Name string `json:"name" validate:"required"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
