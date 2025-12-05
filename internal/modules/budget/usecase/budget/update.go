package budget

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
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
}

func NewUpdateUseCase(budgetRepo interfaces.BudgetRepository) *UpdateUseCase {
	return &UpdateUseCase{budgetRepo: budgetRepo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Budget, error) {
	if input.ID == 0 {
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	budgetEntity, err := uc.budgetRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Budget{}, err
	}

	if budgetEntity.UserID != input.UserID {
		return domain.Budget{}, errors.ErrBudgetNotFound
	}

	if input.Name != nil {
		if *input.Name == "" {
			return domain.Budget{}, errors.ErrBudgetNameRequired
		}
		budgetEntity.Name = *input.Name
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			return domain.Budget{}, errors.ErrInvalidBudgetAmount
		}
		budgetEntity.Amount = *input.Amount
	}

	if input.Period != nil {
		if !isValidBudgetPeriod(*input.Period) {
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
		return domain.Budget{}, errors.ErrInvalidDateRange
	}

	if input.AlertAt != nil {
		budgetEntity.AlertAt = input.AlertAt
	}

	if input.Description != nil {
		budgetEntity.Description = *input.Description
	}

	return uc.budgetRepo.Update(ctx, budgetEntity)
}
