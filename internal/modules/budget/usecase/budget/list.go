package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type ListUseCase struct {
	budgetRepo interfaces.BudgetRepository
}

func NewListUseCase(budgetRepo interfaces.BudgetRepository) *ListUseCase {
	return &ListUseCase{budgetRepo: budgetRepo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Budget, error) {
	return uc.budgetRepo.GetByUserID(ctx, userID)
}
