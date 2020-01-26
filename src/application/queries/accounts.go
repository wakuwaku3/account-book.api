package queries

import (
	"github.com/wakuwaku3/account-book.api/src/application"

	"github.com/wakuwaku3/account-book.api/src/application/usecases"
)

type accounts struct {
	repos      application.AccountsRepository
	usersRepos application.UsersRepository
}

// NewAccounts はインスタンスを生成します
func NewAccounts(repos application.AccountsRepository,
	usersRepos application.UsersRepository) usecases.AccountsQuery {
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
		JwtClaims: application.JwtClaims{
			Email:        *email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: application.JwtRefreshClaims{
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
		JwtClaims: application.JwtClaims{
			Email:        *email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: application.JwtRefreshClaims{
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
		JwtClaims: application.JwtClaims{
			Email:        model.Email,
			UserID:       account.UserID,
			UserName:     user.UserName,
			Culture:      user.Culture,
			UseStartDate: user.UseStartDate,
		},
		JwtRefreshClaims: application.JwtRefreshClaims{
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
func (t *accounts) GetQuitInfo() (*usecases.QuitInfo, error) {
	account, err := t.repos.GetByAuth()
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepos.GetByAuth()
	if err != nil {
		return nil, err
	}
	return &usecases.QuitInfo{
		HashedPassword: account.HashedPassword,
		UserName:       user.UserName,
	}, nil

}
