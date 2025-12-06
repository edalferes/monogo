package gorm

import (
	"context"
	"time"

	gormpkg "gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type gormTransactionRepository struct {
	db     *gormpkg.DB
	mapper mappers.TransactionMapper
}

// NewGormTransactionRepository creates a new GORM-based transaction repository
func NewGormTransactionRepository(db *gormpkg.DB) interfaces.TransactionRepository {
	return &gormTransactionRepository{
		db:     db,
		mapper: mappers.TransactionMapper{},
	}
}

func (r *gormTransactionRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	model := r.mapper.ToModel(transaction)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Transaction{}, err
	}

	// Reload with relationships
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		First(&model, model.ID).Error; err != nil {
		return domain.Transaction{}, err
	}

	return r.mapper.ToDomain(model), nil
}

func (r *gormTransactionRepository) GetByID(ctx context.Context, id uint) (domain.Transaction, error) {
	var model models.TransactionModel
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		First(&model, id).Error; err != nil {
		return domain.Transaction{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormTransactionRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) GetByUserIDPaginated(ctx context.Context, userID uint, limit, offset int) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.TransactionModel{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *gormTransactionRepository) GetByAccountID(ctx context.Context, accountID uint) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	// Get transactions where account is either source OR destination
	if err := r.db.WithContext(ctx).
		Where("account_id = ? OR destination_account_id = ?", accountID, accountID).
		Order("date DESC").
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("date DESC").Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("date DESC").
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) GetByType(ctx context.Context, userID uint, transactionType domain.TransactionType) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, string(transactionType)).
		Order("date DESC").
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) Update(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	model := r.mapper.ToModel(transaction)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Transaction{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormTransactionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TransactionModel{}, id).Error
}

func (r *gormTransactionRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.TransactionModel{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *gormTransactionRepository) GetByUserIDPaginatedWithFilters(ctx context.Context, userID uint, limit, offset int, startDate, endDate *time.Time) ([]domain.Transaction, error) {
	query := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	var transactionModels []models.TransactionModel
	if err := query.
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactionModels).Error; err != nil {
		return nil, err
	}

	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) CountByUserIDWithFilters(ctx context.Context, userID uint, startDate, endDate *time.Time) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&models.TransactionModel{}).
		Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
