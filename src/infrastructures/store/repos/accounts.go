package repos

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
)

type accounts struct {
	provider store.Provider
}

// NewAccounts はインスタンスを生成します
func NewAccounts(provider store.Provider) domains.AccountsRepository {
	return &accounts{provider: provider}
}

func (t *accounts) Get(email *string) (*models.Account, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := client.Collection("accounts").Doc(*email).Get(ctx)
	if err != nil {
		return nil, err
	}
	var model models.Account
	doc.DataTo(&model)
	return &model, nil
}
