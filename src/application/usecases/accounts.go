package usecases

import (
	"errors"

	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/application/services"
	"github.com/wakuwaku3/account-book.api/src/infrastructures/cmn"
)

type (
	accounts struct {
		query             AccountsQuery
		jwt               application.Jwt
		service           services.Accounts
		resetPasswordMail application.ResetPasswordMail
		userCreationMail  application.UserCreationMail
		userExistingMail  application.UserExistingMail
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
		GetSignUpModel(args *GetSignUpModelArgs) (*GetSignUpModelResult, error)
		SignUp(args *SignUpArgs) (*SignUpResult, error)
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
	// GetSignUpModelArgs は 引数です
	GetSignUpModelArgs struct {
		SignUpToken string
	}
	// GetSignUpModelResult は 結果です
	GetSignUpModelResult struct {
		Email string
	}
	// SignUpArgs は 引数です
	SignUpArgs struct {
		SignUpToken string
		Password    string
		UserName    string
		Culture     string
	}
	// SignUpResult は 結果です
	SignUpResult struct {
		Token        string
		RefreshToken string
	}
)

// NewAccounts is create instance.
func NewAccounts(
	query AccountsQuery,
	jwt application.Jwt,
	service services.Accounts,
	resetPasswordMail application.ResetPasswordMail,
	userCreationMail application.UserCreationMail,
	userExistingMail application.UserExistingMail,
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
	err := application.NewClientError()
	if t.Email == "" {
		err.Append(application.RequiredMailAddress)
	}
	if t.Password == "" {
		err.Append(application.RequiredPassword)
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
			args := application.ResetPasswordMailSendArgs{
				Email: args.Email,
				Token: token.PasswordResetToken,
			}
			t.resetPasswordMail.Send(&args)
		}()
	}
	return nil
}
func (t *PasswordResetRequestingArgs) valid() error {
	err := application.NewClientError()
	if t.Email == "" {
		err.Append(application.RequiredMailAddress)
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
	err := application.NewClientError()
	if t.PasswordResetToken == "" {
		err.Append(application.RequiredPasswordToken)
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
	err := application.NewClientError()
	if t.PasswordResetToken == "" {
		err.Append(application.RequiredPasswordToken)
	}
	if t.Password == "" {
		err.Append(application.RequiredPassword)
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
			t.userCreationMail.Send(&application.UserCreationMailSendArgs{
				Email: args.Email,
				Token: token.SignUpToken,
			})
		}()
	}
	if token != nil {
		// メールアドレスが使用されていた場合、UserExistingメールを送る
		go func() {
			t.userExistingMail.Send(&application.UserExistingMailSendArgs{
				Email: args.Email,
				Token: token.PasswordResetToken,
			})
		}()
	}
	return nil
}
func (t *SignUpRequestingArgs) valid() error {
	err := application.NewClientError()
	if t.Email == "" {
		err.Append(application.RequiredMailAddress)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) GetSignUpModel(args *GetSignUpModelArgs) (*GetSignUpModelResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	info, err := t.query.GetSignUpModelInfo(&args.SignUpToken)
	if err != nil {
		return nil, err
	}
	if info.Expires.Before(t.clock.Now()) {
		return nil, application.NewClientError(application.ExpiredURL)
	}
	return &GetSignUpModelResult{
		Email: info.Email,
	}, nil
}
func (t *GetSignUpModelArgs) valid() error {
	err := application.NewClientError()
	if t.SignUpToken == "" {
		err.Append(application.RequiredSignUpToken)
	}
	if err.HasError() {
		return err
	}
	return nil
}
func (t *accounts) SignUp(args *SignUpArgs) (*SignUpResult, error) {
	if err := args.valid(); err != nil {
		return nil, err
	}
	if err := t.service.ValidPassword(&args.Password); err != nil {
		return nil, err
	}
	info, err := t.query.GetSignUpModelInfo(&args.SignUpToken)
	if err != nil {
		return nil, err
	}
	if info.Expires.Before(t.clock.Now()) {
		return nil, application.NewClientError(application.ExpiredURL)
	}

	result, err := t.service.CreateUser(&services.CreateUserArgs{
		Email:    info.Email,
		Password: args.Password,
		Culture:  args.Culture,
		UserName: args.UserName,
	})
	if err != nil {
		return nil, err
	}
	token, err := t.jwt.CreateToken(&result.JwtClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := t.jwt.CreateRefreshToken(&result.JwtRefreshClaims)
	if err != nil {
		return nil, err
	}
	return &SignUpResult{
		Token:        *token,
		RefreshToken: *refreshToken,
	}, nil
}
func (t *SignUpArgs) valid() error {
	err := application.NewClientError()
	if t.SignUpToken == "" {
		err.Append(application.RequiredSignUpToken)
	}
	if t.Password == "" {
		err.Append(application.RequiredPassword)
	}
	if t.UserName == "" {
		err.Append(application.RequiredName)
	}
	if t.Culture != "ja" && t.Culture != "en" {
		err.Append(application.InValidCulture)
	}
	if err.HasError() {
		return err
	}
	return nil
}
