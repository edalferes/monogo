package domain

import "time"

// AuditLog representa um registro de auditoria de ações sensíveis
// Campos: ID, UserID (nullable), Username, Action, Status, IP, Details, CreatedAt
// Exemplos de Actions: "login_success", "login_failed", "user_created", "password_changed"
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
