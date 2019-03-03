package usecases

import (
	"github.com/wakuwaku3/account-book.api/src/0-enterprise-business-rules/models"
)

type (
	// Env は環境変数を取得します
	Env interface {
		Initialize() error
		GetCredentialsFilePath() string
		GetSecret() string
	}
	// UsersRepository は新ユーザーのリポジトリです
	UsersRepository interface {
		Get() (*[]models.User, error)
	}
)
