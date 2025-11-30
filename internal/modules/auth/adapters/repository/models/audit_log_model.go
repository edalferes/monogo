package models

import "time"

// AuditLogModel representa a estrutura de dados para GORM
type AuditLogModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `gorm:"index;constraint:OnDelete:SET NULL"`
	Username  string    `gorm:"not null"`
	Action    string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	IP        string    `gorm:"not null"`
	Details   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	User *UserModel `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}
