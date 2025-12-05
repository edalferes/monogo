package mappers

import (
	"github.com/lib/pq"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// TransactionMapper converts between domain.Transaction and models.TransactionModel
type TransactionMapper struct{}

// ToModel converts domain.Transaction to models.TransactionModel
func (m TransactionMapper) ToModel(transaction domain.Transaction) models.TransactionModel {
	return models.TransactionModel{
		ID:                   transaction.ID,
		UserID:               transaction.UserID,
		AccountID:            transaction.AccountID,
		CategoryID:           transaction.CategoryID,
		Type:                 string(transaction.Type),
		Amount:               transaction.Amount,
		Description:          transaction.Description,
		Date:                 transaction.Date,
		Month:                transaction.Month,
		Status:               string(transaction.Status),
		Tags:                 pq.StringArray(transaction.Tags),
		Attachments:          pq.StringArray(transaction.Attachments),
		IsRecurring:          transaction.IsRecurring,
		RecurrenceRule:       transaction.RecurrenceRule,
		RecurrenceEnd:        transaction.RecurrenceEnd,
		ParentID:             transaction.ParentID,
		DestinationAccountID: transaction.DestinationAccountID,
		TransferFee:          transaction.TransferFee,
		CreatedAt:            transaction.CreatedAt,
		UpdatedAt:            transaction.UpdatedAt,
	}
}

// ToDomain converts models.TransactionModel to domain.Transaction
func (m TransactionMapper) ToDomain(transactionModel models.TransactionModel) domain.Transaction {
	accountMapper := AccountMapper{}
	categoryMapper := CategoryMapper{}

	var account *domain.Account
	if transactionModel.Account.ID != 0 {
		mappedAccount := accountMapper.ToDomain(transactionModel.Account)
		account = &mappedAccount
	}

	var category *domain.Category
	if transactionModel.Category.ID != 0 {
		mappedCategory := categoryMapper.ToDomain(transactionModel.Category)
		category = &mappedCategory
	}

	return domain.Transaction{
		ID:                   transactionModel.ID,
		UserID:               transactionModel.UserID,
		AccountID:            transactionModel.AccountID,
		CategoryID:           transactionModel.CategoryID,
		Type:                 domain.TransactionType(transactionModel.Type),
		Amount:               transactionModel.Amount,
		Description:          transactionModel.Description,
		Date:                 transactionModel.Date,
		Month:                transactionModel.Month,
		Status:               domain.TransactionStatus(transactionModel.Status),
		Tags:                 []string(transactionModel.Tags),
		Attachments:          []string(transactionModel.Attachments),
		IsRecurring:          transactionModel.IsRecurring,
		RecurrenceRule:       transactionModel.RecurrenceRule,
		RecurrenceEnd:        transactionModel.RecurrenceEnd,
		ParentID:             transactionModel.ParentID,
		DestinationAccountID: transactionModel.DestinationAccountID,
		TransferFee:          transactionModel.TransferFee,
		CreatedAt:            transactionModel.CreatedAt,
		UpdatedAt:            transactionModel.UpdatedAt,
		Account:              account,
		Category:             category,
	}
}

// ToDomainSlice converts []models.TransactionModel to []domain.Transaction
func (m TransactionMapper) ToDomainSlice(transactionModels []models.TransactionModel) []domain.Transaction {
	transactions := make([]domain.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		transactions[i] = m.ToDomain(transactionModel)
	}
	return transactions
}
