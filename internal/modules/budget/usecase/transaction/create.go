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
	budgetRepo      interfaces.BudgetRepository
}

func NewCreateUseCase(transactionRepo interfaces.TransactionRepository, accountRepo interfaces.AccountRepository, categoryRepo interfaces.CategoryRepository, budgetRepo interfaces.BudgetRepository) *CreateUseCase {
	return &CreateUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		budgetRepo:      budgetRepo,
	}
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

	createdTx, err := uc.transactionRepo.Create(ctx, tx)
	if err != nil {
		return domain.Transaction{}, err
	}

	// Update budget spent if this is an expense transaction
	if input.Type == domain.TransactionTypeExpense {
		go uc.updateBudgetSpent(context.Background(), input.UserID, input.CategoryID, date)
	}

	return createdTx, nil
}

// updateBudgetSpent updates the spent amount for active budgets of the category
func (uc *CreateUseCase) updateBudgetSpent(ctx context.Context, userID, categoryID uint, transactionDate time.Time) {
	// Get active budgets for this category
	budgets, err := uc.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return
	}

	for _, budget := range budgets {
		// Check if budget is for this category and transaction date is within period
		if budget.CategoryID == categoryID &&
			budget.IsActive &&
			!transactionDate.Before(budget.StartDate) &&
			!transactionDate.After(budget.EndDate) {

			// Calculate total spent for this budget
			transactions, err := uc.transactionRepo.GetByDateRange(ctx, userID, budget.StartDate, budget.EndDate)
			if err != nil {
				continue
			}

			var spent float64
			for _, tx := range transactions {
				if tx.CategoryID == categoryID && tx.Type == domain.TransactionTypeExpense {
					spent += tx.Amount
				}
			}

			// Update budget spent
			_ = uc.budgetRepo.UpdateSpent(ctx, budget.ID, spent)
		}
	}
}
