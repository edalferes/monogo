package usecase

import (
	"context"

	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// ListBudgetsUseCase handles listing user budgets
type ListBudgetsUseCase struct {
	budgetRepo repository.BudgetRepository
}

// NewListBudgetsUseCase creates a new use case instance
func NewListBudgetsUseCase(budgetRepo repository.BudgetRepository) *ListBudgetsUseCase {
	return &ListBudgetsUseCase{
		budgetRepo: budgetRepo,
	}
}

// Execute lists budgets for a user with optional filters
func (uc *ListBudgetsUseCase) Execute(ctx context.Context, input ListBudgetsInput) ([]domain.Budget, error) {
	// Filter by category if provided
	if input.CategoryID != nil {
		return uc.budgetRepo.GetByCategoryID(ctx, *input.CategoryID)
	}

	// Filter active budgets if requested
	if input.ActiveOnly {
		return uc.budgetRepo.GetActive(ctx, input.UserID)
	}

	// Default: get all user budgets
	return uc.budgetRepo.GetByUserID(ctx, input.UserID)
}

// ListBudgetsInput represents filters for listing budgets
type ListBudgetsInput struct {
	UserID     uint  `json:"user_id"`
	CategoryID *uint `json:"category_id,omitempty"`
	ActiveOnly bool  `json:"active_only"`
}
