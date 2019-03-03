package usecases

import (
	"errors"
	"strings"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	accounts struct {
		query   AccountsQuery
		jwt     domains.Jwt
		service services.Accounts
	}
	// Accounts is AccountsController
	Accounts interface {
		SignIn(args *SignInArgs) (*SignInResult, error, error)
		Refresh() (*RefreshResult, error)
		PasswordResetRequesting(args *PasswordResetRequestingArgs) error
		GetResetPasswordModel(args *ResetPasswordArgs) (*GetResetPasswordModelResult, error)
		ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error)
	}
	// SignInArgs は 引数です
	SignInArgs struct {
		Email    string
		Password string
	}
	// Claims は Claimsです
	Claims struct {
		Token    string
		UserID   string
		UserName string
		Email    string
	}
	// SignInResult は 結果です
	SignInResult struct {
		Claims Claims
	}
	RefreshResult struct {
		Claims Claims
	}
	PasswordResetRequestingArgs struct {
		Email string
	}
	ResetPasswordArgs struct {
		ResetPasswordToken string
	}
	GetResetPasswordModelResult struct {
		Email string
	}
	ResetPasswordResult struct {
		Claims Claims
	}
)

// NewAccounts is create instance.
func NewAccounts(
	query AccountsQuery,
	jwt domains.Jwt,
	service services.Accounts) Accounts {
	return &accounts{
		query:   query,
		jwt:     jwt,
		service: service,
	}
}
func (t *accounts) SignIn(args *SignInArgs) (*SignInResult, error, error) {
	err := args.Valid()
	if err != nil {
		return nil, err, nil
	}
	info, err := t.query.GetSignInInfo(&args.Email)
	if err != nil {
		return nil, nil, err
	}
	err = t.service.ValidPassword(&services.ValidPasswordArgs{
		HashedPassword: info.HashedPassword,
		InputPassword:  args.Password,
	})
	if err != nil {
		return nil, err, nil
	}
	token, err := t.jwt.CreateToken(&info.JwtClaims)
	if err != nil {
		return nil, nil, err
	}
	return &SignInResult{
		Claims: Claims{
			Token:    *token,
			UserID:   info.JwtClaims.UserID,
			Email:    info.JwtClaims.Email,
			UserName: info.JwtClaims.UserName,
		},
	}, nil, nil
}

// Valid は SignInArgsを検証します
func (t *SignInArgs) Valid() error {
	array := make([]string, 0)
	if t.Email == "" {
		array = append(array, "メールアドレスが入力されていません。")
	}
	if t.Password == "" {
		array = append(array, "パスワードが入力されていません。")
	}
	if len(array) > 0 {
		return errors.New(strings.Join(array, "\n"))
	}
	return nil
}
func (t *accounts) Refresh() (*RefreshResult, error) {
	return &RefreshResult{}, nil
}
func (t *accounts) PasswordResetRequesting(args *PasswordResetRequestingArgs) error {
	return nil
}
func (t *accounts) GetResetPasswordModel(args *ResetPasswordArgs) (*GetResetPasswordModelResult, error) {
	return &GetResetPasswordModelResult{}, nil
}
func (t *accounts) ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error) {
	return &ResetPasswordResult{}, nil
}
