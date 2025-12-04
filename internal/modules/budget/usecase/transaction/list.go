package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type ListUseCase struct {
	transactionRepo interfaces.TransactionRepository
}

func NewListUseCase(transactionRepo interfaces.TransactionRepository) *ListUseCase {
	return &ListUseCase{transactionRepo: transactionRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	return uc.transactionRepo.GetByUserID(ctx, userID)
}
