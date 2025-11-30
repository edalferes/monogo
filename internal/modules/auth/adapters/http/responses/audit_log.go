package responses

import "time"

type AuditLogResponse struct {
	ID        uint      `json:"id"`
	UserID    *uint     `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	IP        string    `json:"ip"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
