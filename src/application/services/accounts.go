package services

import (
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/wakuwaku3/account-book.api/src/application"
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
		CreateUser(args *CreateUserArgs) (*CreateUserResult, error)
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
	// CreateUserArgs は 引数です
	CreateUserArgs struct {
		Email    string
		Password string
		UserName string
		Culture  string
	}
	// CreateUserResult は 引数です
	CreateUserResult struct {
		JwtClaims        domains.JwtClaims
		JwtRefreshClaims domains.JwtRefreshClaims
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
	err := application.NewClientError()
	if utf8.RuneCountInString(*password) < 8 {
		err.Append(application.LessLengthPathword)
	}
	if !passwordRegex.MatchString(*password) {
		err.Append(application.InvalidCharPassword)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) ComparePassword(args *ComparePasswordArgs) error {
	hashedPassword := *t.crypt.Hash(&args.InputPassword)
	if hashedPassword != args.HashedPassword {
		return application.NewClientError(application.FailureSignIn)
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
func (t *accounts) CreateUser(args *CreateUserArgs) (*CreateUserResult, error) {
	hashedPassword := t.crypt.Hash(&args.Password)
	now := t.clock.Now()
	token := uuid.New().String()
	user, account, err := t.repos.CreateUserAndAccount(&models.User{
		Email:        args.Email,
		Culture:      args.Culture,
		UserName:     args.UserName,
		UseStartDate: now,
	}, &models.Account{
		Email:          args.Email,
		AccountToken:   token,
		HashedPassword: *hashedPassword,
	})
	if err != nil {
		return nil, err
	}
	return &CreateUserResult{
		JwtClaims: domains.JwtClaims{
			Email:        account.Email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        account.Email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}
