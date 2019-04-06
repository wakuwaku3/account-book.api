package repos

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"

	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"google.golang.org/api/iterator"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
)

type (
	transactions struct {
		provider       store.Provider
		clock          cmn.Clock
		claimsProvider domains.ClaimsProvider
	}
)

// NewTransactions はインスタンスを生成します
func NewTransactions(
	provider store.Provider,
	clock cmn.Clock,
	claimsProvider domains.ClaimsProvider,
) domains.TransactionsRepository {
	return &transactions{provider, clock, claimsProvider}
}
func (t *transactions) transactionsRef(client *firestore.Client) *firestore.CollectionRef {
	userID := t.claimsProvider.GetUserID()
	return client.Collection("users").Doc(*userID).Collection("transactions")
}
func (t *transactions) Get(id *string) (*models.Transaction, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.transactionsRef(client).Doc(*id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var transaction models.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}
func (t *transactions) GetByMonth(month *time.Time) (*[]models.Transaction, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	start := t.clock.GetMonthStartDay(month)
	end := start.AddDate(0, 1, 0)

	transactions := make([]models.Transaction, 0)
	iter := t.transactionsRef(client).Where("date", ">=", start).Where("date", "<", end).OrderBy("date", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var transaction models.Transaction
		if err := doc.DataTo(&transaction); err != nil {
			return nil, err
		}
		transaction.TransactionID = doc.Ref.ID
		transactions = append(transactions, transaction)
	}
	return &transactions, nil
}
func (t *transactions) Create(model *models.Transaction) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	ref, _, err := t.transactionsRef(client).Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *transactions) Update(id *string, model *models.Transaction) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.transactionsRef(client).Doc(*id).Set(ctx, model)
	if err != nil {
		return err
	}
	return nil
}
func (t *transactions) Delete(id *string) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	_, err := t.transactionsRef(client).Doc(*id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
