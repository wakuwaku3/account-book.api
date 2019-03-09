package usecases

import "github.com/wakuwaku3/account-book.api/src/domains"

type (
	// AccountsQuery はアカウントのクエリです
	AccountsQuery interface {
		GetSignInInfo(email *string) (*SignInInfo, error)
		GetRefreshInfo(email *string) (*RefreshInfo, error)
	}
	// SignInInfo サインインのために必要な情報です
	SignInInfo struct {
		HashedPassword   string
		JwtClaims        domains.JwtClaims
		JwtRefreshClaims domains.JwtRefreshClaims
	}
	// RefreshInfo トークンリフレッシュのために必要な情報です
	RefreshInfo struct {
		AccountToken     string
		JwtClaims        domains.JwtClaims
		JwtRefreshClaims domains.JwtRefreshClaims
	}
)
