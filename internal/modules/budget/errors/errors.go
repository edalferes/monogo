package errors

import "errors"

var (
	// Account errors
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInvalidAccountType   = errors.New("invalid account type")
	ErrAccountNameRequired  = errors.New("account name is required")
	ErrInvalidAccountID     = errors.New("invalid account ID")

	// Category errors
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryInUse         = errors.New("category is in use by transactions")
	ErrInvalidCategoryType   = errors.New("invalid category type")
	ErrCategoryNameRequired  = errors.New("category name is required")
	ErrInvalidCategoryID     = errors.New("invalid category ID")

	// Transaction errors
	ErrTransactionNotFound            = errors.New("transaction not found")
	ErrInvalidTransactionType         = errors.New("invalid transaction type")
	ErrInvalidTransactionStatus       = errors.New("invalid transaction status")
	ErrInvalidAmount                  = errors.New("invalid amount")
	ErrInvalidDate                    = errors.New("invalid date")
	ErrTransferSameAccount            = errors.New("cannot transfer to same account")
	ErrTransactionCancelled           = errors.New("transaction is already cancelled")
	ErrTransactionDescriptionRequired = errors.New("transaction description is required")
	ErrTransferRequiresDestination    = errors.New("transfer requires destination account")

	// Budget errors
	ErrBudgetNotFound      = errors.New("budget not found")
	ErrBudgetAlreadyExists = errors.New("budget already exists for this period")
	ErrInvalidBudgetPeriod = errors.New("invalid budget period")
	ErrInvalidBudgetDates  = errors.New("invalid budget dates")
	ErrBudgetOverlap       = errors.New("budget period overlaps with existing budget")
	ErrBudgetNameRequired  = errors.New("budget name is required")
	ErrInvalidBudgetAmount = errors.New("invalid budget amount")
	ErrInvalidDateRange    = errors.New("invalid date range: start date must be before end date")

	// General errors
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrUnauthorizedAccess = errors.New("unauthorized access")
	ErrInvalidData        = errors.New("invalid data")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidUserID      = errors.New("invalid user ID")
)
