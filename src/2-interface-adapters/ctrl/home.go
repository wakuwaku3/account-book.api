package ctrl

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	home struct{}
	// Home is HomeController
	Home interface {
		Get(c echo.Context) error
	}
)

// NewHome is create instance.
func NewHome() Home {
	return &home{}
}
func (home *home) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello world")
}
