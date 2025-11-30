package models

import (
	"time"

	"gorm.io/gorm"
)

// CategoryModel represents the database model for Category
type CategoryModel struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index:idx_user_categories;constraint:OnDelete:CASCADE"`
	Name        string `gorm:"not null;size:100"`
	Type        string `gorm:"not null;size:20"`
	Icon        string `gorm:"size:10"`
	Color       string `gorm:"size:20"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Transactions []TransactionModel `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
	Budgets      []BudgetModel      `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

func (CategoryModel) TableName() string {
	return "budget_categories"
}
