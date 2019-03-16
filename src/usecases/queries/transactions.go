package queries

import (
	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"

	"github.com/wakuwaku3/account-book.api/src/usecases"
)

type transactions struct {
	repos domains.TransactionsRepository
}

// NewTransactions はインスタンスを生成します
func NewTransactions(
	repos domains.TransactionsRepository,
) usecases.TransactionsQuery {
	return &transactions{
		repos,
	}
}
func (t *transactions) GetTransactions(
	args *usecases.GetTransactionsArgs,
) (*usecases.GetTransactionsResult, error) {
	records, err := t.repos.GetByMonth(&args.SelectedMonth)
	if err != nil {
		return nil, err
	}
	transactions := make([]usecases.GetTransactionResult, len(*records))
	for i, record := range *records {
		r := &record
		transactions[i] = *convertTransaction(r)
	}
	return &usecases.GetTransactionsResult{Transactions: transactions}, nil
}
func (t *transactions) GetTransaction(id *string) (
	*usecases.GetTransactionResult,
	error,
) {
	model, err := t.repos.Get(id)
	if err != nil {
		return nil, err
	}
	return convertTransaction(model), nil
}
func convertTransaction(model *models.Transaction) *usecases.GetTransactionResult {
	return &usecases.GetTransactionResult{
		Amount:        model.Amount,
		Category:      model.Category,
		Date:          model.Date,
		Notes:         model.Notes,
		TransactionID: model.TransactionID,
		Editable:      model.DailyID == nil,
	}
}
