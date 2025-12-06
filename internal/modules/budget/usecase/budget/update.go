package budget

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type UpdateInput struct {
	ID          uint
	UserID      uint
	Name        *string
	Amount      *float64
	Period      *domain.BudgetPeriod
	StartDate   *time.Time
	EndDate     *time.Time
	AlertAt     *float64
	Description *string
}

type UpdateUseCase struct {
	budgetRepo interfaces.BudgetRepository
	logger     logger.Logger
}

func NewUpdateUseCase(budgetRepo interfaces.BudgetRepository, log logger.Logger) *UpdateUseCase {
	return &UpdateUseCase{
		budgetRepo: budgetRepo,
		logger:     log.With().Str("usecase", "budget.update").Logger(),
	}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Budget, error) {
	uc.logger.Debug().Uint("budget_id", input.ID).Uint("user_id", input.UserID).Msg("updating budget")

	if input.ID == 0 {
		uc.logger.Error().Msg("invalid budget_id: cannot be zero")
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, input.ID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", input.ID).Msg("budget not found")
		return domain.Budget{}, err
	}

	if budgetEntity.UserID != input.UserID {
		uc.logger.Warn().
			Uint("budget_user_id", budgetEntity.UserID).
			Uint("request_user_id", input.UserID).
			Msg("unauthorized access: budget belongs to different user")
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	if input.Name != nil {
		if *input.Name == "" {
			uc.logger.Error().Msg("budget name cannot be empty")
			return domain.Budget{}, errors.ErrBudgetNameRequired
		}
		budgetEntity.Name = *input.Name
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			uc.logger.Error().Msg("invalid amount: must be positive")
			return domain.Budget{}, errors.ErrInvalidBudgetAmount
		}
		budgetEntity.Amount = *input.Amount
	}

	if input.Period != nil {
		if !isValidBudgetPeriod(*input.Period) {
			uc.logger.Error().Str("period", string(*input.Period)).Msg("invalid budget period")
			return domain.Budget{}, errors.ErrInvalidBudgetPeriod
		}
		budgetEntity.Period = *input.Period
	}

	if input.StartDate != nil {
		budgetEntity.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		budgetEntity.EndDate = *input.EndDate
	}

	if budgetEntity.StartDate.After(budgetEntity.EndDate) || budgetEntity.StartDate.Equal(budgetEntity.EndDate) {
		uc.logger.Error().Msg("invalid date range: start_date must be before end_date")
		return domain.Budget{}, errors.ErrInvalidDateRange
	}

	if input.AlertAt != nil {
		budgetEntity.AlertAt = input.AlertAt
	}

	if input.Description != nil {
		budgetEntity.Description = *input.Description
	}

	updatedBudget, err := uc.budgetRepo.Update(ctx, budgetEntity)
	if err != nil {
		uc.logger.Error().Err(err).Uint("budget_id", input.ID).Msg("failed to update budget")
		return domain.Budget{}, err
	}

	uc.logger.Info().Uint("budget_id", input.ID).Msg("budget updated successfully")
	return updatedBudget, nil
}
