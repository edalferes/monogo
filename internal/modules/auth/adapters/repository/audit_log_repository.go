package repository

import (
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type AuditLogRepository struct {
	DB *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{
		DB: db,
	}
}

var _ interfaces.AuditLogRepository = (*AuditLogRepository)(nil)

func (r *AuditLogRepository) Create(log *domain.AuditLog) error {
	return r.DB.Create(log).Error
}

func (r *AuditLogRepository) ListAll() ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	if err := r.DB.Order("created_at desc").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
