package service

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/usecase/interfaces"
)

type AuditService interface {
	Log(userID *uint, username, action, status, ip, details string) error
}

type auditService struct {
	repo interfaces.AuditLogRepository
}

func NewAuditService(repo interfaces.AuditLogRepository) AuditService {
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
