package auth

import (
	"reflect"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/application"
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
			claimsProvider := NewClaimsProvider(email, userID)
			ifs := []reflect.Type{reflect.TypeOf((*application.ClaimsProvider)(nil)).Elem()}
			container.Register(claimsProvider, dijct.RegisterOptions{Interfaces: ifs})
			return next(c)
		}
	}
}
