package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

type DeleteUseCase struct {
	budgetRepo repository.BudgetRepository
}

func NewDeleteUseCase(budgetRepo repository.BudgetRepository) *DeleteUseCase {
	return &DeleteUseCase{budgetRepo: budgetRepo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, userID, budgetID uint) error {
	if budgetID == 0 {
		return errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return err
	}

	if budgetEntity.UserID != userID {
		return errors.ErrBudgetNotFound
	}

	return uc.budgetRepo.Delete(ctx, budgetID)
}
