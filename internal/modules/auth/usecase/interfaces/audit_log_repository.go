package interfaces

import "github.com/edalferes/monogo/internal/modules/auth/domain"

type AuditLogRepository interface {
	Create(log *domain.AuditLog) error
	ListAll() ([]domain.AuditLog, error)
}
