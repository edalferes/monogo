package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type ListUseCase struct {
	budgetRepo      interfaces.BudgetRepository
	transactionRepo interfaces.TransactionRepository
}

func NewListUseCase(budgetRepo interfaces.BudgetRepository, transactionRepo interfaces.TransactionRepository) *ListUseCase {
	return &ListUseCase{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Budget, error) {
	budgets, err := uc.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate spent for each budget based on transactions
	for i := range budgets {
		spent, err := uc.calculateSpent(ctx, budgets[i])
		if err != nil {
			// Log error but continue with other budgets
			continue
		}
		budgets[i].Spent = spent
	}

	return budgets, nil
}

// calculateSpent calculates total spent for a budget based on transactions
func (uc *ListUseCase) calculateSpent(ctx context.Context, budget domain.Budget) (float64, error) {
	// Get transactions by category within budget period
	transactions, err := uc.transactionRepo.GetByDateRange(ctx, budget.UserID, budget.StartDate, budget.EndDate)
	if err != nil {
		return 0, err
	}

	var spent float64
	for _, tx := range transactions {
		// Only count expenses for this category
		if tx.CategoryID == budget.CategoryID && tx.Type == domain.TransactionTypeExpense {
			spent += tx.Amount
		}
	}

	return spent, nil
}
