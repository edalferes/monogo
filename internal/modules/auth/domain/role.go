package domain

// Role represents a role/function in the system
type Role struct {
	ID          uint
	Name        string
	Permissions []Permission
}
