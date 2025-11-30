package transaction_test

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uint) (domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Transaction{}, args.Error(1)
	}
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(ctx context.Context, accountID uint) ([]domain.Transaction, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Transaction, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]domain.Transaction, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByType(ctx context.Context, userID uint, transactionType domain.TransactionType) ([]domain.Transaction, error) {
	args := m.Called(ctx, userID, transactionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id uint) (domain.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.Account{}, args.Error(1)
	}
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Account, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(domain.Account), args.Error(1)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAccountRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
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
