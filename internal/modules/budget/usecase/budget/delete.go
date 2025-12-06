package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type DeleteUseCase struct {
	budgetRepo interfaces.BudgetRepository
	logger     logger.Logger
}

func NewDeleteUseCase(budgetRepo interfaces.BudgetRepository, log logger.Logger) *DeleteUseCase {
	return &DeleteUseCase{
		budgetRepo: budgetRepo,
		logger:     log.With().Str("usecase", "budget.delete").Logger(),
	}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, userID, budgetID uint) error {
	uc.logger.Debug().Uint("user_id", userID).Uint("budget_id", budgetID).Msg("deleting budget")

	if budgetID == 0 {
		uc.logger.Error().Msg("invalid budget_id: cannot be zero")
		return errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", budgetID).Msg("budget not found")
		return err
	}

	if budgetEntity.UserID != userID {
		uc.logger.Warn().
			Uint("budget_user_id", budgetEntity.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: budget belongs to different user")
		return errors.ErrBudgetNotFound
	}

	err = uc.budgetRepo.Delete(ctx, budgetID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", budgetID).Msg("failed to delete budget")
		return err
	}

	uc.logger.Info().Uint("budget_id", budgetID).Msg("budget deleted successfully")
	return nil
}
