package gorm

import (
	"context"

	gormpkg "gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

type gormBudgetRepository struct {
	db     *gormpkg.DB
	mapper mappers.BudgetMapper
}

// NewGormBudgetRepository creates a new GORM-based budget repository
func NewGormBudgetRepository(db *gormpkg.DB) repository.BudgetRepository {
	return &gormBudgetRepository{
		db:     db,
		mapper: mappers.BudgetMapper{},
	}
}

func (r *gormBudgetRepository) Create(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	model := r.mapper.ToModel(budget)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Budget{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormBudgetRepository) GetByID(ctx context.Context, id uint) (domain.Budget, error) {
	var model models.BudgetModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Budget{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormBudgetRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Budget, error) {
	var budgetModels []models.BudgetModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&budgetModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(budgetModels), nil
}

func (r *gormBudgetRepository) GetByCategoryID(ctx context.Context, categoryID uint) ([]domain.Budget, error) {
	var budgetModels []models.BudgetModel
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&budgetModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(budgetModels), nil
}

func (r *gormBudgetRepository) GetActive(ctx context.Context, userID uint) ([]domain.Budget, error) {
	var budgetModels []models.BudgetModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&budgetModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(budgetModels), nil
}

func (r *gormBudgetRepository) Update(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	model := r.mapper.ToModel(budget)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Budget{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormBudgetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.BudgetModel{}, id).Error
}

func (r *gormBudgetRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.BudgetModel{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *gormBudgetRepository) UpdateSpent(ctx context.Context, budgetID uint, spent float64) error {
	return r.db.WithContext(ctx).Model(&models.BudgetModel{}).Where("id = ?", budgetID).Update("spent", spent).Error
}
