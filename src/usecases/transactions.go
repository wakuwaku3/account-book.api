package usecases

import (
	"errors"
	"strings"
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	transactions struct {
		query   TransactionsQuery
		service services.Transactions
	}
	// Transactions is TransactionsUseCases
	Transactions interface {
		GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error, error)
		GetTransaction(id *string) (*GetTransactionResult, error, error)
		Create(args *TransactionArgs) (*CreateTransactionResult, error, error)
		Update(id *string, args *TransactionArgs) (error, error)
		Delete(id *string) (error, error)
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
func (t *transactions) GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error, error) {
	info, err := t.query.GetTransactions(args)
	if err != nil {
		return nil, nil, err
	}
	return info, nil, nil
}
func (t *transactions) GetTransaction(id *string) (*GetTransactionResult, error, error) {
	info, err := t.query.GetTransaction(id)
	if err != nil {
		return nil, nil, err
	}
	return info, nil, nil
}
func (t *transactions) Create(args *TransactionArgs) (*CreateTransactionResult, error, error) {
	if err := args.valid(); err != nil {
		return nil, err, nil
	}
	res, cErr, err := t.service.Create(args.convert())
	if err != nil {
		return nil, nil, err
	}
	if cErr != nil {
		return nil, cErr, nil
	}
	return &CreateTransactionResult{
		TransactionID: res.TransactionID,
	}, nil, nil
}
func (t *TransactionArgs) valid() error {
	array := make([]string, 0)
	if t.Amount == nil {
		array = append(array, "金額が入力されていません。")
	}
	if t.Category == nil {
		array = append(array, "カテゴリが入力されていません。")
	}
	if len(array) > 0 {
		return errors.New(strings.Join(array, "\n"))
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
func (t *transactions) Update(id *string, args *TransactionArgs) (error, error) {
	if err := args.valid(); err != nil {
		return err, nil
	}
	return t.service.Update(id, args.convert())
}
func (t *transactions) Delete(id *string) (error, error) {
	return t.service.Delete(id)
}
