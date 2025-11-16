package repository

import "github.com/edalferes/monetics/internal/modules/auth/domain"

type AuditLogRepository interface {
	Create(log *domain.AuditLog) error
	ListAll() ([]domain.AuditLog, error)
}
