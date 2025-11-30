package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type ListUseCase struct {
	transactionRepo repository.TransactionRepository
}

func NewListUseCase(transactionRepo repository.TransactionRepository) *ListUseCase {
	return &ListUseCase{transactionRepo: transactionRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	return uc.transactionRepo.GetByUserID(ctx, userID)
}
