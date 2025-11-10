package mappers

import (
	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// CategoryMapper converts between domain.Category and models.CategoryModel
type CategoryMapper struct{}

// ToModel converts domain.Category to models.CategoryModel
func (m CategoryMapper) ToModel(category domain.Category) models.CategoryModel {
	return models.CategoryModel{
		ID:          category.ID,
		UserID:      category.UserID,
		Name:        category.Name,
		Type:        string(category.Type),
		Icon:        category.Icon,
		Color:       category.Color,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToDomain converts models.CategoryModel to domain.Category
func (m CategoryMapper) ToDomain(categoryModel models.CategoryModel) domain.Category {
	return domain.Category{
		ID:          categoryModel.ID,
		UserID:      categoryModel.UserID,
		Name:        categoryModel.Name,
		Type:        domain.CategoryType(categoryModel.Type),
		Icon:        categoryModel.Icon,
		Color:       categoryModel.Color,
		Description: categoryModel.Description,
		CreatedAt:   categoryModel.CreatedAt,
		UpdatedAt:   categoryModel.UpdatedAt,
	}
}

// ToDomainSlice converts []models.CategoryModel to []domain.Category
func (m CategoryMapper) ToDomainSlice(categoryModels []models.CategoryModel) []domain.Category {
	categories := make([]domain.Category, len(categoryModels))
	for i, categoryModel := range categoryModels {
		categories[i] = m.ToDomain(categoryModel)
	}
	return categories
}
