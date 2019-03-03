package services

import (
	"errors"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	accounts struct {
		crypt domains.Crypt
	}

	// Accounts は Accountsのサービスです
	Accounts interface {
		ValidPassword(args *ValidPasswordArgs) error
	}
	// ValidPasswordArgs は Password.Validの引数です
	ValidPasswordArgs struct {
		HashedPassword string
		InputPassword  string
	}
)

// NewAccounts is create instance.
func NewAccounts(crypt domains.Crypt) Accounts {
	return &accounts{crypt: crypt}
}
func (t *accounts) ValidPassword(args *ValidPasswordArgs) error {
	hashedPassword := t.crypt.Hash(&args.InputPassword)
	if *hashedPassword != args.HashedPassword {
		return errors.New("パスワードが違います。")
	}
	return nil
}
