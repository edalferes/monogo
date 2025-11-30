package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// DeleteUseCase handles account deletion (soft delete)
type DeleteUseCase struct {
	accountRepo repository.AccountRepository
}

// NewDeleteUseCase creates a new use case instance
func NewDeleteUseCase(accountRepo repository.AccountRepository) *DeleteUseCase {
	return &DeleteUseCase{
		accountRepo: accountRepo,
	}
}

// Execute soft-deletes an account
func (uc *DeleteUseCase) Execute(ctx context.Context, userID, accountID uint) error {
	if accountID == 0 {
		return errors.ErrAccountNotFound
	}

	// Get existing account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}

	// Ensure the account belongs to the user
	if account.UserID != userID {
		return errors.ErrAccountNotFound
	}

	// Soft delete by marking as inactive
	return uc.accountRepo.Delete(ctx, accountID)
}
