package usecases

import (
	"errors"

	"github.com/wakuwaku3/account-book.api/src/domains"
	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"
	"github.com/wakuwaku3/account-book.api/src/domains/services"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	accounts struct {
		query             AccountsQuery
		jwt               domains.Jwt
		service           services.Accounts
		resetPasswordMail domains.ResetPasswordMail
		userCreationMail  domains.UserCreationMail
		userExistingMail  domains.UserExistingMail
		clock             cmn.Clock
	}
	// Accounts is AccountsController
	Accounts interface {
		SignIn(args *SignInArgs) (*SignInResult, error)
		Refresh(args *RefreshArgs) (*RefreshResult, error)
		PasswordResetRequesting(args *PasswordResetRequestingArgs) error
		GetResetPasswordModel(args *GetResetPasswordModelArgs) (*GetResetPasswordModelResult, error)
		ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error)
		SignUpRequesting(args *SignUpRequestingArgs) error
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
	// PasswordResetRequestingArgs は 引数です
	PasswordResetRequestingArgs struct {
		Email string
	}
	// GetResetPasswordModelArgs は 引数です
	GetResetPasswordModelArgs struct {
		PasswordResetToken string
	}
	// GetResetPasswordModelResult は 結果です
	GetResetPasswordModelResult struct {
		Email string
	}
	// ResetPasswordArgs は 引数です
	ResetPasswordArgs struct {
		PasswordResetToken string
		Password           string
	}
	// ResetPasswordResult は 結果です
	ResetPasswordResult struct {
		Token        string
		RefreshToken string
	}
	// SignUpRequestingArgs は 引数です
	SignUpRequestingArgs struct {
		Email string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	query AccountsQuery,
	jwt domains.Jwt,
	service services.Accounts,
	resetPasswordMail domains.ResetPasswordMail,
	userCreationMail domains.UserCreationMail,
	userExistingMail domains.UserExistingMail,
	clock cmn.Clock,
) Accounts {
	return &accounts{
		query,
		jwt,
		service,
		resetPasswordMail,
		userCreationMail,
		userExistingMail,
		clock,
	}
}
func (t *accounts) SignIn(args *SignInArgs) (*SignInResult, error) {
	err := args.valid()
	if err != nil {
		return nil, err
	}
	info, err := t.query.GetSignInInfo(&args.Email)
	if err != nil {
		return nil, err
	}
	err = t.service.ComparePassword(&services.ComparePasswordArgs{
		HashedPassword: info.HashedPassword,
		InputPassword:  args.Password,
	})
	if err != nil {
		return nil, err
	}
	token, err := t.jwt.CreateToken(&info.JwtClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.jwt.CreateRefreshToken(&info.JwtRefreshClaims)
	if err != nil {
		return nil, err
	}
	return &SignInResult{
		Token:        *token,
		RefreshToken: *refreshToken,
	}, nil
}
func (t *SignInArgs) valid() error {
	err := apperrors.NewClientError()
	if t.Email == "" {
		err.Append(apperrors.RequiredMailAddress)
	}
	if t.Password == "" {
		err.Append(apperrors.RequiredPassword)
	}
	if err.HasError() {
		return err
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
func (t *accounts) PasswordResetRequesting(args *PasswordResetRequestingArgs) error {
	if err := args.valid(); err != nil {
		return err
	}
	token, err := t.service.CreatePasswordResetToken(&services.CreatePasswordResetTokenArgs{
		Email: args.Email,
	})
	if err != nil {
		return err
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
	return nil
}
func (t *PasswordResetRequestingArgs) valid() error {
	err := apperrors.NewClientError()
	if t.Email == "" {
		err.Append(apperrors.RequiredMailAddress)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) GetResetPasswordModel(args *GetResetPasswordModelArgs) (*GetResetPasswordModelResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	info, err := t.query.GetResetPasswordModelInfo(&args.PasswordResetToken)
	if err != nil {
		return nil, err
	}
	if info.Expires.Before(t.clock.Now()) {
		return nil, errors.New("URLの有効期限が切れています。")
	}
	return &GetResetPasswordModelResult{
		Email: info.Email,
	}, nil
}
func (t *GetResetPasswordModelArgs) valid() error {
	err := apperrors.NewClientError()
	if t.PasswordResetToken == "" {
		err.Append(apperrors.RequiredPasswordToken)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) ResetPassword(args *ResetPasswordArgs) (*ResetPasswordResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	if err := t.service.ValidPassword(&args.Password); err != nil {
		return nil, err
	}
	info, err := t.query.GetResetPasswordInfo(&args.PasswordResetToken)
	if err != nil {
		return nil, err
	}
	if info.Expires.Before(t.clock.Now()) {
		return nil, errors.New("URLの有効期限が切れています。")
	}
	setPasswordArgs := &services.SetPasswordArgs{
		Email:    info.Email,
		Password: args.Password,
	}
	if err := t.service.SetPassword(setPasswordArgs); err != nil {
		return nil, err
	}
	token, err := t.jwt.CreateToken(&info.JwtClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.jwt.CreateRefreshToken(&info.JwtRefreshClaims)
	if err != nil {
		return nil, err
	}
	return &ResetPasswordResult{
		Token:        *token,
		RefreshToken: *refreshToken,
	}, nil
}
func (t *ResetPasswordArgs) valid() error {
	err := apperrors.NewClientError()
	if t.PasswordResetToken == "" {
		err.Append(apperrors.RequiredPasswordToken)
	}
	if t.Password == "" {
		err.Append(apperrors.RequiredPassword)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) SignUpRequesting(args *SignUpRequestingArgs) error {
	if err := args.valid(); err != nil {
		return err
	}

	// 既にメールアドレスが使用されている場合、Tokenが生成される
	token, err := t.service.CreatePasswordResetToken(&services.CreatePasswordResetTokenArgs{
		Email: args.Email,
	})
	if err != nil {
		return err
	}
	if token == nil {
		// メールアドレスが使用されてない場合、UserCreationメールを送る
		token, err := t.service.CreateSignUpToken(&services.CreateSignUpTokenArgs{
			Email: args.Email,
		})
		if err != nil {
			return err
		}
		go func() {
			t.userCreationMail.Send(&domains.UserCreationMailSendArgs{
				Email: args.Email,
				Token: token.SignUpToken,
			})
		}()
	}
	if token != nil {
		// メールアドレスが使用されていた場合、UserExistingメールを送る
		go func() {
			t.userExistingMail.Send(&domains.UserExistingMailSendArgs{
				Email: args.Email,
				Token: token.PasswordResetToken,
			})
		}()
	}
	return nil
}
func (t *SignUpRequestingArgs) valid() error {
	err := apperrors.NewClientError()
	if t.Email == "" {
		err.Append(apperrors.RequiredMailAddress)
	}
	if err.HasError() {
		return err
	}
	return nil
}
