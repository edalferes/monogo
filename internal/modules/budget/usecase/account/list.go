package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// ListUseCase handles listing user accounts
type ListUseCase struct {
	accountRepo repository.AccountRepository
}

// NewListUseCase creates a new use case instance
func NewListUseCase(accountRepo repository.AccountRepository) *ListUseCase {
	return &ListUseCase{
		accountRepo: accountRepo,
	}
}

// Execute lists all accounts for a user
func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Account, error) {
	return uc.accountRepo.GetByUserID(ctx, userID)
}
