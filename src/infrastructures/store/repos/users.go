package repos

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
)

type users struct {
	provider store.Provider
}

// NewUsers はインスタンスを生成します
func NewUsers(provider store.Provider) domains.UsersRepository {
	return &users{provider: provider}
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
