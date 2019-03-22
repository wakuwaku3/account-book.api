package repos

import (
	"context"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"google.golang.org/api/iterator"
)

type accounts struct {
	provider store.Provider
	clock    cmn.Clock
}

// NewAccounts はインスタンスを生成します
func NewAccounts(provider store.Provider, clock cmn.Clock) domains.AccountsRepository {
	return &accounts{provider, clock}
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
func (t *accounts) CreatePasswordResetToken(model *models.PasswordResetToken) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("passwordResetTokens")
	ref, _, err := passwordRestTokensRef.Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *accounts) CleanUp() error {
	now := t.clock.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("passwordResetTokens")

	iter := passwordRestTokensRef.Where("expires", "<=", now).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		batch.Delete(doc.Ref)
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (t *accounts) CleanUpByEmail(email *string) error {
	now := t.clock.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := client.Collection("passwordResetTokens")

	iter := passwordRestTokensRef.Where("email", "==", *email).Where("expires", "<=", now).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		batch.Delete(doc.Ref)
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (t *accounts) GetPasswordResetToken(passwordResetToken *string) (*models.PasswordResetToken, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := client.Collection("passwordResetTokens").Doc(*passwordResetToken).Get(ctx)
	if err != nil {
		return nil, err
	}
	var model models.PasswordResetToken
	doc.DataTo(&model)
	return &model, nil
}
func (t *accounts) SetPassword(email *string, hashedPassword *string) error {
	client := t.provider.GetClient()
	ctx := context.Background()
	ref := client.Collection("accounts").Doc(*email)
	doc, err := ref.Get(ctx)
	if err != nil {
		return err
	}
	var model models.Account
	doc.DataTo(&model)
	model.HashedPassword = *hashedPassword
	_, err = ref.Set(ctx, &model)
	return err
}
