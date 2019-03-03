package web

import (
	"net/http"

	"github.com/labstack/echo"
	ctrl "github.com/wakuwaku3/account-book.api/src/ctrls"
)

func (web *web) setRoute() *web {
	web.echo.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello")
	})
	web.echo.GET("/_ah/warmup", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// accounts
	// sign-in
	web.echo.POST("/accounts/sign-in", func(c echo.Context) error {
		return web.container.Invoke(func(accounts ctrl.Accounts) error {
			return accounts.SignIn(c)
		})
	})
	// refresh
	web.echo.POST("/accounts/refresh", func(c echo.Context) error {
		return web.container.Invoke(func(accounts ctrl.Accounts) error {
			return accounts.Refresh(c)
		})
	})
	// password-reset-requesting
	web.echo.PUT("/accounts/password-reset-requesting", func(c echo.Context) error {
		return web.container.Invoke(func(accounts ctrl.Accounts) error {
			return accounts.PasswordResetRequesting(c)
		})
	})
	// reset-password
	web.echo.GET("/accounts/reset-password", func(c echo.Context) error {
		return web.container.Invoke(func(accounts ctrl.Accounts) error {
			return accounts.GetResetPasswordModel(c)
		})
	})
	web.echo.POST("/accounts/reset-password", func(c echo.Context) error {
		return web.container.Invoke(func(accounts ctrl.Accounts) error {
			return accounts.ResetPassword(c)
		})
	})

	return web
}
