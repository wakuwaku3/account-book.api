package queries

import (
	"github.com/wakuwaku3/account-book.api/src/domains"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"
)

type accounts struct {
	repos      domains.AccountsRepository
	usersRepos domains.UsersRepository
}

// NewAccounts はインスタンスを生成します
func NewAccounts(repos domains.AccountsRepository,
	usersRepos domains.UsersRepository) usecases.AccountsQuery {
	return &accounts{
		repos:      repos,
		usersRepos: usersRepos,
	}
}

func (t *accounts) GetSignInInfo(email *string) (*usecases.SignInInfo, error) {
	account, err := t.repos.Get(email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepos.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.SignInInfo{
		HashedPassword: account.HashedPassword,
		JwtClaims: domains.JwtClaims{
			Email:        *email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        *email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}

func (t *accounts) GetRefreshInfo(email *string) (*usecases.RefreshInfo, error) {
	account, err := t.repos.Get(email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepos.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.RefreshInfo{
		AccountToken: account.AccountToken,
		JwtClaims: domains.JwtClaims{
			Email:        *email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        *email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}

func (t *accounts) GetResetPasswordModelInfo(passwordResetToken *string) (*usecases.ResetPasswordModelInfo, error) {
	model, err := t.repos.GetPasswordResetToken(passwordResetToken)
	if err != nil {
		return nil, err
	}
	return &usecases.ResetPasswordModelInfo{
		Email:   model.Email,
		Expires: model.Expires,
	}, nil
}

func (t *accounts) GetResetPasswordInfo(passwordResetToken *string) (*usecases.ResetPasswordInfo, error) {
	model, err := t.repos.GetPasswordResetToken(passwordResetToken)
	if err != nil {
		return nil, err
	}
	account, err := t.repos.Get(&model.Email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepos.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.ResetPasswordInfo{
		Email:   model.Email,
		Expires: model.Expires,
		JwtClaims: domains.JwtClaims{
			Email:        model.Email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        model.Email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}

func (t *accounts) GetSignUpModelInfo(signUpToken *string) (*usecases.SignUpModelInfo, error) {
	model, err := t.repos.GetSignUpToken(signUpToken)
	if err != nil {
		return nil, err
	}
	return &usecases.SignUpModelInfo{
		Email:   model.Email,
		Expires: model.Expires,
	}, nil
}
