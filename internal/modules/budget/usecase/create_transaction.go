package usecase

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// CreateTransactionUseCase handles transaction creation
type CreateTransactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
}

// NewCreateTransactionUseCase creates a new use case instance
func NewCreateTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// Execute creates a new transaction
func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInput) (domain.Transaction, error) {
	// Validate input
	if input.UserID == 0 {
		return domain.Transaction{}, errors.ErrInvalidUserID
	}
	if input.AccountID == 0 {
		return domain.Transaction{}, errors.ErrInvalidAccountID
	}
	if input.CategoryID == 0 {
		return domain.Transaction{}, errors.ErrInvalidCategoryID
	}
	if input.Amount <= 0 {
		return domain.Transaction{}, errors.ErrInvalidAmount
	}
	if input.Description == "" {
		return domain.Transaction{}, errors.ErrTransactionDescriptionRequired
	}
	if !isValidTransactionType(input.Type) {
		return domain.Transaction{}, errors.ErrInvalidTransactionType
	}

	// Verify account exists and belongs to user
	account, err := uc.accountRepo.GetByID(ctx, input.AccountID)
	if err != nil {
		return domain.Transaction{}, errors.ErrAccountNotFound
	}
	if account.UserID != input.UserID {
		return domain.Transaction{}, errors.ErrUnauthorizedAccess
	}

	// Verify category exists and belongs to user
	category, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		return domain.Transaction{}, errors.ErrCategoryNotFound
	}
	if category.UserID != input.UserID {
		return domain.Transaction{}, errors.ErrUnauthorizedAccess
	}

	// Validate transfer if applicable
	if input.Type == domain.TransactionTypeTransfer {
		if input.DestinationAccountID == nil {
			return domain.Transaction{}, errors.ErrTransferRequiresDestination
		}
		destAccount, err := uc.accountRepo.GetByID(ctx, *input.DestinationAccountID)
		if err != nil {
			return domain.Transaction{}, errors.ErrAccountNotFound
		}
		if destAccount.UserID != input.UserID {
			return domain.Transaction{}, errors.ErrUnauthorizedAccess
		}
	}

	// Create transaction domain entity
	transaction := domain.Transaction{
		UserID:               input.UserID,
		AccountID:            input.AccountID,
		CategoryID:           input.CategoryID,
		Type:                 input.Type,
		Amount:               input.Amount,
		Description:          input.Description,
		Date:                 input.Date,
		Status:               domain.TransactionStatusCompleted,
		Tags:                 input.Tags,
		Attachments:          input.Attachments,
		IsRecurring:          input.IsRecurring,
		RecurrenceRule:       input.RecurrenceRule,
		RecurrenceEnd:        input.RecurrenceEnd,
		DestinationAccountID: input.DestinationAccountID,
		TransferFee:          input.TransferFee,
	}

	// Save to repository
	return uc.transactionRepo.Create(ctx, transaction)
}

// CreateTransactionInput represents the input for creating a transaction
type CreateTransactionInput struct {
	UserID               uint                   `json:"user_id"`
	AccountID            uint                   `json:"account_id"`
	CategoryID           uint                   `json:"category_id"`
	Type                 domain.TransactionType `json:"type"`
	Amount               float64                `json:"amount"`
	Description          string                 `json:"description"`
	Date                 time.Time              `json:"date"`
	Tags                 []string               `json:"tags,omitempty"`
	Attachments          []string               `json:"attachments,omitempty"`
	IsRecurring          bool                   `json:"is_recurring"`
	RecurrenceRule       string                 `json:"recurrence_rule,omitempty"`
	RecurrenceEnd        *time.Time             `json:"recurrence_end,omitempty"`
	DestinationAccountID *uint                  `json:"destination_account_id,omitempty"`
	TransferFee          *float64               `json:"transfer_fee,omitempty"`
}

func isValidTransactionType(transactionType domain.TransactionType) bool {
	switch transactionType {
	case domain.TransactionTypeIncome, domain.TransactionTypeExpense, domain.TransactionTypeTransfer:
		return true
	default:
		return false
	}
}
