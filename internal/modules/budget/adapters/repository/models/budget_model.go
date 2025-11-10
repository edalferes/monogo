package models

import (
	"time"

	"gorm.io/gorm"
)

// BudgetModel represents the database model for Budget
type BudgetModel struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index:idx_user_budgets"`
	CategoryID  uint      `gorm:"not null;index:idx_category_budgets"`
	Name        string    `gorm:"not null;size:100"`
	Amount      float64   `gorm:"type:decimal(15,2);not null"`
	Spent       float64   `gorm:"type:decimal(15,2);default:0"`
	Period      string    `gorm:"not null;size:20"`
	StartDate   time.Time `gorm:"not null;index:idx_budget_period"`
	EndDate     time.Time `gorm:"not null;index:idx_budget_period"`
	AlertAt     *float64  `gorm:"type:decimal(5,2)"`
	IsActive    bool      `gorm:"default:true"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Category CategoryModel `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
}

func (BudgetModel) TableName() string {
	return "budget_budgets"
}
