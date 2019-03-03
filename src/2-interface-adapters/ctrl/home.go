package ctrl

import (
	"net/http"

	"github.com/wakuwaku3/account-book.api/src/1-application-business-rules/usecases"

	"github.com/labstack/echo"
)

type (
	home struct {
		envProvider usecases.Env
	}
	// Home is HomeController
	Home interface {
		Get(c echo.Context) error
	}
)

// NewHome is create instance.
func NewHome(envProvider usecases.Env) Home {
	return &home{envProvider: envProvider}
}
func (home *home) Get(c echo.Context) error {
	secret := home.envProvider.GetSecret()
	return c.JSON(http.StatusOK, secret)
}
