package usecase

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// CreateBudgetUseCase handles budget creation
type CreateBudgetUseCase struct {
	budgetRepo   repository.BudgetRepository
	categoryRepo repository.CategoryRepository
}

// NewCreateBudgetUseCase creates a new use case instance
func NewCreateBudgetUseCase(
	budgetRepo repository.BudgetRepository,
	categoryRepo repository.CategoryRepository,
) *CreateBudgetUseCase {
	return &CreateBudgetUseCase{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

// Execute creates a new budget
func (uc *CreateBudgetUseCase) Execute(ctx context.Context, input CreateBudgetInput) (domain.Budget, error) {
	// Validate input
	if input.UserID == 0 {
		return domain.Budget{}, errors.ErrInvalidUserID
	}
	if input.CategoryID == 0 {
		return domain.Budget{}, errors.ErrInvalidCategoryID
	}
	if input.Name == "" {
		return domain.Budget{}, errors.ErrBudgetNameRequired
	}
	if input.Amount <= 0 {
		return domain.Budget{}, errors.ErrInvalidBudgetAmount
	}
	if !isValidBudgetPeriod(input.Period) {
		return domain.Budget{}, errors.ErrInvalidBudgetPeriod
	}
	if input.StartDate.After(input.EndDate) {
		return domain.Budget{}, errors.ErrInvalidDateRange
	}

	// Verify category exists and belongs to user
	category, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		return domain.Budget{}, errors.ErrCategoryNotFound
	}
	if category.UserID != input.UserID {
		return domain.Budget{}, errors.ErrUnauthorizedAccess
	}

	// Create budget domain entity
	budget := domain.Budget{
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Amount:      input.Amount,
		Spent:       0, // Initial spent is 0
		Period:      input.Period,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		AlertAt:     input.AlertAt,
		IsActive:    true,
		Description: input.Description,
	}

	// Save to repository
	return uc.budgetRepo.Create(ctx, budget)
}

// CreateBudgetInput represents the input for creating a budget
type CreateBudgetInput struct {
	UserID      uint                `json:"user_id"`
	CategoryID  uint                `json:"category_id"`
	Name        string              `json:"name"`
	Amount      float64             `json:"amount"`
	Period      domain.BudgetPeriod `json:"period"`
	StartDate   time.Time           `json:"start_date"`
	EndDate     time.Time           `json:"end_date"`
	AlertAt     *float64            `json:"alert_at,omitempty"`
	Description string              `json:"description,omitempty"`
}

func isValidBudgetPeriod(period domain.BudgetPeriod) bool {
	switch period {
	case domain.BudgetPeriodMonthly, domain.BudgetPeriodQuarterly, domain.BudgetPeriodYearly, domain.BudgetPeriodCustom:
		return true
	default:
		return false
	}
}
