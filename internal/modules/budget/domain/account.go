package domain

import "time"

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeChecking AccountType = "checking"   // Conta corrente
	AccountTypeSavings  AccountType = "savings"    // Poupança
	AccountTypeCredit   AccountType = "credit"     // Cartão de crédito
	AccountTypeCash     AccountType = "cash"       // Dinheiro
	AccountTypeInvest   AccountType = "investment" // Investimentos
)

// Account represents a financial account (bank account, credit card, cash, etc.)
//
// An Account belongs to a User and contains the balance and transaction history.
// It can represent different types of accounts like checking, savings, credit cards.
//
// Business rules:
//   - Each account must belong to a user
//   - Initial balance can be positive or negative (debt)
//   - Account name must be unique per user
//   - Balance is updated automatically with transactions
//
// Example:
//
//	account := &Account{
//		UserID:      1,
//		Name:        "Banco Inter",
//		Type:        AccountTypeChecking,
//		Balance:     5000.00,
//		Currency:    "BRL",
//		Description: "Conta corrente principal",
//	}
type Account struct {
	ID          uint        `json:"id"`
	UserID      uint        `json:"user_id"`
	Name        string      `json:"name"`
	Type        AccountType `json:"type"`
	Balance     float64     `json:"balance"`
	Currency    string      `json:"currency"`
	Description string      `json:"description,omitempty"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
