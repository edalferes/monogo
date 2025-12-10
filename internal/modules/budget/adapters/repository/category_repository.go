package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
)

type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new GORM-based category repository
func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (domain.Category, error) {
	var category domain.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetByType(ctx context.Context, userID uint, categoryType domain.CategoryType) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, string(categoryType)).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	if err := r.db.WithContext(ctx).Save(&category).Error; err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Category{}, id).Error
}

func (r *CategoryRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
