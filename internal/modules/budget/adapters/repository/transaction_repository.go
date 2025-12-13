package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new GORM-based transaction repository
func NewTransactionRepository(db *gorm.DB) interfaces.TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(&transaction).Error; err != nil {
		return domain.Transaction{}, err
	}

	// Reload with relationships
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		First(&transaction, transaction.ID).Error; err != nil {
		return domain.Transaction{}, err
	}

	return transaction, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uint) (domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		First(&transaction, id).Error; err != nil {
		return domain.Transaction{}, err
	}
	return transaction, nil
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByUserIDPaginated(ctx context.Context, userID uint, limit, offset int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	// Get transactions where account is either source OR destination
	if err := r.db.WithContext(ctx).
		Where("account_id = ? OR destination_account_id = ?", accountID, accountID).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetByType(ctx context.Context, userID uint, transactionType domain.TransactionType) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, string(transactionType)).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	if err := r.db.WithContext(ctx).Save(&transaction).Error; err != nil {
		return domain.Transaction{}, err
	}
	return transaction, nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Transaction{}, id).Error
}

func (r *TransactionRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Transaction{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *TransactionRepository) GetByUserIDPaginatedWithFilters(ctx context.Context, userID uint, limit, offset int, startDate, endDate *time.Time) ([]domain.Transaction, error) {
	query := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	var transactions []domain.Transaction
	if err := query.
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) CountByUserIDWithFilters(ctx context.Context, userID uint, startDate, endDate *time.Time) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetByUserIDPaginatedWithAllFilters retrieves paginated transactions with all filters
func (r *TransactionRepository) GetByUserIDPaginatedWithAllFilters(
	ctx context.Context,
	userID uint,
	limit, offset int,
	txType *domain.TransactionType,
	accountID, categoryID *uint,
	startDate, endDate *time.Time,
) ([]domain.Transaction, error) {
	query := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID)

	if txType != nil {
		query = query.Where("type = ?", *txType)
	}
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	var transactions []domain.Transaction
	if err := query.
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// CountByUserIDWithAllFilters counts transactions with all filters
func (r *TransactionRepository) CountByUserIDWithAllFilters(
	ctx context.Context,
	userID uint,
	txType *domain.TransactionType,
	accountID, categoryID *uint,
	startDate, endDate *time.Time,
) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&domain.Transaction{}).
		Where("user_id = ?", userID)

	if txType != nil {
		query = query.Where("type = ?", *txType)
	}
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
