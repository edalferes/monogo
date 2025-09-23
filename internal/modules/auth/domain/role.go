package domain

// Role representa um papel/função no sistema
type Role struct {
	ID          uint
	Name        string
	Permissions []Permission
}
