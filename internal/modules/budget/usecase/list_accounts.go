package usecase

import (
	"context"

	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// ListAccountsUseCase handles listing user accounts
type ListAccountsUseCase struct {
	accountRepo repository.AccountRepository
}

// NewListAccountsUseCase creates a new use case instance
func NewListAccountsUseCase(accountRepo repository.AccountRepository) *ListAccountsUseCase {
	return &ListAccountsUseCase{
		accountRepo: accountRepo,
	}
}

// Execute lists all accounts for a user
func (uc *ListAccountsUseCase) Execute(ctx context.Context, userID uint) ([]domain.Account, error) {
	return uc.accountRepo.GetByUserID(ctx, userID)
}
