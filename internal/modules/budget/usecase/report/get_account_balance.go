package report

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// GetAccountBalanceUseCase handles getting account with calculated balance
type GetAccountBalanceUseCase struct {
	accountRepo     interfaces.AccountRepository
	transactionRepo interfaces.TransactionRepository
}

// NewGetAccountBalanceUseCase creates a new use case instance
func NewGetAccountBalanceUseCase(
	accountRepo interfaces.AccountRepository,
	transactionRepo interfaces.TransactionRepository,
) *GetAccountBalanceUseCase {
	return &GetAccountBalanceUseCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// AccountBalanceOutput represents the account with calculated balance
type AccountBalanceOutput struct {
	Account        domain.Account `json:"account"`
	CurrentBalance float64        `json:"current_balance"`
	TotalIncome    float64        `json:"total_income"`
	TotalExpense   float64        `json:"total_expense"`
	TotalTransfers float64        `json:"total_transfers"`
}

// Execute gets account and calculates its current balance
func (uc *GetAccountBalanceUseCase) Execute(ctx context.Context, userID uint, accountID uint) (AccountBalanceOutput, error) {
	// Get account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return AccountBalanceOutput{}, errors.ErrAccountNotFound
	}

	// Verify ownership
	if account.UserID != userID {
		return AccountBalanceOutput{}, errors.ErrUnauthorizedAccess
	}

	// Get all transactions for this account
	transactions, err := uc.transactionRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return AccountBalanceOutput{}, err
	}

	// Calculate balance
	var totalIncome, totalExpense, totalTransfers float64
	for _, tx := range transactions {
		switch tx.Type {
		case domain.TransactionTypeIncome:
			totalIncome += tx.Amount
		case domain.TransactionTypeExpense:
			totalExpense += tx.Amount
		case domain.TransactionTypeTransfer:
			// If this account is the source, it's a debit
			if tx.AccountID == accountID {
				totalTransfers -= tx.Amount
				// Include transfer fee if applicable
				if tx.TransferFee != nil {
					totalTransfers -= *tx.TransferFee
				}
			}
			// If this account is the destination, it's a credit
			if tx.DestinationAccountID != nil && *tx.DestinationAccountID == accountID {
				totalTransfers += tx.Amount
			}
		}
	}

	// Current balance = initial balance + income - expense + transfers
	currentBalance := account.Balance + totalIncome - totalExpense + totalTransfers

	return AccountBalanceOutput{
		Account:        account,
		CurrentBalance: currentBalance,
		TotalIncome:    totalIncome,
		TotalExpense:   totalExpense,
		TotalTransfers: totalTransfers,
	}, nil
}
