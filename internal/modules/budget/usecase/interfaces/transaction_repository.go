package interfaces

import (
	"context"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// TransactionRepository defines the contract for transaction persistence operations
type TransactionRepository interface {
	Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	GetByID(ctx context.Context, id uint) (domain.Transaction, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uint) ([]domain.Transaction, error)
	GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Transaction, error)
	GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]domain.Transaction, error)
	GetByType(ctx context.Context, userID uint, transactionType domain.TransactionType) ([]domain.Transaction, error)
	Update(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error)
	Delete(ctx context.Context, id uint) error
	ExistsByID(ctx context.Context, id uint) (bool, error)
}
