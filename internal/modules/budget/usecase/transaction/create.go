package transaction

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type CreateUseCase struct {
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
	categoryRepo    interfaces.CategoryRepository
}

func NewCreateUseCase(transactionRepo interfaces.TransactionRepository, accountRepo interfaces.AccountRepository, categoryRepo interfaces.CategoryRepository) *CreateUseCase {
	return &CreateUseCase{transactionRepo: transactionRepo, accountRepo: accountRepo, categoryRepo: categoryRepo}
}

type CreateInput struct {
	UserID               uint
	AccountID            uint
	CategoryID           uint
	Type                 domain.TransactionType
	Amount               float64
	Description          string
	Date                 string
	DestinationAccountID *uint
}

func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.Transaction, error) {
	if input.UserID == 0 {
		return domain.Transaction{}, errors.ErrInvalidUserID
	}
	if input.Amount <= 0 {
		return domain.Transaction{}, errors.ErrInvalidAmount
	}
	if !isValidTransactionType(input.Type) {
		return domain.Transaction{}, errors.ErrInvalidTransactionType
	}

	account, err := uc.accountRepo.GetByID(ctx, input.AccountID)
	if err != nil {
		return domain.Transaction{}, errors.ErrAccountNotFound
	}
	if account.UserID != input.UserID {
		return domain.Transaction{}, errors.ErrUnauthorizedAccess
	}

	category, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		return domain.Transaction{}, errors.ErrCategoryNotFound
	}
	if category.UserID != input.UserID {
		return domain.Transaction{}, errors.ErrUnauthorizedAccess
	}

	// Normalize destination_account_id: treat 0 as nil
	var destinationAccountID *uint
	if input.DestinationAccountID != nil && *input.DestinationAccountID > 0 {
		destinationAccountID = input.DestinationAccountID

		// Validate destination account for transfers
		destAccount, err := uc.accountRepo.GetByID(ctx, *input.DestinationAccountID)
		if err != nil {
			return domain.Transaction{}, errors.ErrAccountNotFound
		}
		if destAccount.UserID != input.UserID {
			return domain.Transaction{}, errors.ErrUnauthorizedAccess
		}
	}

	// Parse date
	date, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		return domain.Transaction{}, errors.ErrInvalidDate
	}

	// Create the main transaction (debit from source account)
	tx := domain.Transaction{
		UserID:               input.UserID,
		AccountID:            input.AccountID,
		CategoryID:           input.CategoryID,
		Type:                 input.Type,
		Amount:               input.Amount,
		Description:          input.Description,
		Date:                 date,
		Status:               domain.TransactionStatusCompleted,
		DestinationAccountID: destinationAccountID,
	}

	return uc.transactionRepo.Create(ctx, tx)
}
