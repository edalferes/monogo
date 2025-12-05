package interfaces

type AuditService interface {
	Log(userID *uint, username, action, status, ip, details string) error
}
