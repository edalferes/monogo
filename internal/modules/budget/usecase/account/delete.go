package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// DeleteUseCase handles account deletion (soft delete)
type DeleteUseCase struct {
	accountRepo interfaces.AccountRepository
	logger      logger.Logger
}

// NewDeleteUseCase creates a new use case instance
func NewDeleteUseCase(accountRepo interfaces.AccountRepository, log logger.Logger) *DeleteUseCase {
	return &DeleteUseCase{
		accountRepo: accountRepo,
		logger:      log.With().Str("usecase", "account.delete").Logger(),
	}
}

// Execute soft-deletes an account
func (uc *DeleteUseCase) Execute(ctx context.Context, userID, accountID uint) error {
	uc.logger.Debug().Uint("user_id", userID).Uint("account_id", accountID).Msg("deleting account")

	if accountID == 0 {
		uc.logger.Error().Msg("invalid account_id: cannot be zero")
		return errors.ErrAccountNotFound
	}

	// Get existing account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("account_id", accountID).Msg("account not found")
		return err
	}

	// Ensure the account belongs to the user
	if account.UserID != userID {
		uc.logger.Warn().
			Uint("account_user_id", account.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: account belongs to different user")
		return errors.ErrAccountNotFound
	}

	// Soft delete by marking as inactive
	err = uc.accountRepo.Delete(ctx, accountID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("account_id", accountID).Msg("failed to delete account")
		return err
	}

	uc.logger.Info().Uint("account_id", accountID).Msg("account deleted successfully")
	return nil
}
