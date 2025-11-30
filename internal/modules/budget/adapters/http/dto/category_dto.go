package dto

import (
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=income expense"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// CategoryResponse represents a category in API responses
type CategoryResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToCategoryResponse converts domain.Category to CategoryResponse
func ToCategoryResponse(category domain.Category) CategoryResponse {
	return CategoryResponse{
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

// ToCategoryResponseList converts []domain.Category to []CategoryResponse
func ToCategoryResponseList(categories []domain.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = ToCategoryResponse(category)
	}
	return responses
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=100"`
	Type        *string `json:"type" validate:"omitempty,oneof=income expense"`
	Icon        *string `json:"icon" validate:"omitempty,max=50"`
	Color       *string `json:"color" validate:"omitempty,max=7"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}
