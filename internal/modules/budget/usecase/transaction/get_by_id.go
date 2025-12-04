package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

type GetByIDUseCase struct {
	transactionRepo interfaces.TransactionRepository
}

func NewGetByIDUseCase(transactionRepo interfaces.TransactionRepository) *GetByIDUseCase {
	return &GetByIDUseCase{transactionRepo: transactionRepo}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, transactionID uint) (domain.Transaction, error) {
	if transactionID == 0 {
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	tx, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return domain.Transaction{}, err
	}

	if tx.UserID != userID {
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	return tx, nil
}
