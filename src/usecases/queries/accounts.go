package queries

import (
	"github.com/wakuwaku3/account-book.api/src/domains"

	"github.com/wakuwaku3/account-book.api/src/usecases"
)

type accounts struct {
	repository      domains.AccountsRepository
	usersRepository domains.UsersRepository
}

// NewAccounts はインスタンスを生成します
func NewAccounts(repository domains.AccountsRepository,
	usersRepository domains.UsersRepository) usecases.AccountsQuery {
	return &accounts{
		repository:      repository,
		usersRepository: usersRepository,
	}
}

func (t *accounts) GetSignInInfo(email *string) (*usecases.SignInInfo, error) {
	account, err := t.repository.Get(email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepository.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.SignInInfo{
		HashedPassword: account.HashedPassword,
		JwtClaims: domains.JwtClaims{
			Email:    *email,
			UserID:   account.UserID,
			UserName: user.UserName,
			Culture:  user.Culture,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        *email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}

func (t *accounts) GetRefreshInfo(email *string) (*usecases.RefreshInfo, error) {
	account, err := t.repository.Get(email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepository.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.RefreshInfo{
		AccountToken: account.AccountToken,
		JwtClaims: domains.JwtClaims{
			Email:    *email,
			UserID:   account.UserID,
			UserName: user.UserName,
			Culture:  user.Culture,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        *email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}

func (t *accounts) GetResetPasswordModelInfo(passwordResetToken *string) (*usecases.ResetPasswordModelInfo, error) {
	model, err := t.repository.GetPasswordResetToken(passwordResetToken)
	if err != nil {
		return nil, err
	}
	return &usecases.ResetPasswordModelInfo{
		Email:   model.Email,
		Expires: model.Expires,
	}, nil
}

func (t *accounts) GetResetPasswordInfo(passwordResetToken *string) (*usecases.ResetPasswordInfo, error) {
	model, err := t.repository.GetPasswordResetToken(passwordResetToken)
	if err != nil {
		return nil, err
	}
	account, err := t.repository.Get(&model.Email)
	if err != nil {
		return nil, err
	}
	user, err := t.usersRepository.Get(&account.UserID)
	if err != nil {
		return nil, err
	}
	return &usecases.ResetPasswordInfo{
		Email:   model.Email,
		Expires: model.Expires,
		JwtClaims: domains.JwtClaims{
			Email:    model.Email,
			UserID:   account.UserID,
			UserName: user.UserName,
			Culture:  user.Culture,
		},
		JwtRefreshClaims: domains.JwtRefreshClaims{
			Email:        model.Email,
			UserID:       account.UserID,
			AccountToken: account.AccountToken,
		},
	}, nil
}
