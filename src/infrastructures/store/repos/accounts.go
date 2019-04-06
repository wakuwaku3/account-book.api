package repos

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/store"
	"google.golang.org/api/iterator"
)

type accounts struct {
	provider store.Provider
	clock    cmn.Clock
}

// NewAccounts はインスタンスを生成します
func NewAccounts(provider store.Provider, clock cmn.Clock) application.AccountsRepository {
	return &accounts{provider, clock}
}
func (t *accounts) usersRef(client *firestore.Client) *firestore.CollectionRef {
	return client.Collection("users")
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
	passwordResetTokensRef := t.passwordResetTokensRef(client)
	ref, _, err := passwordResetTokensRef.Add(ctx, model)
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
	passwordResetTokensRef := t.passwordResetTokensRef(client)

	iter := passwordResetTokensRef.Where("expires", "<=", now).Documents(ctx)
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
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()
	passwordResetTokensRef := t.passwordResetTokensRef(client)

	iter := passwordResetTokensRef.Where("email", "==", *email).Documents(ctx)
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
func (t *accounts) GetSignUpToken(signUpToken *string) (*models.SignUpToken, error) {
	client := t.provider.GetClient()
	ctx := context.Background()
	doc, err := t.signUpTokensRef(client).Doc(*signUpToken).Get(ctx)
	if err != nil {
		return nil, err
	}
	var model models.SignUpToken
	doc.DataTo(&model)
	return &model, nil
}
func (t *accounts) CreateUserAndAccount(user *models.User, account *models.Account) (*models.User, *models.Account, error) {
	client := t.provider.GetClient()
	batch := client.Batch()
	ctx := context.Background()

	usersRef := t.usersRef(client)
	userRef := usersRef.NewDoc()
	batch.Set(userRef, user)
	user.UserID = userRef.ID

	accountsRef := t.accountsRef(client)
	accountRef := accountsRef.Doc(user.Email)
	account.UserID = user.UserID
	batch.Set(accountRef, account)

	signUpTokensRef := t.signUpTokensRef(client)
	email := user.Email
	iter := signUpTokensRef.Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		batch.Delete(doc.Ref)
	}
	if _, err := batch.Commit(ctx); err != nil {
		return nil, nil, err
	}
	return user, account, nil
}
