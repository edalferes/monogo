package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type GetByIDUseCase struct {
	transactionRepo interfaces.TransactionRepository
	logger          logger.Logger
}

func NewGetByIDUseCase(transactionRepo interfaces.TransactionRepository, log logger.Logger) *GetByIDUseCase {
	return &GetByIDUseCase{
		transactionRepo: transactionRepo,
		logger:          log.With().Str("usecase", "transaction.getbyid").Logger(),
	}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, transactionID uint) (domain.Transaction, error) {
	uc.logger.Debug().Uint("user_id", userID).Uint("transaction_id", transactionID).Msg("getting transaction by id")

	if transactionID == 0 {
		uc.logger.Error().Msg("invalid transaction_id: cannot be zero")
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	tx, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("transaction_id", transactionID).Msg("transaction not found")
		return domain.Transaction{}, err
	}

	if tx.UserID != userID {
		uc.logger.Warn().
			Uint("transaction_user_id", tx.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: transaction belongs to different user")
		return domain.Transaction{}, errors.ErrTransactionNotFound
	}

	uc.logger.Info().Uint("transaction_id", transactionID).Msg("transaction retrieved successfully")
	return tx, nil
}
