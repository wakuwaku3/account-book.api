package ctrl

import (
	"log"
	"net/http"

	"github.com/wakuwaku3/account-book.api/src/1-application-business-rules/usecases"

	"github.com/labstack/echo"
)

type (
	home struct {
		usersRepository usecases.UsersRepository
	}
	// Home is HomeController
	Home interface {
		Get(c echo.Context) error
	}
)

// NewHome is create instance.
func NewHome(usersRepository usecases.UsersRepository) Home {
	return &home{usersRepository: usersRepository}
}
func (home *home) Get(c echo.Context) error {
	users, err := home.usersRepository.Get()
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, *users)
}
