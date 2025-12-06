package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// ListUseCase handles listing user accounts
type ListUseCase struct {
	accountRepo interfaces.AccountRepository
	logger      logger.Logger
}

// NewListUseCase creates a new use case instance
func NewListUseCase(accountRepo interfaces.AccountRepository, log logger.Logger) *ListUseCase {
	return &ListUseCase{
		accountRepo: accountRepo,
		logger:      log.With().Str("usecase", "account.list").Logger(),
	}
}

// Execute lists all accounts for a user
func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Account, error) {
	uc.logger.Debug().Uint("user_id", userID).Msg("listing accounts")

	accounts, err := uc.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("user_id", userID).Msg("failed to list accounts")
		return nil, err
	}

	uc.logger.Info().Uint("user_id", userID).Int("count", len(accounts)).Msg("accounts listed successfully")
	return accounts, nil
}
