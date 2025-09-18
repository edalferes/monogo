package gorm

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	gormpkg "gorm.io/gorm"
)

type AuditLogRepositoryGorm struct {
	DB *gormpkg.DB
}

func NewAuditLogRepositoryGorm(db *gormpkg.DB) *AuditLogRepositoryGorm {
	return &AuditLogRepositoryGorm{DB: db}
}

var _ repository.AuditLogRepository = (*AuditLogRepositoryGorm)(nil)

func (r *AuditLogRepositoryGorm) Create(log *domain.AuditLog) error {
	return r.DB.Create(log).Error
}

func (r *AuditLogRepositoryGorm) ListAll() ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	if err := r.DB.Order("created_at desc").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
