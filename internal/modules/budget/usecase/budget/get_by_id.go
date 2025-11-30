package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

type GetByIDUseCase struct {
	budgetRepo repository.BudgetRepository
}

func NewGetByIDUseCase(budgetRepo repository.BudgetRepository) *GetByIDUseCase {
	return &GetByIDUseCase{budgetRepo: budgetRepo}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, budgetID uint) (domain.Budget, error) {
	if budgetID == 0 {
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return domain.Budget{}, err
	}

	if budgetEntity.UserID != userID {
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	return budgetEntity, nil
}
