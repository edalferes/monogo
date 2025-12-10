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
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id" gorm:"not null;index:idx_user_transactions;constraint:OnDelete:CASCADE"`
	AccountID   uint              `json:"account_id" gorm:"not null;index:idx_account_transactions"`
	CategoryID  uint              `json:"category_id" gorm:"not null;index:idx_category_transactions"`
	Type        TransactionType   `json:"type" gorm:"not null;size:20"`
	Amount      float64           `json:"amount" gorm:"type:decimal(15,2);not null"`
	Description string            `json:"description" gorm:"type:text"`
	Date        time.Time         `json:"date" gorm:"not null;index:idx_transaction_date"`
	Month       string            `json:"month" gorm:"size:7;index:idx_transaction_month"`
	Status      TransactionStatus `json:"status" gorm:"not null;size:20;default:'completed'"`
	// Transfer specific fields
	DestinationAccountID *uint    `json:"destination_account_id,omitempty" gorm:"index:idx_destination_account"`
	TransferFee          *float64 `json:"transfer_fee,omitempty" gorm:"type:decimal(15,2)"`
	// Recurrence
	IsRecurring    bool       `json:"is_recurring" gorm:"default:false"`
	RecurrenceRule string     `json:"recurrence_rule,omitempty" gorm:"size:50"`
	RecurrenceEnd  *time.Time `json:"recurrence_end,omitempty"`
	ParentID       *uint      `json:"parent_id,omitempty" gorm:"index:idx_parent_transaction"`
	// Metadata
	Tags        []string  `json:"tags,omitempty" gorm:"type:text;serializer:json"`
	Attachments []string  `json:"attachments,omitempty" gorm:"type:text;serializer:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Relations (populated by repository with Preload)
	Account  *Account  `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Transaction) TableName() string {
	return "budget_transactions"
}
