package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type ListUseCase struct {
	budgetRepo repository.BudgetRepository
}

func NewListUseCase(budgetRepo repository.BudgetRepository) *ListUseCase {
	return &ListUseCase{budgetRepo: budgetRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Budget, error) {
	return uc.budgetRepo.GetByUserID(ctx, userID)
}
