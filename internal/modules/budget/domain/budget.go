package domain

import "time"

// BudgetPeriod represents the budget period type
type BudgetPeriod string

const (
	BudgetPeriodMonthly   BudgetPeriod = "monthly"   // Mensal
	BudgetPeriodQuarterly BudgetPeriod = "quarterly" // Trimestral
	BudgetPeriodYearly    BudgetPeriod = "yearly"    // Anual
	BudgetPeriodCustom    BudgetPeriod = "custom"    // Personalizado
)

// Budget represents a budget plan for a category
//
// Budgets help users set spending limits for categories over time periods.
// They enable tracking of planned vs actual spending.
//
// Business rules:
//   - Each budget must belong to a user and category
//   - Amount must be positive
//   - Period dates must be valid (start before end)
//   - Cannot have overlapping budgets for same category
//   - Spent amount is calculated from transactions
//
// Example:
//
//	budget := &Budget{
//		UserID:     1,
//		CategoryID: 5,
//		Name:       "Orçamento Alimentação Janeiro",
//		Amount:     2000.00,
//		Period:     BudgetPeriodMonthly,
//		StartDate:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
//		EndDate:    time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
//	}
type Budget struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	CategoryID  uint         `json:"category_id"`
	Name        string       `json:"name"`
	Amount      float64      `json:"amount"` // Valor planejado
	Spent       float64      `json:"spent"`  // Valor gasto (calculado)
	Period      BudgetPeriod `json:"period"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     time.Time    `json:"end_date"`
	AlertAt     *float64     `json:"alert_at,omitempty"` // Percentual para alerta (ex: 80.0)
	IsActive    bool         `json:"is_active"`
	Description string       `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// RemainingAmount returns how much budget is left
func (b *Budget) RemainingAmount() float64 {
	return b.Amount - b.Spent
}

// PercentageUsed returns the percentage of budget used
func (b *Budget) PercentageUsed() float64 {
	if b.Amount == 0 {
		return 0
	}
	return (b.Spent / b.Amount) * 100
}

// IsOverBudget checks if spending exceeded the budget
func (b *Budget) IsOverBudget() bool {
	return b.Spent > b.Amount
}

// ShouldAlert checks if alert threshold is reached
func (b *Budget) ShouldAlert() bool {
	if b.AlertAt == nil {
		return false
	}
	return b.PercentageUsed() >= *b.AlertAt
}
