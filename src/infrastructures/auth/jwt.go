package auth

import (
	"time"

	"github.com/google/uuid"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	j struct {
		env domains.Env
	}
	// CustomClaims は CustomClaimsです
	CustomClaims struct {
		UserID        string `json:"nonce"`
		UserName      string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		jwt.StandardClaims
	}
)

// NewJwt is create instance
func NewJwt(env domains.Env) domains.Jwt {
	return &j{env: env}
}
func (t *j) CreateToken(claims *domains.JwtClaims) (*string, error) {
	now := time.Now()
	cc := CustomClaims{
		UserID:        claims.UserID,
		UserName:      claims.UserName,
		Email:         claims.Email,
		EmailVerified: true,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "prj-account-book.appspot.com",
			Subject:   "https://prj-account-book.firebaseapp.com/",
			Audience:  "https://prj-account-book.firebaseapp.com/",
			ExpiresAt: now.AddDate(0, 0, 15).Unix(),
			NotBefore: now.Unix(),
			IssuedAt:  now.Unix(),
			Id:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	secret := t.env.GetJwtSecret()
	tokenString, err := token.SignedString(*secret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
func (t *j) Parse(token *string) (*domains.JwtClaims, error) {
	return &domains.JwtClaims{}, nil
}
