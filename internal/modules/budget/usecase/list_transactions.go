package usecase

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// ListTransactionsUseCase handles listing user transactions
type ListTransactionsUseCase struct {
	transactionRepo repository.TransactionRepository
}

// NewListTransactionsUseCase creates a new use case instance
func NewListTransactionsUseCase(transactionRepo repository.TransactionRepository) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{
		transactionRepo: transactionRepo,
	}
}

// Execute lists transactions for a user with optional filters
func (uc *ListTransactionsUseCase) Execute(ctx context.Context, input ListTransactionsInput) ([]domain.Transaction, error) {
	// Filter by date range if provided
	if input.StartDate != nil && input.EndDate != nil {
		return uc.transactionRepo.GetByDateRange(ctx, input.UserID, *input.StartDate, *input.EndDate)
	}

	// Filter by account if provided
	if input.AccountID != nil {
		return uc.transactionRepo.GetByAccountID(ctx, *input.AccountID)
	}

	// Filter by category if provided
	if input.CategoryID != nil {
		return uc.transactionRepo.GetByCategoryID(ctx, *input.CategoryID)
	}

	// Filter by type if provided
	if input.Type != nil {
		return uc.transactionRepo.GetByType(ctx, input.UserID, *input.Type)
	}

	// Default: get all user transactions
	return uc.transactionRepo.GetByUserID(ctx, input.UserID)
}

// ListTransactionsInput represents filters for listing transactions
type ListTransactionsInput struct {
	UserID     uint                    `json:"user_id"`
	AccountID  *uint                   `json:"account_id,omitempty"`
	CategoryID *uint                   `json:"category_id,omitempty"`
	Type       *domain.TransactionType `json:"type,omitempty"`
	StartDate  *time.Time              `json:"start_date,omitempty"`
	EndDate    *time.Time              `json:"end_date,omitempty"`
}
