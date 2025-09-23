package domain

// User representa um usuário no sistema
type User struct {
	ID       uint
	Username string
	Password string
	Roles    []Role
}
