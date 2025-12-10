package domain

import "time"

// CategoryType represents if category is for income or expense
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"  // Receita
	CategoryTypeExpense CategoryType = "expense" // Despesa
)

// Category represents a transaction category for budget organization
//
// Categories help users organize their transactions and track spending
// by different areas (food, transport, salary, etc.).
//
// Business rules:
//   - Each category must belong to a user
//   - Category name must be unique per user and type
//   - System can have default categories for new users
//   - Categories can have an associated icon/color for UI
//
// Example:
//
//	category := &Category{
//		UserID:      1,
//		Name:        "Alimentação",
//		Type:        CategoryTypeExpense,
//		Icon:        "🍽️",
//		Color:       "#FF5733",
//		Description: "Gastos com supermercado e restaurants",
//	}
type Category struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UserID      uint         `json:"user_id" gorm:"not null;index:idx_user_categories;constraint:OnDelete:CASCADE"`
	Name        string       `json:"name" gorm:"not null;size:100"`
	Type        CategoryType `json:"type" gorm:"not null;size:20"`
	Icon        string       `json:"icon,omitempty" gorm:"size:50"`
	Color       string       `json:"color,omitempty" gorm:"size:20"`
	Description string       `json:"description,omitempty" gorm:"type:text"`
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (Category) TableName() string {
	return "budget_categories"
}
