package domain

import "time"

// AuditLog represents an audit record of sensitive actions
// Fields: ID, UserID (nullable), Username, Action, Status, IP, Details, CreatedAt
// Examples of Actions: "login_success", "login_failed", "user_created", "password_changed"
type AuditLog struct {
	ID        uint
	UserID    *uint
	Username  string
	Action    string
	Status    string
	IP        string
	Details   string
	CreatedAt time.Time
}
