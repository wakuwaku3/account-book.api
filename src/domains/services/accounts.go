package services

import (
	"errors"
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	accounts struct {
		crypt domains.Crypt
		repos domains.AccountsRepository
	}

	// Accounts は Accountsのサービスです
	Accounts interface {
		ValidPassword(args *ValidPasswordArgs) error
		CreatePasswordResetToken(
			args *CreatePasswordResetTokenArgs) (
			*CreatePasswordResetTokenResult, error)
	}
	// ValidPasswordArgs は Password.Validの引数です
	ValidPasswordArgs struct {
		HashedPassword string
		InputPassword  string
	}
	// CreatePasswordResetTokenArgs は 引数です
	CreatePasswordResetTokenArgs struct {
		Email string
	}
	// CreatePasswordResetTokenResult は 結果です
	CreatePasswordResetTokenResult struct {
		PasswordResetToken string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	crypt domains.Crypt,
	repos domains.AccountsRepository) Accounts {
	return &accounts{crypt, repos}
}
func (t *accounts) ValidPassword(args *ValidPasswordArgs) error {
	hashedPassword := t.crypt.Hash(&args.InputPassword)
	if *hashedPassword != args.HashedPassword {
		return errors.New("パスワードが違います。")
	}
	return nil
}
func (t *accounts) CreatePasswordResetToken(
	args *CreatePasswordResetTokenArgs,
) (*CreatePasswordResetTokenResult, error) {
	if _, err := t.repos.Get(&args.Email); err != nil {
		return nil, nil
	}
	expires := time.Now().AddDate(0, 0, 2)
	id, err := t.repos.CreatePasswordResetToken(&args.Email, &expires)
	if err != nil {
		return nil, err
	}
	go t.repos.CleanUp()
	return &CreatePasswordResetTokenResult{
		PasswordResetToken: *id,
	}, nil
}
