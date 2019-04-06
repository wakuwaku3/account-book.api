package usecases

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	// AccountsQuery はアカウントのクエリです
	AccountsQuery interface {
		GetSignInInfo(email *string) (*SignInInfo, error)
		GetRefreshInfo(email *string) (*RefreshInfo, error)
		GetResetPasswordModelInfo(passwordResetToken *string) (*ResetPasswordModelInfo, error)
		GetResetPasswordInfo(passwordResetToken *string) (*ResetPasswordInfo, error)
		GetSignUpModelInfo(signUpToken *string) (*SignUpModelInfo, error)
	}
	// SignInInfo サインインのために必要な情報です
	SignInInfo struct {
		HashedPassword   string
		JwtClaims        application.JwtClaims
		JwtRefreshClaims application.JwtRefreshClaims
	}
	// RefreshInfo トークンリフレッシュのために必要な情報です
	RefreshInfo struct {
		AccountToken     string
		JwtClaims        application.JwtClaims
		JwtRefreshClaims application.JwtRefreshClaims
	}
	// ResetPasswordModelInfo はパスワードリセット画面表示のために必要な情報です
	ResetPasswordModelInfo struct {
		Email   string
		Expires time.Time
	}
	// SignUpModelInfo はサインアップ画面表示のために必要な情報です
	SignUpModelInfo struct {
		Email   string
		Expires time.Time
	}
	// ResetPasswordInfo はパスワードリセットのために必要な情報です
	ResetPasswordInfo struct {
		Email            string
		Expires          time.Time
		JwtClaims        application.JwtClaims
		JwtRefreshClaims application.JwtRefreshClaims
	}
	// TransactionsQuery はアカウントのクエリです
	TransactionsQuery interface {
		GetTransactions(args *GetTransactionsArgs) (*GetTransactionsResult, error)
		GetTransaction(id *string) (*GetTransactionResult, error)
	}
	// PlansQuery は計画のクエリです
	PlansQuery interface {
		GetPlans() (*GetPlansResult, error)
		GetPlan(id *string) (*GetPlanResult, error)
	}
	// DashboardQuery はダッシュボードのクエリです
	DashboardQuery interface {
		GetSummary(args *GetDashboardArgs) (*GetDashboardResult, error)
	}
	// ActualQuery は実績のクエリです
	ActualQuery interface {
		Get(args *GetActualArgs) (*GetActualResult, error)
		GetActualInfo(key *application.ActualKey) (*ActualInfo, error)
	}
)
