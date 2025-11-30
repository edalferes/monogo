package account

import "github.com/edalferes/monetics/internal/modules/budget/domain"

// isValidAccountType checks if the account type is valid
func isValidAccountType(accountType domain.AccountType) bool {
	switch accountType {
	case domain.AccountTypeChecking,
		domain.AccountTypeSavings,
		domain.AccountTypeCredit,
		domain.AccountTypeCash,
		domain.AccountTypeInvest:
		return true
	default:
		return false
	}
}
