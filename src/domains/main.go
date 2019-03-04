package domains

import (
	"github.com/wakuwaku3/account-book.api/src/domains/models"
)

type (
	// Env は環境変数を取得します
	Env interface {
		Initialize() error
		GetCredentialsFilePath() *string
		GetPasswordHashedKey() *[]byte
		GetJwtSecret() *[]byte
	}
	// Crypt はハッシュ化のサービスです
	Crypt interface {
		Hash(text *string) *string
	}
	// Jwt はJwtのサービスです
	Jwt interface {
		CreateToken(claims *JwtClaims) (*string, error)
		Parse(token *string) (*JwtClaims, error)
	}
	// JwtClaims はJwtTokenにうめこまれます
	JwtClaims struct {
		UserID   string
		UserName string
		Email    string
	}
	// ClaimsProvider は Claimsを取得します
	ClaimsProvider interface {
		GetUserID() *string
		GetEmail() *string
	}
	// UsersRepository は新ユーザーのリポジトリです
	UsersRepository interface {
		Get(userID *string) (*models.User, error)
	}
	// AccountsRepository はアカウントのリポジトリです
	AccountsRepository interface {
		Get(email *string) (*models.Account, error)
	}
)
