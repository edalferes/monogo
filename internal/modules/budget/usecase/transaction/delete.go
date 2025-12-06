package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type DeleteUseCase struct {
	transactionRepo interfaces.TransactionRepository
	logger          logger.Logger
}

func NewDeleteUseCase(transactionRepo interfaces.TransactionRepository, log logger.Logger) *DeleteUseCase {
	return &DeleteUseCase{
		transactionRepo: transactionRepo,
		logger:          log.With().Str("usecase", "transaction.delete").Logger(),
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, userID, transactionID uint) error {
	uc.logger.Debug().Uint("user_id", userID).Uint("transaction_id", transactionID).Msg("deleting transaction")

	if transactionID == 0 {
		uc.logger.Error().Msg("invalid transaction_id: cannot be zero")
		return errors.ErrTransactionNotFound
	}

	tx, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("transaction_id", transactionID).Msg("transaction not found")
		return err
	}

	if tx.UserID != userID {
		uc.logger.Warn().
			Uint("transaction_user_id", tx.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: transaction belongs to different user")
		return errors.ErrTransactionNotFound
	}

	err = uc.transactionRepo.Delete(ctx, transactionID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("transaction_id", transactionID).Msg("failed to delete transaction")
		return err
	}

	uc.logger.Info().Uint("transaction_id", transactionID).Msg("transaction deleted successfully")
	return nil
}
