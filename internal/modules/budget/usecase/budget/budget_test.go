package budget_test

import (
	"context"
	"testing"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/budget"
	"github.com/stretchr/testify/assert"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockCategoryRepo := new(MockCategoryRepository)
		uc := budget.NewCreateUseCase(mockBudgetRepo, mockCategoryRepo)

		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

		input := budget.CreateInput{
			UserID:     1,
			CategoryID: 1,
			Name:       "Janeiro 2025",
			Amount:     2000.00,
			Period:     domain.BudgetPeriodMonthly,
			StartDate:  startDate,
			EndDate:    endDate,
		}

		category := domain.Category{ID: 1, UserID: 1, Name: "Alimentação"}
		mockCategoryRepo.On("GetByID", ctx, uint(1)).Return(category, nil)

		expectedBudget := domain.Budget{
			ID:         1,
			UserID:     1,
			CategoryID: 1,
			Name:       "Janeiro 2025",
			Amount:     2000.00,
			Period:     domain.BudgetPeriodMonthly,
			StartDate:  startDate,
			EndDate:    endDate,
			IsActive:   true,
			Spent:      0,
		}

		mockBudgetRepo.On("Create", ctx, domain.Budget{
			UserID:     1,
			CategoryID: 1,
			Name:       "Janeiro 2025",
			Amount:     2000.00,
			Period:     domain.BudgetPeriodMonthly,
			StartDate:  startDate,
			EndDate:    endDate,
			IsActive:   true,
			Spent:      0,
		}).Return(expectedBudget, nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedBudget.ID, result.ID)
		mockBudgetRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("error - invalid amount", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockCategoryRepo := new(MockCategoryRepository)
		uc := budget.NewCreateUseCase(mockBudgetRepo, mockCategoryRepo)

		input := budget.CreateInput{
			UserID:     1,
			CategoryID: 1,
			Name:       "Test",
			Amount:     -100,
			Period:     domain.BudgetPeriodMonthly,
		}

		_, err := uc.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidBudgetAmount, err)
	})
}

func TestListUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockBudgetRepository)
		uc := budget.NewListUseCase(mockRepo)

		expected := []domain.Budget{
			{ID: 1, UserID: 1, Name: "Janeiro"},
			{ID: 2, UserID: 1, Name: "Fevereiro"},
		}

		mockRepo.On("GetByUserID", ctx, uint(1)).Return(expected, nil)

		result, err := uc.Execute(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByIDUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockBudgetRepository)
		uc := budget.NewGetByIDUseCase(mockRepo)

		expected := domain.Budget{ID: 1, UserID: 1, Name: "Janeiro"}
		mockRepo.On("GetByID", ctx, uint(1)).Return(expected, nil)

		result, err := uc.Execute(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - unauthorized", func(t *testing.T) {
		mockRepo := new(MockBudgetRepository)
		uc := budget.NewGetByIDUseCase(mockRepo)

		budgetFromDB := domain.Budget{ID: 1, UserID: 1}
		mockRepo.On("GetByID", ctx, uint(1)).Return(budgetFromDB, nil)

		_, err := uc.Execute(ctx, 2, 1)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrBudgetNotFound, err)
	})
}

func TestDeleteUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockBudgetRepository)
		uc := budget.NewDeleteUseCase(mockRepo)

		budgetFromDB := domain.Budget{ID: 1, UserID: 1}
		mockRepo.On("GetByID", ctx, uint(1)).Return(budgetFromDB, nil)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := uc.Execute(ctx, 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
