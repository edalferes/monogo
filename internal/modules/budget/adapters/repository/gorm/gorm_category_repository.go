package gorm

import (
	"context"

	gormpkg "gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type gormCategoryRepository struct {
	db     *gormpkg.DB
	mapper mappers.CategoryMapper
}

// NewGormCategoryRepository creates a new GORM-based category repository
func NewGormCategoryRepository(db *gormpkg.DB) interfaces.CategoryRepository {
	return &gormCategoryRepository{
		db:     db,
		mapper: mappers.CategoryMapper{},
	}
}

func (r *gormCategoryRepository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	model := r.mapper.ToModel(category)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Category{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormCategoryRepository) GetByID(ctx context.Context, id uint) (domain.Category, error) {
	var model models.CategoryModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Category{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormCategoryRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Category, error) {
	var categoryModels []models.CategoryModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categoryModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(categoryModels), nil
}

func (r *gormCategoryRepository) GetByType(ctx context.Context, userID uint, categoryType domain.CategoryType) ([]domain.Category, error) {
	var categoryModels []models.CategoryModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, string(categoryType)).Find(&categoryModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(categoryModels), nil
}

func (r *gormCategoryRepository) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	model := r.mapper.ToModel(category)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Category{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormCategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CategoryModel{}, id).Error
}

func (r *gormCategoryRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CategoryModel{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
