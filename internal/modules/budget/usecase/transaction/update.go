package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// UpdateInput represents the input for updating a transaction
type UpdateInput struct {
	ID          uint
	UserID      uint
	AccountID   *uint
	CategoryID  *uint
	Type        *domain.TransactionType
	Amount      *float64
	Description *string
	Date        *string
}

// UpdateUseCase handles transaction updates
type UpdateUseCase struct {
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
	categoryRepo    interfaces.CategoryRepository
}

// NewUpdateUseCase creates a new use case instance
func NewUpdateUseCase(
	transactionRepo interfaces.TransactionRepository,
	accountRepo interfaces.AccountRepository,
	categoryRepo interfaces.CategoryRepository,
) *UpdateUseCase {
	return &UpdateUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
	}
}

// Execute updates an existing transaction
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Transaction, error) {
	if input.ID == 0 {
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Transaction{}, err
	}

	// Ensure the transaction belongs to the user
	if tx.UserID != input.UserID {
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	// Validate and update account if provided
	if input.AccountID != nil {
		account, err := uc.accountRepo.GetByID(ctx, *input.AccountID)
		if err != nil {
			return domain.Transaction{}, errors.ErrAccountNotFound
		}
		if account.UserID != input.UserID {
			return domain.Transaction{}, errors.ErrUnauthorizedAccess
		}
		tx.AccountID = *input.AccountID
	}

	// Validate and update category if provided
	if input.CategoryID != nil {
		category, err := uc.categoryRepo.GetByID(ctx, *input.CategoryID)
		if err != nil {
			return domain.Transaction{}, errors.ErrCategoryNotFound
		}
		if category.UserID != input.UserID {
			return domain.Transaction{}, errors.ErrUnauthorizedAccess
		}
		tx.CategoryID = *input.CategoryID
	}

	// Update type if provided
	if input.Type != nil {
		if !isValidTransactionType(*input.Type) {
			return domain.Transaction{}, errors.ErrInvalidTransactionType
		}
		tx.Type = *input.Type
	}

	// Update amount if provided
	if input.Amount != nil {
		if *input.Amount <= 0 {
			return domain.Transaction{}, errors.ErrInvalidAmount
		}
		tx.Amount = *input.Amount
	}

	// Update description if provided
	if input.Description != nil {
		tx.Description = *input.Description
	}

	return uc.transactionRepo.Update(ctx, tx)
}
