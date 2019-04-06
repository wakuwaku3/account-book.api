package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"
	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	transactions struct {
		query   TransactionsQuery
		service services.Transactions
	}
	// Transactions is TransactionsUseCases
	Transactions interface {
		GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error)
		GetTransaction(id *string) (*GetTransactionResult, error)
		Create(args *TransactionArgs) (*CreateTransactionResult, error)
		Update(id *string, args *TransactionArgs) error
		Delete(id *string) error
	}
	// GetTransactionsArgs は引数です
	GetTransactionsArgs struct {
		SelectedMonth time.Time
	}
	// GetTransactionsResult は結果です
	GetTransactionsResult struct {
		Transactions []GetTransactionResult
	}
	// GetTransactionResult は結果です
	GetTransactionResult struct {
		TransactionID string
		Amount        int
		Category      int
		Date          time.Time
		Notes         *string
		Editable      bool
	}
	// TransactionArgs は引数です
	TransactionArgs struct {
		Amount   *int
		Category *int
		Notes    *string
	}
	// CreateTransactionResult は結果です
	CreateTransactionResult struct {
		TransactionID string
	}
)

// NewTransactions is create instance
func NewTransactions(
	query TransactionsQuery,
	service services.Transactions,
) Transactions {
	return &transactions{
		query,
		service,
	}
}
func (t *transactions) GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error) {
	info, err := t.query.GetTransactions(args)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *transactions) GetTransaction(id *string) (*GetTransactionResult, error) {
	info, err := t.query.GetTransaction(id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (t *transactions) Create(args *TransactionArgs) (*CreateTransactionResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	res, err := t.service.Create(args.convert())
	if err != nil {
		return nil, err
	}
	return &CreateTransactionResult{
		TransactionID: res.TransactionID,
	}, nil
}
func (t *TransactionArgs) valid() error {
	err := apperrors.NewClientError()
	if t.Amount == nil {
		err.Append(apperrors.RequiredAmount)
	}
	if t.Category == nil {
		err.Append(apperrors.RequiredCategory)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *TransactionArgs) convert() *services.TransactionArgs {
	return &services.TransactionArgs{
		Amount:   *t.Amount,
		Category: *t.Category,
		Notes:    t.Notes,
	}
}
func (t *transactions) Update(id *string, args *TransactionArgs) error {
	if err := args.valid(); err != nil {
		return err
	}
	return t.service.Update(id, args.convert())
}
func (t *transactions) Delete(id *string) error {
	return t.service.Delete(id)
}
