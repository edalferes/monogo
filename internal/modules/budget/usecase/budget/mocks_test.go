package budget_test

import (
"context"

"github.com/edalferes/monetics/internal/modules/budget/domain"
"github.com/stretchr/testify/mock"
)

type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) Create(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	args := m.Called(ctx, budget)
	return args.Get(0).(domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetByID(ctx context.Context, id uint) (domain.Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Budget{}, args.Error(1)
	}
	return args.Get(0).(domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Budget, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Budget, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetActive(ctx context.Context, userID uint) ([]domain.Budget, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Update(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	args := m.Called(ctx, budget)
	return args.Get(0).(domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBudgetRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockBudgetRepository) UpdateSpent(ctx context.Context, budgetID uint, spent float64) error {
	args := m.Called(ctx, budgetID, spent)
	return args.Error(0)
}

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uint) (domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Category{}, args.Error(1)
	}
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Category, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByType(ctx context.Context, userID uint, categoryType domain.CategoryType) ([]domain.Category, error) {
	args := m.Called(ctx, userID, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
