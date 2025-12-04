package audit

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type auditService struct {
	repo interfaces.AuditLogRepository
}

func NewAuditService(repo interfaces.AuditLogRepository) interfaces.AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(userID *uint, username, action, status, ip, details string) error {
	log := &domain.AuditLog{
		UserID:   userID,
		Username: username,
		Action:   action,
		Status:   status,
		IP:       ip,
		Details:  details,
	}
	return s.repo.Create(log)
}
