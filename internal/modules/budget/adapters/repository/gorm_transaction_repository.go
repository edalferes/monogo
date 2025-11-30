package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type gormTransactionRepository struct {
	db     *gorm.DB
	mapper mappers.TransactionMapper
}

// NewGormTransactionRepository creates a new GORM-based transaction repository
func NewGormTransactionRepository(db *gorm.DB) TransactionRepository {
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
	return r.mapper.ToDomain(model), nil
}

func (r *gormTransactionRepository) GetByID(ctx context.Context, id uint) (domain.Transaction, error) {
	var model models.TransactionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Transaction{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormTransactionRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date DESC").Find(&transactionModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(transactionModels), nil
}

func (r *gormTransactionRepository) GetByAccountID(ctx context.Context, accountID uint) ([]domain.Transaction, error) {
	var transactionModels []models.TransactionModel
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).Order("date DESC").Find(&transactionModels).Error; err != nil {
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
