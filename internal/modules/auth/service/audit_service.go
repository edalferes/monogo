package service

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type AuditService interface {
	Log(userID *uint, username, action, status, ip, details string) error
}

type auditService struct {
	repo repository.AuditLogRepository
}

func NewAuditService(repo repository.AuditLogRepository) AuditService {
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
