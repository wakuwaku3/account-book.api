package usecases

import (
	"errors"
	"strings"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/services"
)

type (
	accounts struct {
		query             AccountsQuery
		jwt               domains.Jwt
		claimsProvider    domains.ClaimsProvider
		service           services.Accounts
		resetPasswordMail domains.ResetPasswordMail
	}
	// Accounts is AccountsController
	Accounts interface {
		SignIn(args *SignInArgs) (*SignInResult, error, error)
		Refresh(args *RefreshArgs) (*RefreshResult, error)
		PasswordResetRequesting(args *PasswordResetRequestingArgs) (error, error)
		GetResetPasswordModel(args *ResetPasswordArgs) (*GetResetPasswordModelResult, error)
		ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error)
	}
	// SignInArgs は 引数です
	SignInArgs struct {
		Email    string
		Password string
	}
	// SignInResult は 結果です
	SignInResult struct {
		Token        string
		RefreshToken string
	}
	// RefreshArgs は 引数です
	RefreshArgs struct {
		RefreshToken string
	}
	// RefreshResult は 結果です
	RefreshResult struct {
		Token        string
		RefreshToken string
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
		Token        string
		RefreshToken string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	query AccountsQuery,
	jwt domains.Jwt,
	service services.Accounts,
	resetPasswordMail domains.ResetPasswordMail,
) Accounts {
	return &accounts{
		query:             query,
		jwt:               jwt,
		service:           service,
		resetPasswordMail: resetPasswordMail,
	}
}
func (t *accounts) SignIn(args *SignInArgs) (*SignInResult, error, error) {
	err := args.valid()
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
	refreshToken, err := t.jwt.CreateRefreshToken(&info.JwtRefreshClaims)
	if err != nil {
		return nil, nil, err
	}
	return &SignInResult{
		Token:        *token,
		RefreshToken: *refreshToken,
	}, nil, nil
}
func (t *SignInArgs) valid() error {
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
func (t *accounts) Refresh(args *RefreshArgs) (*RefreshResult, error) {
	claims, err := t.jwt.ParseRefreshToken(&args.RefreshToken)
	if err != nil {
		return nil, err
	}
	info, err := t.query.GetRefreshInfo(&claims.Email)
	if err != nil {
		return nil, err
	}
	if info.AccountToken != claims.AccountToken {
		return nil, errors.New("accountToken does not match")
	}
	token, err := t.jwt.CreateToken(&info.JwtClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.jwt.CreateRefreshToken(&info.JwtRefreshClaims)
	if err != nil {
		return nil, err
	}
	return &RefreshResult{
		Token:        *token,
		RefreshToken: *refreshToken,
	}, nil
}
func (t *accounts) PasswordResetRequesting(args *PasswordResetRequestingArgs) (error, error) {
	if err := args.valid(); err != nil {
		return err, nil
	}
	token, err := t.service.CreatePasswordResetToken(&services.CreatePasswordResetTokenArgs{
		Email: args.Email,
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		go func() {
			args := domains.ResetPasswordMailSendArgs{
				Email: args.Email,
				Token: token.PasswordResetToken,
			}
			t.resetPasswordMail.Send(&args)
		}()
	}
	return nil, nil
}
func (t *PasswordResetRequestingArgs) valid() error {
	array := make([]string, 0)
	if t.Email == "" {
		array = append(array, "メールアドレスが入力されていません。")
	}
	if len(array) > 0 {
		return errors.New(strings.Join(array, "\n"))
	}
	return nil
}
func (t *accounts) GetResetPasswordModel(args *ResetPasswordArgs) (*GetResetPasswordModelResult, error) {
	return &GetResetPasswordModelResult{}, nil
}
func (t *accounts) ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error) {
	return &ResetPasswordResult{}, nil
}
