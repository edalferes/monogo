package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
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
	accountRepo repository.AccountRepository
}

// NewUpdateUseCase creates a new use case instance
func NewUpdateUseCase(accountRepo repository.AccountRepository) *UpdateUseCase {
	return &UpdateUseCase{
		accountRepo: accountRepo,
	}
}

// Execute updates an existing account
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Account, error) {
	if input.ID == 0 {
		return domain.Account{}, errors.ErrAccountNotFound
	}

	// Get existing account
	account, err := uc.accountRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Account{}, err
	}

	// Ensure the account belongs to the user
	if account.UserID != input.UserID {
		return domain.Account{}, errors.ErrAccountNotFound
	}

	// Update fields if provided
	if input.Name != nil {
		if *input.Name == "" {
			return domain.Account{}, errors.ErrAccountNameRequired
		}
		account.Name = *input.Name
	}

	if input.Type != nil {
		if !isValidAccountType(*input.Type) {
			return domain.Account{}, errors.ErrInvalidAccountType
		}
		account.Type = *input.Type
	}

	if input.Currency != nil {
		if *input.Currency == "" {
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
		return domain.Account{}, err
	}

	return updatedAccount, nil
}
