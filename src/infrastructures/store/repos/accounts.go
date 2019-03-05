package repos

import (
	"context"
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"google.golang.org/api/iterator"
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
func (t *accounts) CreatePasswordResetToken(email *string, expires *time.Time) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("password-reset-tokens")
	ref, _, err := passwordRestTokensRef.Add(ctx, map[string]interface{}{
		"expires": *expires,
	})
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *accounts) CleanUp() error {
	now := time.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("password-reset-tokens")

	iter := passwordRestTokensRef.Where("expires", "<=", now).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var old models.PasswordResetToken
		doc.DataTo(&old)
		if old.Expires.Equal(now) || old.Expires.Before(now) {
			batch.Delete(doc.Ref)
		}
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (t *accounts) CleanUpByEmail(email string) error {
	now := time.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("password-reset-tokens")

	iter := passwordRestTokensRef.Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var old models.PasswordResetToken
		doc.DataTo(&old)
		if old.Expires.Equal(now) || old.Expires.Before(now) {
			batch.Delete(doc.Ref)
		}
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}
