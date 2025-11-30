package category_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/category"
)

func TestCreateUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success with all fields", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewCreateUseCase(mockRepo)

		input := category.CreateInput{
			UserID:      1,
			Name:        "Alimentação",
			Type:        domain.CategoryTypeExpense,
			Icon:        "🍽️",
			Color:       "#FF5733",
			Description: "Gastos com comida",
		}

		expectedCategory := domain.Category{
			ID:          1,
			UserID:      1,
			Name:        "Alimentação",
			Type:        domain.CategoryTypeExpense,
			Icon:        "🍽️",
			Color:       "#FF5733",
			Description: "Gastos com comida",
			IsActive:    true,
		}

		mockRepo.On("Create", ctx, domain.Category{
			UserID:      input.UserID,
			Name:        input.Name,
			Type:        input.Type,
			Icon:        input.Icon,
			Color:       input.Color,
			Description: input.Description,
			IsActive:    true,
		}).Return(expectedCategory, nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - empty name", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewCreateUseCase(mockRepo)

		input := category.CreateInput{
			UserID: 1,
			Name:   "",
			Type:   domain.CategoryTypeExpense,
		}

		result, err := uc.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrCategoryNameRequired, err)
		assert.Equal(t, domain.Category{}, result)
	})

	t.Run("error - invalid category type", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewCreateUseCase(mockRepo)

		input := category.CreateInput{
			UserID: 1,
			Name:   "Test Category",
			Type:   "invalid",
		}

		result, err := uc.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidCategoryType, err)
		assert.Equal(t, domain.Category{}, result)
	})
}

func TestListUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success with categories", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewListUseCase(mockRepo)

		expectedCategories := []domain.Category{
			{ID: 1, UserID: 1, Name: "Salário", Type: domain.CategoryTypeIncome},
			{ID: 2, UserID: 1, Name: "Alimentação", Type: domain.CategoryTypeExpense},
		}

		mockRepo.On("GetByUserID", ctx, uint(1)).Return(expectedCategories, nil)

		result, err := uc.Execute(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByIDUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewGetByIDUseCase(mockRepo)

		expectedCategory := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "Alimentação",
			Type:   domain.CategoryTypeExpense,
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(expectedCategory, nil)

		result, err := uc.Execute(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - unauthorized access", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewGetByIDUseCase(mockRepo)

		categoryFromDB := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "User 1 Category",
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(categoryFromDB, nil)

		result, err := uc.Execute(ctx, 2, 1)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrCategoryNotFound, err)
		assert.Equal(t, domain.Category{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success - update name", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewUpdateUseCase(mockRepo)

		existingCategory := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "Old Name",
			Type:   domain.CategoryTypeExpense,
		}

		newName := "New Name"
		input := category.UpdateInput{
			ID:     1,
			UserID: 1,
			Name:   &newName,
		}

		updatedCategory := existingCategory
		updatedCategory.Name = newName

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingCategory, nil)
		mockRepo.On("Update", ctx, updatedCategory).Return(updatedCategory, nil)

		result, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - empty name", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewUpdateUseCase(mockRepo)

		existingCategory := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "Original Name",
		}

		emptyName := ""
		input := category.UpdateInput{
			ID:     1,
			UserID: 1,
			Name:   &emptyName,
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingCategory, nil)

		result, err := uc.Execute(ctx, input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrCategoryNameRequired, err)
		assert.Equal(t, domain.Category{}, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewDeleteUseCase(mockRepo)

		existingCategory := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "Category to Delete",
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingCategory, nil)
		mockRepo.On("Delete", ctx, uint(1)).Return(nil)

		err := uc.Execute(ctx, 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - unauthorized access", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		uc := category.NewDeleteUseCase(mockRepo)

		existingCategory := domain.Category{
			ID:     1,
			UserID: 1,
			Name:   "User 1 Category",
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(existingCategory, nil)

		err := uc.Execute(ctx, 1, 2)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrCategoryNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
