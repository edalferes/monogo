package category

import "github.com/edalferes/monetics/internal/modules/budget/domain"

// isValidCategoryType checks if the category type is valid
func isValidCategoryType(categoryType domain.CategoryType) bool {
	switch categoryType {
	case domain.CategoryTypeIncome,
		domain.CategoryTypeExpense:
		return true
	default:
		return false
	}
}
