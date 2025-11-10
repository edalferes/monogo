package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TransactionModel represents the database model for Transaction
type TransactionModel struct {
	ID                   uint      `gorm:"primaryKey"`
	UserID               uint      `gorm:"not null;index:idx_user_transactions"`
	AccountID            uint      `gorm:"not null;index:idx_account_transactions"`
	CategoryID           uint      `gorm:"not null;index:idx_category_transactions"`
	Type                 string    `gorm:"not null;size:20"`
	Amount               float64   `gorm:"type:decimal(15,2);not null"`
	Description          string    `gorm:"type:text"`
	Date                 time.Time `gorm:"not null;index:idx_transaction_date"`
	Status               string    `gorm:"not null;size:20;default:'completed'"`
	DestinationAccountID *uint     `gorm:"index:idx_destination_account"`
	TransferFee          *float64  `gorm:"type:decimal(15,2)"`
	IsRecurring          bool      `gorm:"default:false"`
	RecurrenceRule       string    `gorm:"size:50"`
	RecurrenceEnd        *time.Time
	ParentID             *uint
	Tags                 pq.StringArray `gorm:"type:text[]"`
	Attachments          pq.StringArray `gorm:"type:text[]"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

func (TransactionModel) TableName() string {
	return "budget_transactions"
}
