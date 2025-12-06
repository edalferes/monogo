package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// GetByIDUseCase handles retrieving a specific account
type GetByIDUseCase struct {
	accountRepo interfaces.AccountRepository
	logger      logger.Logger
}

// NewGetByIDUseCase creates a new use case instance
func NewGetByIDUseCase(accountRepo interfaces.AccountRepository, log logger.Logger) *GetByIDUseCase {
	return &GetByIDUseCase{
		accountRepo: accountRepo,
		logger:      log.With().Str("usecase", "account.getbyid").Logger(),
	}
}

// Execute retrieves an account by ID
func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, accountID uint) (domain.Account, error) {
	uc.logger.Debug().Uint("user_id", userID).Uint("account_id", accountID).Msg("getting account by id")

	if accountID == 0 {
		uc.logger.Error().Msg("invalid account_id: cannot be zero")
		return domain.Account{}, errors.ErrAccountNotFound
	}

	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("account_id", accountID).Msg("account not found")
		return domain.Account{}, err
	}

	// Ensure the account belongs to the user
	if account.UserID != userID {
		uc.logger.Warn().
			Uint("account_user_id", account.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: account belongs to different user")
		return domain.Account{}, errors.ErrAccountNotFound
	}

	uc.logger.Info().Uint("account_id", accountID).Msg("account retrieved successfully")
	return account, nil
}
