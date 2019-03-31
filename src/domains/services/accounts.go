package services

import (
	"regexp"
	"unicode/utf8"

	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"
	"github.com/wakuwaku3/account-book.api/src/domains/models"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	accounts struct {
		crypt domains.Crypt
		repos domains.AccountsRepository
		clock cmn.Clock
	}

	// Accounts は Accountsのサービスです
	Accounts interface {
		ValidPassword(password *string) error
		ComparePassword(args *ComparePasswordArgs) error
		CreatePasswordResetToken(
			args *CreatePasswordResetTokenArgs) (
			*CreatePasswordResetTokenResult, error)
		SetPassword(args *SetPasswordArgs) error
		CreateSignUpToken(
			args *CreateSignUpTokenArgs) (
			*CreateSignUpTokenResult, error)
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
	// CreateSignUpTokenArgs は 引数です
	CreateSignUpTokenArgs struct {
		Email string
	}
	// CreateSignUpTokenResult は 結果です
	CreateSignUpTokenResult struct {
		SignUpToken string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	crypt domains.Crypt,
	repos domains.AccountsRepository,
	clock cmn.Clock,
) Accounts {
	return &accounts{crypt, repos, clock}
}

var passwordRegex = regexp.MustCompile(`^.*[0-9].*[a-z].*[A-Z]$|^.*[0-9].*[A-Z].*[a-z]$|^.*[a-z].*[0-9].*[A-Z]$|^.*[a-z].*[A-Z].*[0-9]$|^.*[A-Z].*[0-9].*[a-z]$|^.*[A-Z].*[a-z].*[0-9]$`)

func (t *accounts) ValidPassword(password *string) error {
	err := apperrors.NewClientError()
	if utf8.RuneCountInString(*password) < 8 {
		err.Append(apperrors.LessLengthPathword)
	}
	if !passwordRegex.MatchString(*password) {
		err.Append(apperrors.InvalidCharPassword)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) ComparePassword(args *ComparePasswordArgs) error {
	hashedPassword := *t.crypt.Hash(&args.InputPassword)
	if hashedPassword != args.HashedPassword {
		return apperrors.NewClientError(apperrors.FailureSignIn)
	}
	return nil
}
func (t *accounts) CreatePasswordResetToken(
	args *CreatePasswordResetTokenArgs,
) (*CreatePasswordResetTokenResult, error) {
	if _, err := t.repos.Get(&args.Email); err != nil {
		return nil, nil
	}
	expires := t.clock.Now().AddDate(0, 0, 2)
	id, err := t.repos.CreatePasswordResetToken(&models.PasswordResetToken{
		Email:   args.Email,
		Expires: expires,
	})
	if err != nil {
		return nil, err
	}
	go t.repos.CleanUpPasswordResetToken()
	return &CreatePasswordResetTokenResult{
		PasswordResetToken: *id,
	}, nil
}
func (t *accounts) SetPassword(args *SetPasswordArgs) error {
	hashedPassword := t.crypt.Hash(&args.Password)
	if err := t.repos.SetPassword(&args.Email, hashedPassword); err != nil {
		return err
	}
	go t.repos.CleanUpPasswordResetTokenByEmail(&args.Email)
	return nil
}
func (t *accounts) CreateSignUpToken(
	args *CreateSignUpTokenArgs) (
	*CreateSignUpTokenResult, error) {
	expires := t.clock.Now().AddDate(0, 0, 2)
	id, err := t.repos.CreateSignUpToken(&models.SignUpToken{
		Email:   args.Email,
		Expires: expires,
	})
	if err != nil {
		return nil, err
	}
	go t.repos.CleanUpSignUpToken()
	return &CreateSignUpTokenResult{
		SignUpToken: *id,
	}, nil
}
