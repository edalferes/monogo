package gorm

import (
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	gormpkg "gorm.io/gorm"
)

type AuditLogRepositoryGorm struct {
	DB     *gormpkg.DB
	mapper mappers.AuditLogMapper
}

func NewAuditLogRepositoryGorm(db *gormpkg.DB) *AuditLogRepositoryGorm {
	return &AuditLogRepositoryGorm{
		DB:     db,
		mapper: mappers.AuditLogMapper{},
	}
}

var _ interfaces.AuditLogRepository = (*AuditLogRepositoryGorm)(nil)

func (r *AuditLogRepositoryGorm) Create(log *domain.AuditLog) error {
	logModel := r.mapper.ToModel(*log)
	return r.DB.Create(&logModel).Error
}

func (r *AuditLogRepositoryGorm) ListAll() ([]domain.AuditLog, error) {
	var logModels []models.AuditLogModel
	if err := r.DB.Order("created_at desc").Find(&logModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(logModels), nil
}
