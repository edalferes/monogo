package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type ListUseCase struct {
	budgetRepo      interfaces.BudgetRepository
	transactionRepo interfaces.TransactionRepository
	logger          logger.Logger
}

func NewListUseCase(budgetRepo interfaces.BudgetRepository, transactionRepo interfaces.TransactionRepository, log logger.Logger) *ListUseCase {
	return &ListUseCase{
		budgetRepo:      budgetRepo,
		transactionRepo: transactionRepo,
		logger:          log.With().Str("usecase", "budget.list").Logger(),
	}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Budget, error) {
	uc.logger.Debug().Uint("user_id", userID).Msg("listing budgets")

	budgets, err := uc.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("user_id", userID).Msg("failed to get budgets")
		return nil, err
	}

	// Calculate spent for each budget based on transactions
	for i := range budgets {
		spent, err := uc.calculateSpent(ctx, budgets[i])
		if err != nil {
			uc.logger.Error().Err(err).Uint("budget_id", budgets[i].ID).Msg("failed to calculate spent, skipping")
			continue
		}
		budgets[i].Spent = spent
	}

	uc.logger.Info().Uint("user_id", userID).Int("count", len(budgets)).Msg("budgets listed successfully")
	return budgets, nil
}

// calculateSpent calculates total spent for a budget based on transactions
func (uc *ListUseCase) calculateSpent(ctx context.Context, budget domain.Budget) (float64, error) {
	// Get transactions by category within budget period
	transactions, err := uc.transactionRepo.GetByDateRange(ctx, budget.UserID, budget.StartDate, budget.EndDate)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", budget.ID).Msg("failed to get transactions for spent calculation")
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
