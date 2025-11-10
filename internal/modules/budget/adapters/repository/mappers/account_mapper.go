package mappers

import (
	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// AccountMapper converts between domain.Account and models.AccountModel
type AccountMapper struct{}

// ToModel converts domain.Account to models.AccountModel
func (m AccountMapper) ToModel(account domain.Account) models.AccountModel {
	return models.AccountModel{
		ID:          account.ID,
		UserID:      account.UserID,
		Name:        account.Name,
		Type:        string(account.Type),
		Balance:     account.Balance,
		Currency:    account.Currency,
		Description: account.Description,
		IsActive:    account.IsActive,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}

// ToDomain converts models.AccountModel to domain.Account
func (m AccountMapper) ToDomain(accountModel models.AccountModel) domain.Account {
	return domain.Account{
		ID:          accountModel.ID,
		UserID:      accountModel.UserID,
		Name:        accountModel.Name,
		Type:        domain.AccountType(accountModel.Type),
		Balance:     accountModel.Balance,
		Currency:    accountModel.Currency,
		Description: accountModel.Description,
		IsActive:    accountModel.IsActive,
		CreatedAt:   accountModel.CreatedAt,
		UpdatedAt:   accountModel.UpdatedAt,
	}
}

// ToDomainSlice converts []models.AccountModel to []domain.Account
func (m AccountMapper) ToDomainSlice(accountModels []models.AccountModel) []domain.Account {
	accounts := make([]domain.Account, len(accountModels))
	for i, accountModel := range accountModels {
		accounts[i] = m.ToDomain(accountModel)
	}
	return accounts
}
