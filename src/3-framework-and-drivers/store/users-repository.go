package store

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/0-enterprise-business-rules/models"
	"github.com/wakuwaku3/account-book.api/src/1-application-business-rules/usecases"
	"google.golang.org/api/iterator"
)

type usersRepository struct {
	provider Provider
}

// CreateUsersRepository はインスタンスを生成します
func CreateUsersRepository(provider Provider) usecases.UsersRepository {
	return &usersRepository{provider: provider}
}

func (usersRepository *usersRepository) Get() (*[]models.User, error) {
	client := usersRepository.provider.GetClient()
	ctx := context.Background()
	iter := client.Collection("users").Documents(ctx)
	users := make([]models.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user models.User
		doc.DataTo(&user)
		users = append(users, user)
	}
	return &users, nil
}
