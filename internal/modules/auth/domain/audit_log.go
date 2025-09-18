package domain

import "time"

// AuditLog represents an audit record of sensitive actions
// Fields: ID, UserID (nullable), Username, Action, Status, IP, Details, CreatedAt
// Example Actions: "login_success", "login_failed", "user_created", "password_changed"
type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action" gorm:"not null"`
	Status    string    `json:"status"`
	IP        string    `json:"ip"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}
