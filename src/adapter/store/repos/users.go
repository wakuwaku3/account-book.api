package repos

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/adapter/store"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
)

type users struct {
	provider       store.Provider
	claimsProvider core.ClaimsProvider
}

// NewUsers はインスタンスを生成します
func NewUsers(provider store.Provider, claimsProvider core.ClaimsProvider) application.UsersRepository {
	return &users{provider: provider, claimsProvider: claimsProvider}
}

func (t *users) GetByAuth() (*models.User, error) {
	return t.Get(t.claimsProvider.GetUserID())
}
func (t *users) Get(userID *string) (*models.User, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := client.Collection("users").Doc(*userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var model models.User
	doc.DataTo(&model)
	return &model, nil
}
