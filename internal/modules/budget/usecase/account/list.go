package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// ListUseCase handles listing user accounts
type ListUseCase struct {
	accountRepo interfaces.AccountRepository
}

// NewListUseCase creates a new use case instance
func NewListUseCase(accountRepo interfaces.AccountRepository) *ListUseCase {
	return &ListUseCase{
		accountRepo: accountRepo,
	}
}

// Execute lists all accounts for a user
func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Account, error) {
	return uc.accountRepo.GetByUserID(ctx, userID)
}
