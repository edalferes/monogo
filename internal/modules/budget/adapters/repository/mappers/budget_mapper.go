package mappers

import (
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// BudgetMapper converts between domain.Budget and models.BudgetModel
type BudgetMapper struct{}

// ToModel converts domain.Budget to models.BudgetModel
func (m BudgetMapper) ToModel(budget domain.Budget) models.BudgetModel {
	return models.BudgetModel{
		ID:          budget.ID,
		UserID:      budget.UserID,
		CategoryID:  budget.CategoryID,
		Name:        budget.Name,
		Amount:      budget.Amount,
		Spent:       budget.Spent,
		Period:      string(budget.Period),
		StartDate:   budget.StartDate,
		EndDate:     budget.EndDate,
		AlertAt:     budget.AlertAt,
		IsActive:    budget.IsActive,
		Description: budget.Description,
		CreatedAt:   budget.CreatedAt,
		UpdatedAt:   budget.UpdatedAt,
	}
}

// ToDomain converts models.BudgetModel to domain.Budget
func (m BudgetMapper) ToDomain(budgetModel models.BudgetModel) domain.Budget {
	return domain.Budget{
		ID:          budgetModel.ID,
		UserID:      budgetModel.UserID,
		CategoryID:  budgetModel.CategoryID,
		Name:        budgetModel.Name,
		Amount:      budgetModel.Amount,
		Spent:       budgetModel.Spent,
		Period:      domain.BudgetPeriod(budgetModel.Period),
		StartDate:   budgetModel.StartDate,
		EndDate:     budgetModel.EndDate,
		AlertAt:     budgetModel.AlertAt,
		IsActive:    budgetModel.IsActive,
		Description: budgetModel.Description,
		CreatedAt:   budgetModel.CreatedAt,
		UpdatedAt:   budgetModel.UpdatedAt,
	}
}

// ToDomainSlice converts []models.BudgetModel to []domain.Budget
func (m BudgetMapper) ToDomainSlice(budgetModels []models.BudgetModel) []domain.Budget {
	budgets := make([]domain.Budget, len(budgetModels))
	for i, budgetModel := range budgetModels {
		budgets[i] = m.ToDomain(budgetModel)
	}
	return budgets
}
