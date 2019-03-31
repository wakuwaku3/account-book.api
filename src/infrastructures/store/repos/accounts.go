package repos

import (
	"context"

	"cloud.google.com/go/firestore"
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
func (t *accounts) accountsRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection("accounts")
}
func (t *accounts) passwordResetTokensRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection("passwordResetTokens")
}
func (t *accounts) signUpTokensRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection("signUpTokens")
}

func (t *accounts) Get(email *string) (*models.Account, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.accountsRef(client).Doc(*email).Get(ctx)
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
	passwordRestTokensRef := t.passwordResetTokensRef(client)
	ref, _, err := passwordRestTokensRef.Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *accounts) CleanUpPasswordResetToken() error {
	now := t.clock.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := t.passwordResetTokensRef(client)

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
func (t *accounts) CleanUpPasswordResetTokenByEmail(email *string) error {
	now := t.clock.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordRestTokensRef := t.passwordResetTokensRef(client)

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
	doc, err := t.passwordResetTokensRef(client).Doc(*passwordResetToken).Get(ctx)
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
	ref := t.accountsRef(client).Doc(*email)
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
func (t *accounts) CreateSignUpToken(model *models.SignUpToken) (*string, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	signUpTokensRef := t.signUpTokensRef(client)
	ref, _, err := signUpTokensRef.Add(ctx, model)
	if err != nil {
		return nil, err
	}
	return &ref.ID, nil
}
func (t *accounts) CleanUpSignUpToken() error {
	now := t.clock.Now()
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	signUpTokensRef := t.signUpTokensRef(client)

	iter := signUpTokensRef.Where("expires", "<=", now).Documents(ctx)
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
