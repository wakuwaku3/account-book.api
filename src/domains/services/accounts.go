package services

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/wakuwaku3/account-book.api/src/domains/models"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	accounts struct {
		crypt domains.Crypt
		repos domains.AccountsRepository
	}

	// Accounts は Accountsのサービスです
	Accounts interface {
		ValidPassword(password *string) error
		ComparePassword(args *ComparePasswordArgs) error
		CreatePasswordResetToken(
			args *CreatePasswordResetTokenArgs) (
			*CreatePasswordResetTokenResult, error)
		SetPassword(args *SetPasswordArgs) error
	}
	// ComparePasswordArgs は 引数です
	ComparePasswordArgs struct {
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
	// SetPasswordArgs は 引数です
	SetPasswordArgs struct {
		Password string
		Email    string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	crypt domains.Crypt,
	repos domains.AccountsRepository) Accounts {
	return &accounts{crypt, repos}
}

var passwordRegex = regexp.MustCompile(`^.*[0-9].*[a-z].*[A-Z]$|^.*[0-9].*[A-Z].*[a-z]$|^.*[a-z].*[0-9].*[A-Z]$|^.*[a-z].*[A-Z].*[0-9]$|^.*[A-Z].*[0-9].*[a-z]$|^.*[A-Z].*[a-z].*[0-9]$`)

func (t *accounts) ValidPassword(password *string) error {
	array := make([]string, 0)
	if utf8.RuneCountInString(*password) < 8 {
		array = append(array, "パスワードは8文字以上設定してください。")
	}
	if !passwordRegex.MatchString(*password) {
		array = append(array, "パスワードには、半角英小文字、大文字、数字をそれぞれ1種類以上使用してください。")
	}
	if len(array) > 0 {
		return errors.New(strings.Join(array, "\n"))
	}
	return nil
}
func (t *accounts) ComparePassword(args *ComparePasswordArgs) error {
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
	id, err := t.repos.CreatePasswordResetToken(&models.PasswordResetToken{
		Email:   args.Email,
		Expires: expires,
	})
	if err != nil {
		return nil, err
	}
	go t.repos.CleanUp()
	return &CreatePasswordResetTokenResult{
		PasswordResetToken: *id,
	}, nil
}
func (t *accounts) SetPassword(args *SetPasswordArgs) error {
	hashedPassword := t.crypt.Hash(&args.Password)
	if err := t.repos.SetPassword(&args.Email, hashedPassword); err != nil {
		return err
	}
	go t.repos.CleanUpByEmail(&args.Email)
	return nil
}
