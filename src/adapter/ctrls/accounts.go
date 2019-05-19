package ctrls

import (
	"github.com/labstack/gommon/log"
	"github.com/wakuwaku3/account-book.api/src/adapter/ctrls/responses"
	"github.com/wakuwaku3/account-book.api/src/application"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"

	"github.com/labstack/echo"
)

type (
	accounts struct {
		useCase usecases.Accounts
	}
	// Accounts is AccountsController
	Accounts interface {
		SignIn(c echo.Context) error
		Refresh(c echo.Context) error
		PasswordResetRequesting(c echo.Context) error
		GetResetPasswordModel(c echo.Context) error
		ResetPassword(c echo.Context) error
		SignUpRequesting(c echo.Context) error
		GetSignUpModel(c echo.Context) error
		SignUp(c echo.Context) error
	}
	signInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	signInResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}
	refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	refreshResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}
	passwordResetRequestingRequest struct {
		Email string `json:"email"`
	}
	getResetPasswordModelRequest struct {
		PasswordResetToken string `query:"passwordResetToken"`
	}
	getResetPasswordModelResponse struct {
		Email string `json:"email"`
	}
	resetPasswordRequest struct {
		PasswordResetToken string `json:"passwordResetToken"`
		Password           string `json:"password"`
	}
	resetPasswordResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}
	signUpRequestingRequest struct {
		Email string `json:"email"`
	}
	getSignUpModelRequest struct {
		SignUpToken string
	}
	getSignUpModelResponse struct {
		Email string `json:"email"`
	}
	signUpRequest struct {
		SignUpToken string `json:"signUpToken"`
		Password    string `json:"password"`
		UserName    string `json:"userName"`
		Culture     string `json:"culture"`
		Agreement   bool   `json:"agreement"`
	}
	signUpResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}
)

// NewAccounts is create instance.
func NewAccounts(useCase usecases.Accounts) Accounts {
	return &accounts{useCase: useCase}
}
func (t *accounts) SignIn(c echo.Context) error {
	request := new(signInRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.SignIn(request.Convert())
	if err != nil {
		if _, ok := err.(application.ClientError); !ok {
			log.Error(err)
			return responses.WriteErrorResponse(c, application.NewClientError(application.FailureSignIn))
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, signInResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
}
func (t *signInRequest) Convert() *usecases.SignInArgs {
	return &usecases.SignInArgs{
		Email:    t.Email,
		Password: t.Password,
	}
}
func (t *accounts) Refresh(c echo.Context) error {
	request := new(refreshRequest)
	if err := c.Bind(&request); err != nil {
		c.Logger().Error(err, c.Request())
		return responses.WriteUnAuthorizedErrorResponse(c)
	}
	res, err := t.useCase.Refresh(request.Convert())

	// エラー原因がわからないためエラーを詳しく出してみる
	if err != nil {
		req := c.Request()
		if cErr, ok := err.(application.ClientError); ok {
			if request == nil {
				c.Logger().Error(err, cErr.GetErrorCodes(), req)
				return responses.WriteUnAuthorizedErrorResponse(c)
			}
			if req == nil {
				c.Logger().Error(err, cErr.GetErrorCodes())
				return responses.WriteUnAuthorizedErrorResponse(c)
			}
			c.Logger().Error(err, cErr.GetErrorCodes(), request, req)
			return responses.WriteUnAuthorizedErrorResponse(c)
		}
		if request == nil {
			c.Logger().Error(err, req)
			return responses.WriteUnAuthorizedErrorResponse(c)
		}
		if req == nil {
			c.Logger().Error(err)
			return responses.WriteUnAuthorizedErrorResponse(c)
		}
		c.Logger().Error(err, request, req)
		return responses.WriteUnAuthorizedErrorResponse(c)
	}

	return responses.WriteResponse(c, refreshResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
}
func (t *refreshRequest) Convert() *usecases.RefreshArgs {
	return &usecases.RefreshArgs{
		RefreshToken: t.RefreshToken,
	}
}
func (t *accounts) PasswordResetRequesting(c echo.Context) error {
	request := new(passwordResetRequestingRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	err := t.useCase.PasswordResetRequesting(&usecases.PasswordResetRequestingArgs{Email: request.Email})
	if err != nil {
		if _, ok := err.(application.ClientError); !ok {
			log.Error(err)
			return responses.WriteEmptyResponse(c)
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *accounts) GetResetPasswordModel(c echo.Context) error {
	passwordResetToken := c.QueryParam("passwordResetToken")
	request := &getResetPasswordModelRequest{
		PasswordResetToken: passwordResetToken,
	}
	res, err := t.useCase.GetResetPasswordModel(&usecases.GetResetPasswordModelArgs{
		PasswordResetToken: request.PasswordResetToken,
	})
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, getResetPasswordModelResponse{
		Email: res.Email,
	})
}
func (t *accounts) ResetPassword(c echo.Context) error {
	request := new(resetPasswordRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	res, err := t.useCase.ResetPassword(&usecases.ResetPasswordArgs{
		PasswordResetToken: request.PasswordResetToken,
		Password:           request.Password,
	})
	if err != nil {
		if _, ok := err.(application.ClientError); !ok {
			log.Error(err)
			return responses.WriteErrorResponse(c, application.NewClientError(application.FailurePasswordReset))
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, resetPasswordResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
}
func (t *accounts) SignUpRequesting(c echo.Context) error {
	request := new(signUpRequestingRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	err := t.useCase.SignUpRequesting(&usecases.SignUpRequestingArgs{Email: request.Email})
	if err != nil {
		if _, ok := err.(application.ClientError); !ok {
			log.Error(err)
			return responses.WriteEmptyResponse(c)
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *accounts) GetSignUpModel(c echo.Context) error {
	signUpToken := c.QueryParam("signUpToken")
	request := &getSignUpModelRequest{
		SignUpToken: signUpToken,
	}
	res, err := t.useCase.GetSignUpModel(&usecases.GetSignUpModelArgs{
		SignUpToken: request.SignUpToken,
	})
	if err != nil {
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, getSignUpModelResponse{
		Email: res.Email,
	})
}
func (t *accounts) SignUp(c echo.Context) error {
	request := new(signUpRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	if !request.Agreement {
		return responses.WriteErrorResponse(c, application.NewClientError(application.RequiredAgreement))
	}
	res, err := t.useCase.SignUp(&usecases.SignUpArgs{
		SignUpToken: request.SignUpToken,
		Password:    request.Password,
		UserName:    request.UserName,
		Culture:     request.Culture,
	})
	if err != nil {
		if _, ok := err.(application.ClientError); !ok {
			log.Error(err)
			return responses.WriteErrorResponse(c, application.NewClientError(application.FailureSignUp))
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, signUpResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
}
