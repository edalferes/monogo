package domain

import "time"

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"   // Receita
	TransactionTypeExpense  TransactionType = "expense"  // Despesa
	TransactionTypeTransfer TransactionType = "transfer" // Transfer between accounts
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"   // Pendente
	TransactionStatusCompleted TransactionStatus = "completed" // Completed
	TransactionStatusCancelled TransactionStatus = "cancelled" // Cancelada
)

// Transaction represents a financial transaction (income, expense, or transfer)
//
// Transactions are the core of the budget system, tracking money flow
// between accounts and categories.
//
// Business rules:
//   - Each transaction must belong to a user
//   - Must have an account (source for expenses/income)
//   - Transfers must have source and destination accounts
//   - Amount must be positive
//   - Date can be future for planned transactions
//   - Completed transactions update account balance
//
// Example:
//
//	transaction := &Transaction{
//		UserID:      1,
//		AccountID:   1,
//		CategoryID:  5,
//		Type:        TransactionTypeExpense,
//		Amount:      150.00,
//		Description: "Supermercado",
//		Date:        time.Now(),
//		Status:      TransactionStatusCompleted,
//	}
type Transaction struct {
	ID          uint              `json:"id"`
	UserID      uint              `json:"user_id"`
	AccountID   uint              `json:"account_id"` // Conta origem
	CategoryID  uint              `json:"category_id"`
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	Description string            `json:"description"`
	Date        time.Time         `json:"date"`
	Month       string            `json:"month"` // Format: "2025-01" for easy grouping and filtering
	Status      TransactionStatus `json:"status"`
	// Transfer specific fields
	DestinationAccountID *uint    `json:"destination_account_id,omitempty"` // For transfers
	TransferFee          *float64 `json:"transfer_fee,omitempty"`           // Transfer fee
	// Recurrence
	IsRecurring    bool       `json:"is_recurring"`
	RecurrenceRule string     `json:"recurrence_rule,omitempty"` // "monthly", "weekly", etc.
	RecurrenceEnd  *time.Time `json:"recurrence_end,omitempty"`
	ParentID       *uint      `json:"parent_id,omitempty"` // For recurring transactions
	// Metadata
	Tags        []string  `json:"tags,omitempty"`
	Attachments []string  `json:"attachments,omitempty"` // URLs de comprovantes
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
