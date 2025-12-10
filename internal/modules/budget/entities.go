package budget

import (
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// Entities returns all database entities for the budget module
func Entities() []interface{} {
	return []interface{}{
		&domain.Account{},
		&domain.Category{},
		&domain.Transaction{},
		&domain.Budget{},
	}
}
