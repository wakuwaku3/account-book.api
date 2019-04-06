package repos

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"github.com/wakuwaku3/account-book.api/src/drivers/store"
)

type users struct {
	provider store.Provider
}

// NewUsers はインスタンスを生成します
func NewUsers(provider store.Provider) application.UsersRepository {
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
