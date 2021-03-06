package web

import (
	"reflect"

	"github.com/wakuwaku3/account-book.api/src/adapter/auth"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/tampopos/dijct"
)

const userKey = "user"

// Authenticate は Authenticateのミドルウェアです
func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			container := c.Get("di-container").(dijct.Container)
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			email := claims["email"].(string)
			userID := claims["nonce"].(string)
			claimsProvider := auth.NewClaimsProvider(email, userID, true)
			ifs := []reflect.Type{reflect.TypeOf((*core.ClaimsProvider)(nil)).Elem()}
			container.Register(claimsProvider, dijct.RegisterOptions{Interfaces: ifs})
			return next(c)
		}
	}
}
