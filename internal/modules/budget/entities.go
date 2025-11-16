package budget

import (
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
)

// Entities returns all database entities for the budget module
func Entities() []interface{} {
	return []interface{}{
		&models.AccountModel{},
		&models.CategoryModel{},
		&models.TransactionModel{},
		&models.BudgetModel{},
	}
}
