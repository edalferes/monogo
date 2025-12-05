package account

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// CreateUseCase handles account creation
type CreateUseCase struct {
	accountRepo interfaces.AccountRepository
}

// NewCreateUseCase creates a new use case instance
func NewCreateUseCase(accountRepo interfaces.AccountRepository) *CreateUseCase {
	return &CreateUseCase{
		accountRepo: accountRepo,
	}
}

// CreateInput represents the input for creating an account
type CreateInput struct {
	UserID         uint               `json:"user_id" validate:"required"`
	Name           string             `json:"name" validate:"required"`
	Type           domain.AccountType `json:"type" validate:"required"`
	InitialBalance float64            `json:"initial_balance"`
	Currency       string             `json:"currency"`
	Description    string             `json:"description,omitempty"`
}

// Execute creates a new account
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.Account, error) {
	// Validate input
	if input.UserID == 0 {
		return domain.Account{}, errors.ErrInvalidUserID
	}
	if input.Name == "" {
		return domain.Account{}, errors.ErrAccountNameRequired
	}
	if !isValidAccountType(input.Type) {
		return domain.Account{}, errors.ErrInvalidAccountType
	}
	if input.Currency == "" {
		input.Currency = "BRL" // Default currency
	}

	// Create account domain entity
	account := domain.Account{
		UserID:      input.UserID,
		Name:        input.Name,
		Type:        input.Type,
		Balance:     input.InitialBalance,
		Currency:    input.Currency,
		Description: input.Description,
		IsActive:    true,
	}

	// Save to repository
	return uc.accountRepo.Create(ctx, account)
}
