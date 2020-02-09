package auth

import (
	"fmt"

	"github.com/google/uuid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	j struct {
		env   application.Env
		clock core.Clock
	}
	customClaims struct {
		UserID        string `json:"nonce"`
		UserName      string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"locale"`
		jwt.StandardClaims
	}
	customRefreshClaims struct {
		UserID        string `json:"nonce"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		AccountToken  string `json:"https://prj-account-book.appspot.com//claim-types/user-name"`
		jwt.StandardClaims
	}
)

// NewJwt is create instance
func NewJwt(env application.Env, clock core.Clock) application.Jwt {
	return &j{env, clock}
}
func (t *j) CreateToken(claims *application.JwtClaims) (*string, error) {
	now := t.clock.Now()
	url := t.env.GetFrontEndURL()
	cc := customClaims{
		UserID:        claims.UserID,
		UserName:      claims.UserName,
		Email:         claims.Email,
		EmailVerified: true,
		Locale:        claims.Culture,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "prj-account-book.appspot.com",
			Subject:   *url,
			Audience:  *url,
			ExpiresAt: now.AddDate(0, 0, 1).Unix(),
			NotBefore: claims.UseStartDate.Unix(),
			IssuedAt:  now.Unix(),
			Id:        uuid.New().String(),
		},
	}
	return t.createToken(cc)
}
func (t *j) CreateRefreshToken(claims *application.JwtRefreshClaims) (*string, error) {
	now := t.clock.Now()
	url := t.env.GetFrontEndURL()
	cc := customRefreshClaims{
		UserID:        claims.UserID,
		AccountToken:  claims.AccountToken,
		Email:         claims.Email,
		EmailVerified: true,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "prj-account-book.appspot.com",
			Subject:   *url,
			Audience:  *url,
			ExpiresAt: now.AddDate(0, 0, 15).Unix(),
			NotBefore: now.Unix(),
			IssuedAt:  now.Unix(),
			Id:        uuid.New().String(),
		},
	}
	return t.createToken(cc)
}
func (t *j) ParseRefreshToken(refreshToken *string) (*application.JwtRefreshClaims, error) {
	var c customRefreshClaims
	secret := t.env.GetJwtSecret()
	jwt.ParseWithClaims(*refreshToken, &c, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return secret, nil
	})
	return &application.JwtRefreshClaims{
		AccountToken: c.AccountToken,
		Email:        c.Email,
		UserID:       c.UserID,
	}, nil
}
func (t *j) createToken(claims jwt.Claims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := t.env.GetJwtSecret()
	tokenString, err := token.SignedString(*secret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
