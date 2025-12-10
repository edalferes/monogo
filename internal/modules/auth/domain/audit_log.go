package domain

import "time"

// AuditLog represents an audit record of sensitive security actions for compliance and monitoring.
//
// Audit logs provide a complete trail of security-relevant events in the system.
// They are essential for compliance, security monitoring, and forensic analysis.
// All authentication attempts, authorization changes, and sensitive operations are logged.
//
// Action naming conventions:
//   - "login_success" - successful authentication
//   - "login_failed" - failed authentication attempt
//   - "user_created" - new user account created
//   - "user_updated" - user account modified
//   - "user_deleted" - user account deleted
//   - "password_changed" - password was updated
//   - "role_assigned" - role was assigned to user
//   - "role_removed" - role was removed from user
//   - "permission_granted" - permission was granted
//   - "permission_revoked" - permission was revoked
//
// Status values:
//   - "success" - operation completed successfully
//   - "failed" - operation failed
//   - "blocked" - operation was blocked by security policy
//
// Persistence considerations:
//   - ID should be mapped to audit_logs.id (primary key)
//   - UserID should be mapped to audit_logs.user_id (nullable foreign key)
//   - Username should be mapped to audit_logs.username (for when UserID is null)
//   - Action should be mapped to audit_logs.action (indexed for queries)
//   - Status should be mapped to audit_logs.status (indexed for queries)
//   - IP should be mapped to audit_logs.ip (client IP address)
//   - Details should be mapped to audit_logs.details (JSON string with additional context)
//   - CreatedAt should be mapped to audit_logs.created_at (timestamp with index)
//
// Example:
//
//	auditLog := &AuditLog{
//		UserID:   &user.ID,
//		Username: user.Username,
//		Action:   "login_success",
//		Status:   "success",
//		IP:       "192.168.1.100",
//		Details:  `{"user_agent": "Mozilla/5.0...", "session_id": "abc123"}`,
//	}
type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id,omitempty" gorm:"index"`
	Username  string    `json:"username" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null;index"`
	Status    string    `json:"status" gorm:"not null;index"`
	IP        string    `json:"ip" gorm:"size:45"`
	Details   string    `json:"details,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
