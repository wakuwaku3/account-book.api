package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/domains"
)

type (
	// AccountsQuery はアカウントのクエリです
	AccountsQuery interface {
		GetSignInInfo(email *string) (*SignInInfo, error)
		GetRefreshInfo(email *string) (*RefreshInfo, error)
		GetResetPasswordModelInfo(passwordResetToken *string) (*ResetPasswordModelInfo, error)
		GetResetPasswordInfo(passwordResetToken *string) (*ResetPasswordInfo, error)
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
	// ResetPasswordModelInfo はパスワードリセット画面表示のために必要な情報です
	ResetPasswordModelInfo struct {
		Email   string
		Expires time.Time
	}
	// ResetPasswordInfo はパスワードリセットのために必要な情報です
	ResetPasswordInfo struct {
		Email            string
		Expires          time.Time
		JwtClaims        domains.JwtClaims
		JwtRefreshClaims domains.JwtRefreshClaims
	}
	// TransactionsQuery はアカウントのクエリです
	TransactionsQuery interface {
		GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error)
		GetTransaction(id *string) (*GetTransactionResult, error)
	}
)
