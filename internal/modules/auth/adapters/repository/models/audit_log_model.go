package models

import "time"

// AuditLogModel representa a estrutura de dados para GORM
type AuditLogModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `gorm:"index"`
	Username  string    `gorm:"not null"`
	Action    string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	IP        string    `gorm:"not null"`
	Details   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}
