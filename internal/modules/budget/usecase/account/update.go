package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// UpdateInput represents the input for updating an account
type UpdateInput struct {
	ID          uint
	UserID      uint
	Name        *string
	Type        *domain.AccountType
	Currency    *string
	Description *string
}

// UpdateUseCase handles account updates
type UpdateUseCase struct {
	accountRepo interfaces.AccountRepository
	logger      logger.Logger
}

// NewUpdateUseCase creates a new use case instance
func NewUpdateUseCase(accountRepo interfaces.AccountRepository, log logger.Logger) *UpdateUseCase {
	return &UpdateUseCase{
		accountRepo: accountRepo,
		logger:      log.With().Str("usecase", "account.update").Logger(),
	}
}

// Execute updates an existing account
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Account, error) {
	uc.logger.Debug().Uint("user_id", input.UserID).Uint("account_id", input.ID).Msg("updating account")

	if input.ID == 0 {
		uc.logger.Error().Msg("invalid account_id: cannot be zero")
		return domain.Account{}, errors.ErrAccountNotFound
	}

	// Get existing account
	account, err := uc.accountRepo.GetByID(ctx, input.ID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("account_id", input.ID).Msg("account not found")
		return domain.Account{}, err
	}

	// Ensure the account belongs to the user
	if account.UserID != input.UserID {
		uc.logger.Warn().
			Uint("account_user_id", account.UserID).
			Uint("request_user_id", input.UserID).
			Msg("unauthorized access: account belongs to different user")
		return domain.Account{}, errors.ErrAccountNotFound
	}

	// Update fields if provided
	if input.Name != nil {
		if *input.Name == "" {
			uc.logger.Error().Msg("validation error: account name cannot be empty")
			return domain.Account{}, errors.ErrAccountNameRequired
		}
		account.Name = *input.Name
	}

	if input.Type != nil {
		if !isValidAccountType(*input.Type) {
			uc.logger.Error().Str("type", string(*input.Type)).Msg("validation error: invalid account type")
			return domain.Account{}, errors.ErrInvalidAccountType
		}
		account.Type = *input.Type
	}

	if input.Currency != nil {
		if *input.Currency == "" {
			uc.logger.Error().Msg("validation error: currency cannot be empty")
			return domain.Account{}, errors.ErrInvalidData
		}
		account.Currency = *input.Currency
	}

	if input.Description != nil {
		account.Description = *input.Description
	}

	// Save updates
	updatedAccount, err := uc.accountRepo.Update(ctx, account)
	if err != nil {
		uc.logger.Error().Err(err).Uint("account_id", input.ID).Msg("failed to update account")
		return domain.Account{}, err
	}

	uc.logger.Info().Uint("account_id", input.ID).Msg("account updated successfully")
	return updatedAccount, nil
}
