package mappers

import (
	"github.com/edalferes/monogo/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monogo/internal/modules/auth/domain"
)

// AuditLogMapper converte entre domain.AuditLog e models.AuditLogModel
type AuditLogMapper struct{}

// ToModel converte domain.AuditLog para models.AuditLogModel
func (m AuditLogMapper) ToModel(auditLog domain.AuditLog) models.AuditLogModel {
	return models.AuditLogModel{
		ID:        auditLog.ID,
		UserID:    auditLog.UserID,
		Username:  auditLog.Username,
		Action:    auditLog.Action,
		Status:    auditLog.Status,
		IP:        auditLog.IP,
		Details:   auditLog.Details,
		CreatedAt: auditLog.CreatedAt,
	}
}

// ToDomain converte models.AuditLogModel para domain.AuditLog
func (m AuditLogMapper) ToDomain(auditLogModel models.AuditLogModel) domain.AuditLog {
	return domain.AuditLog{
		ID:        auditLogModel.ID,
		UserID:    auditLogModel.UserID,
		Username:  auditLogModel.Username,
		Action:    auditLogModel.Action,
		Status:    auditLogModel.Status,
		IP:        auditLogModel.IP,
		Details:   auditLogModel.Details,
		CreatedAt: auditLogModel.CreatedAt,
	}
}

// ToModelSlice converte []domain.AuditLog para []models.AuditLogModel
func (m AuditLogMapper) ToModelSlice(auditLogs []domain.AuditLog) []models.AuditLogModel {
	auditLogModels := make([]models.AuditLogModel, len(auditLogs))
	for i, auditLog := range auditLogs {
		auditLogModels[i] = m.ToModel(auditLog)
	}
	return auditLogModels
}

// ToDomainSlice converte []models.AuditLogModel para []domain.AuditLog
func (m AuditLogMapper) ToDomainSlice(auditLogModels []models.AuditLogModel) []domain.AuditLog {
	auditLogs := make([]domain.AuditLog, len(auditLogModels))
	for i, auditLogModel := range auditLogModels {
		auditLogs[i] = m.ToDomain(auditLogModel)
	}
	return auditLogs
}
