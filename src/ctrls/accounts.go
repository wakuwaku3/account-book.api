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
	claimResponse struct {
		Token    string `json:"token"`
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
		Email    string `json:"email"`
	}
	signInResponse struct {
		Claim claimResponse `json:"claim"`
	}
	refreshResponse struct {
		Claim claimResponse `json:"claim"`
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
		return c.JSON(http.StatusBadRequest, responses.Response{
			Errors: []string{clientErr.Error()},
		})
	}
	return c.JSON(http.StatusOK, responses.Response{
		Result: signInResponse{
			Claim: claimResponse{
				Token:    res.Claims.Token,
				UserID:   res.Claims.UserID,
				UserName: res.Claims.UserName,
				Email:    res.Claims.Email,
			},
		},
	})
}
func (t *signInRequest) Convert() *usecases.SignInArgs {
	return &usecases.SignInArgs{
		Email:    t.Email,
		Password: t.Password,
	}
}
func (t *accounts) Refresh(c echo.Context) error {
	res, err := t.useCase.Refresh()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responses.Response{
		Result: refreshResponse{
			Claim: claimResponse{
				Token:    res.Claims.Token,
				UserID:   res.Claims.UserID,
				UserName: res.Claims.UserName,
				Email:    res.Claims.Email,
			},
		},
	})
}
func (t *accounts) PasswordResetRequesting(c echo.Context) error {
	request := new(passwordResetRequestingRequest)
	if err := c.Bind(&request); err != nil {
		return err
	}
	cErr, err := t.useCase.PasswordResetRequesting(&usecases.PasswordResetRequestingArgs{Email: request.Email})
	if err != nil {
		return err
	}
	if cErr != nil {
		return c.JSON(http.StatusBadRequest, responses.Response{
			Errors: []string{cErr.Error()},
		})
	}
	return c.NoContent(http.StatusNoContent)
}
func (t *accounts) GetResetPasswordModel(c echo.Context) error {
	return c.JSON(http.StatusOK, "res")
}
func (t *accounts) ResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, "res")
}
