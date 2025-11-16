package dto

import (
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// CreateTransactionRequest represents the request to create a transaction
type CreateTransactionRequest struct {
	AccountID            uint       `json:"account_id" validate:"required"`
	CategoryID           uint       `json:"category_id" validate:"required"`
	Type                 string     `json:"type" validate:"required,oneof=income expense transfer"`
	Amount               float64    `json:"amount" validate:"required,gt=0"`
	Description          string     `json:"description" validate:"required"`
	Date                 time.Time  `json:"date" validate:"required"`
	Tags                 []string   `json:"tags"`
	Attachments          []string   `json:"attachments"`
	IsRecurring          bool       `json:"is_recurring"`
	RecurrenceRule       string     `json:"recurrence_rule"`
	RecurrenceEnd        *time.Time `json:"recurrence_end"`
	DestinationAccountID *uint      `json:"destination_account_id"`
	TransferFee          *float64   `json:"transfer_fee"`
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	ID                   uint       `json:"id"`
	UserID               uint       `json:"user_id"`
	AccountID            uint       `json:"account_id"`
	CategoryID           uint       `json:"category_id"`
	Type                 string     `json:"type"`
	Amount               float64    `json:"amount"`
	Description          string     `json:"description"`
	Date                 time.Time  `json:"date"`
	Status               string     `json:"status"`
	Tags                 []string   `json:"tags,omitempty"`
	Attachments          []string   `json:"attachments,omitempty"`
	IsRecurring          bool       `json:"is_recurring"`
	RecurrenceRule       string     `json:"recurrence_rule,omitempty"`
	RecurrenceEnd        *time.Time `json:"recurrence_end,omitempty"`
	ParentID             *uint      `json:"parent_id,omitempty"`
	DestinationAccountID *uint      `json:"destination_account_id,omitempty"`
	TransferFee          *float64   `json:"transfer_fee,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// ToTransactionResponse converts domain.Transaction to TransactionResponse
func ToTransactionResponse(transaction domain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:                   transaction.ID,
		UserID:               transaction.UserID,
		AccountID:            transaction.AccountID,
		CategoryID:           transaction.CategoryID,
		Type:                 string(transaction.Type),
		Amount:               transaction.Amount,
		Description:          transaction.Description,
		Date:                 transaction.Date,
		Status:               string(transaction.Status),
		Tags:                 transaction.Tags,
		Attachments:          transaction.Attachments,
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

// ToTransactionResponseList converts []domain.Transaction to []TransactionResponse
func ToTransactionResponseList(transactions []domain.Transaction) []TransactionResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = ToTransactionResponse(transaction)
	}
	return responses
}

// UpdateTransactionRequest represents the request to update a transaction
type UpdateTransactionRequest struct {
	CategoryID  *uint      `json:"category_id"`
	Amount      *float64   `json:"amount"`
	Description string     `json:"description"`
	Date        *time.Time `json:"date"`
	Tags        []string   `json:"tags"`
}
