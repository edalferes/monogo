package usecase

import (
	"context"

	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
	"github.com/edalferes/monogo/internal/modules/budget/errors"
)

// CreateAccountUseCase handles account creation
type CreateAccountUseCase struct {
	accountRepo repository.AccountRepository
}

// NewCreateAccountUseCase creates a new use case instance
func NewCreateAccountUseCase(accountRepo repository.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepo: accountRepo,
	}
}

// Execute creates a new account
func (uc *CreateAccountUseCase) Execute(ctx context.Context, input CreateAccountInput) (domain.Account, error) {
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

// CreateAccountInput represents the input for creating an account
type CreateAccountInput struct {
	UserID         uint               `json:"user_id"`
	Name           string             `json:"name"`
	Type           domain.AccountType `json:"type"`
	InitialBalance float64            `json:"initial_balance"`
	Currency       string             `json:"currency"`
	Description    string             `json:"description,omitempty"`
}

func isValidAccountType(accountType domain.AccountType) bool {
	switch accountType {
	case domain.AccountTypeChecking, domain.AccountTypeSavings, domain.AccountTypeCredit, domain.AccountTypeCash, domain.AccountTypeInvest:
		return true
	default:
		return false
	}
}
