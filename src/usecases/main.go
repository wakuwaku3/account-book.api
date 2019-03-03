package usecases

import "github.com/wakuwaku3/account-book.api/src/domains"

type (
	// AccountsQuery はアカウントのクエリです
	AccountsQuery interface {
		GetSignInInfo(email *string) (*SignInInfo, error)
	}
	// SignInInfo サインインのために必要な情報です
	SignInInfo struct {
		HashedPassword string
		JwtClaims      domains.JwtClaims
	}
)
