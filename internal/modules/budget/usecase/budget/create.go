package budget

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

type CreateUseCase struct {
	budgetRepo   interfaces.BudgetRepository
	categoryRepo interfaces.CategoryRepository
}

func NewCreateUseCase(budgetRepo interfaces.BudgetRepository, categoryRepo interfaces.CategoryRepository) *CreateUseCase {
	return &CreateUseCase{budgetRepo: budgetRepo, categoryRepo: categoryRepo}
}

type CreateInput struct {
	UserID      uint
	CategoryID  uint
	Name        string
	Amount      float64
	Period      domain.BudgetPeriod
	StartDate   time.Time
	EndDate     time.Time
	AlertAt     *float64
	Description string
}

func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.Budget, error) {
	if input.UserID == 0 {
		return domain.Budget{}, errors.ErrInvalidUserID
	}
	if input.Name == "" {
		return domain.Budget{}, errors.ErrBudgetNameRequired
	}
	if input.Amount <= 0 {
		return domain.Budget{}, errors.ErrInvalidBudgetAmount
	}
	if input.CategoryID == 0 {
		return domain.Budget{}, errors.ErrInvalidCategoryID
	}
	if !isValidBudgetPeriod(input.Period) {
		return domain.Budget{}, errors.ErrInvalidBudgetPeriod
	}
	if input.StartDate.After(input.EndDate) || input.StartDate.Equal(input.EndDate) {
		return domain.Budget{}, errors.ErrInvalidDateRange
	}

	category, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		return domain.Budget{}, errors.ErrCategoryNotFound
	}
	if category.UserID != input.UserID {
		return domain.Budget{}, errors.ErrUnauthorizedAccess
	}

	budgetEntity := domain.Budget{
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Amount:      input.Amount,
		Period:      input.Period,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		AlertAt:     input.AlertAt,
		Description: input.Description,
		IsActive:    true,
		Spent:       0,
	}

	return uc.budgetRepo.Create(ctx, budgetEntity)
}
