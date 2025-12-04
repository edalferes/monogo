package transaction

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
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

	// Validate destination account for transfers
	if input.DestinationAccountID != nil {
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

	tx := domain.Transaction{
		UserID:               input.UserID,
		AccountID:            input.AccountID,
		CategoryID:           input.CategoryID,
		Type:                 input.Type,
		Amount:               input.Amount,
		Description:          input.Description,
		Date:                 date,
		Status:               domain.TransactionStatusCompleted,
		DestinationAccountID: input.DestinationAccountID,
	}

	return uc.transactionRepo.Create(ctx, tx)
}
