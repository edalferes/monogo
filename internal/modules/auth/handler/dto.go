package handler

// RoleDTO representa o payload para criação de uma role
type RoleDTO struct {
	Name string `json:"name" validate:"required"`
}

// PermissionDTO representa o payload para criação de uma permission
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
