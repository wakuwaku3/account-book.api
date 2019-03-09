package ctrls

import (
	"net/http"

	"github.com/wakuwaku3/account-book.api/src/ctrls/responses"

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
	res, clientErr, err := t.useCase.SignIn(request.Convert())
	if err != nil {
		return err
	}
	if clientErr != nil {
		return responses.WriteErrorResponse(c, clientErr)
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
	clientErr, err := t.useCase.PasswordResetRequesting(&usecases.PasswordResetRequestingArgs{Email: request.Email})
	if err != nil {
		return err
	}
	if clientErr != nil {
		return responses.WriteErrorResponse(c, clientErr)
	}
	return responses.WriteEmptyResponse(c)
}
func (t *accounts) GetResetPasswordModel(c echo.Context) error {
	return c.JSON(http.StatusOK, "res")
}
func (t *accounts) ResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, "res")
}
