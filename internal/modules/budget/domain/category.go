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
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	Name        string       `json:"name"`
	Type        CategoryType `json:"type"`
	Icon        string       `json:"icon,omitempty"`
	Color       string       `json:"color,omitempty"`
	Description string       `json:"description,omitempty"`
	IsActive    bool         `json:"is_active"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
