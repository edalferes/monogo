package models

import (
	"time"

	"gorm.io/gorm"
)

// CategoryModel represents the database model for Category
type CategoryModel struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index:idx_user_categories"`
	Name        string `gorm:"not null;size:100"`
	Type        string `gorm:"not null;size:20"`
	Icon        string `gorm:"size:10"`
	Color       string `gorm:"size:20"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (CategoryModel) TableName() string {
	return "budget_categories"
}
