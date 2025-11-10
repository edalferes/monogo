package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountModel represents the database model for Account
type AccountModel struct {
	ID          uint    `gorm:"primaryKey"`
	UserID      uint    `gorm:"not null;index:idx_user_accounts"`
	Name        string  `gorm:"not null;size:100"`
	Type        string  `gorm:"not null;size:20"`
	Balance     float64 `gorm:"type:decimal(15,2);default:0"`
	Currency    string  `gorm:"size:3;default:'BRL'"`
	Description string  `gorm:"type:text"`
	IsActive    bool    `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (AccountModel) TableName() string {
	return "budget_accounts"
}
