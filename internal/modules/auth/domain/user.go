package domain

// User represents a user in the system
type User struct {
	ID       uint
	Username string
	Password string
	Roles    []Role
}
