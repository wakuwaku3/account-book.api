package ctrls

import (
	"github.com/labstack/gommon/log"
	"github.com/wakuwaku3/account-book.api/src/ctrls/responses"
	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"

	"github.com/wakuwaku3/account-book.api/src/usecases"

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
		if _, ok := err.(apperrors.ClientError); !ok {
			log.Error(err)
			return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.FailureSignIn))
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
		return responses.WriteUnAuthorizedErrorResponse(c)
	}
	res, err := t.useCase.Refresh(request.Convert())
	if err != nil {
		if _, ok := err.(apperrors.ClientError); !ok {
			log.Error(err)
			return responses.WriteUnAuthorizedErrorResponse(c)
		}
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
		if _, ok := err.(apperrors.ClientError); !ok {
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
		if _, ok := err.(apperrors.ClientError); !ok {
			log.Error(err)
			return responses.WriteErrorResponse(c, apperrors.NewClientError(apperrors.FailurePasswordReset))
		}
		return responses.WriteErrorResponse(c, err)
	}
	return responses.WriteResponse(c, resetPasswordResponse{
		Token:        res.Token,
		RefreshToken: res.RefreshToken,
	})
}
