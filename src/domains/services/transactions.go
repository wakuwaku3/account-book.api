package services

import (
	"errors"
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	transactions struct {
		repos domains.TransactionsRepository
		clock cmn.Clock
	}
	// Transactions is TransactionsService
	Transactions interface {
		Create(args *TransactionArgs) (*CreateTransactionResult, error, error)
		Update(id *string, args *TransactionArgs) (error, error)
		Delete(id *string) (error, error)
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
func NewTransactions(repos domains.TransactionsRepository, clock cmn.Clock) Transactions {
	return &transactions{repos, clock}
}
func (t *transactions) Create(args *TransactionArgs) (*CreateTransactionResult, error, error) {
	id, err := t.repos.Create(args.convert(t.clock.Now()))
	if err != nil {
		return nil, nil, err
	}
	return &CreateTransactionResult{TransactionID: *id}, nil, nil
}
func (t *TransactionArgs) convert(now time.Time) *models.Transaction {
	return &models.Transaction{
		Amount:   t.Amount,
		Category: t.Category,
		Notes:    t.Notes,
		Date:     now,
	}
}
func (t *transactions) Update(id *string, args *TransactionArgs) (error, error) {
	model, err := t.repos.Get(id)
	if err != nil {
		return nil, err
	}
	if model.DailyID != nil {
		return errors.New("締め処理後の取引は変更できません。"), nil
	}

	model.Amount = args.Amount
	model.Category = args.Category
	model.Notes = args.Notes

	if err := t.repos.Update(id, model); err != nil {
		return nil, err
	}
	return nil, nil
}
func (t *transactions) Delete(id *string) (error, error) {
	model, err := t.repos.Get(id)
	if err != nil {
		return nil, err
	}
	if model.DailyID != nil {
		return errors.New("締め処理後の取引は削除できません。"), nil
	}

	if err := t.repos.Delete(id); err != nil {
		return nil, err
	}
	return nil, nil
}
