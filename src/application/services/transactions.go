package services

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"

	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	accountbook "github.com/wakuwaku3/account-book.api/src/enterprise/domains/accountBook"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
)

type (
	transactions struct {
		repos              application.TransactionsRepository
		clock              core.Clock
		assetsChangedEvent accountbook.AssetsChangedEvent
	}
	// Transactions is TransactionsService
	Transactions interface {
		Create(args *TransactionArgs) (*CreateTransactionResult, error)
		Update(id *string, args *TransactionArgs) error
		Delete(id *string) error
	}
	// TransactionArgs は引数です
	TransactionArgs struct {
		Amount   int
		Category int
		Notes    *string
	}
	// CreateTransactionResult は結果です
	CreateTransactionResult struct {
		TransactionID string
	}
)

// NewTransactions is create instance
func NewTransactions(
	repos application.TransactionsRepository,
	clock core.Clock,
	assetsChangedEvent accountbook.AssetsChangedEvent,
) Transactions {
	return &transactions{repos, clock, assetsChangedEvent}
}
func (t *transactions) Create(args *TransactionArgs) (*CreateTransactionResult, error) {
	id, err := t.repos.Create(args.convert(t.clock.Now()))
	if err != nil {
		return nil, err
	}
	t.assetsChangedEvent.Trigger()
	return &CreateTransactionResult{TransactionID: *id}, nil
}
func (t *TransactionArgs) convert(now time.Time) *models.Transaction {
	return &models.Transaction{
		Amount:   t.Amount,
		Category: t.Category,
		Notes:    t.Notes,
		Date:     now,
	}
}
func (t *transactions) Update(id *string, args *TransactionArgs) error {
	model, err := t.repos.Get(id)
	if err != nil {
		return err
	}
	if model.DailyID != nil {
		return core.NewError(application.ClosedTransaction)
	}

	model.Amount = args.Amount
	model.Category = args.Category
	model.Notes = args.Notes

	if err := t.repos.Update(id, model); err != nil {
		return err
	}
	t.assetsChangedEvent.Trigger()
	return nil
}
func (t *transactions) Delete(id *string) error {
	model, err := t.repos.Get(id)
	if err != nil {
		return err
	}
	if model.DailyID != nil {
		return core.NewError(application.ClosedTransaction)
	}

	if err := t.repos.Delete(id); err != nil {
		return err
	}
	t.assetsChangedEvent.Trigger()
	return nil
}
