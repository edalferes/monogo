package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

type DeleteUseCase struct {
	transactionRepo repository.TransactionRepository
}

func NewDeleteUseCase(transactionRepo repository.TransactionRepository) *DeleteUseCase {
	return &DeleteUseCase{transactionRepo: transactionRepo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, userID, transactionID uint) error {
	if transactionID == 0 {
		return errors.ErrTransactionNotFound
	}

	tx, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if tx.UserID != userID {
		return errors.ErrTransactionNotFound
	}

	return uc.transactionRepo.Delete(ctx, transactionID)
}
