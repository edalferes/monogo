package budget

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type GetByIDUseCase struct {
	budgetRepo interfaces.BudgetRepository
	logger     logger.Logger
}

func NewGetByIDUseCase(budgetRepo interfaces.BudgetRepository, log logger.Logger) *GetByIDUseCase {
	return &GetByIDUseCase{
		budgetRepo: budgetRepo,
		logger:     log.With().Str("usecase", "budget.getbyid").Logger(),
	}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, budgetID uint) (domain.Budget, error) {
	uc.logger.Debug().Uint("user_id", userID).Uint("budget_id", budgetID).Msg("getting budget by id")

	if budgetID == 0 {
		uc.logger.Error().Msg("invalid budget_id: cannot be zero")
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", budgetID).Msg("budget not found")
		return domain.Budget{}, err
	}

	if budgetEntity.UserID != userID {
		uc.logger.Warn().
			Uint("budget_user_id", budgetEntity.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: budget belongs to different user")
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	uc.logger.Info().Uint("budget_id", budgetID).Msg("budget retrieved successfully")
	return budgetEntity, nil
}
