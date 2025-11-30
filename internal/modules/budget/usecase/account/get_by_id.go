package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// GetByIDUseCase handles retrieving a specific account
type GetByIDUseCase struct {
	accountRepo repository.AccountRepository
}

// NewGetByIDUseCase creates a new use case instance
func NewGetByIDUseCase(accountRepo repository.AccountRepository) *GetByIDUseCase {
	return &GetByIDUseCase{
		accountRepo: accountRepo,
	}
}

// Execute retrieves an account by ID
func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, accountID uint) (domain.Account, error) {
	if accountID == 0 {
		return domain.Account{}, errors.ErrAccountNotFound
	}

	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return domain.Account{}, err
	}

	// Ensure the account belongs to the user
	if account.UserID != userID {
		return domain.Account{}, errors.ErrAccountNotFound
	}

	return account, nil
}
