package transaction

import "github.com/edalferes/monetics/internal/modules/budget/domain"

func isValidTransactionType(txType domain.TransactionType) bool {
	switch txType {
	case domain.TransactionTypeIncome,
		domain.TransactionTypeExpense,
		domain.TransactionTypeTransfer:
		return true
	default:
		return false
	}
}
