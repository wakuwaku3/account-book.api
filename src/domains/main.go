package domains

import (
	"time"

	"github.com/labstack/gommon/log"

	"github.com/wakuwaku3/account-book.api/src/domains/models"
)

type (
	// Env は環境変数を取得します
	Env interface {
		Initialize() error
		GetCredentialsFilePath() *string
		GetPasswordHashedKey() *[]byte
		GetJwtSecret() *[]byte
		GetSendGridAPIKey() *string
		GetFrontEndURL() *string
		IsProduction() bool
		GetAllowOrigins() *[]string
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
		CreatePasswordResetToken(email *string, expires *time.Time) (*string, error)
		CleanUp() error
		CleanUpByEmail(email string) error
	}
	// ResetPasswordMail はパスワード再設定メール送信サービスです
	ResetPasswordMail interface {
		Send(args *ResetPasswordMailSendArgs) error
	}
	// ResetPasswordMailSendArgs はパスワード再設定メール送信用パラメータです
	ResetPasswordMailSendArgs struct {
		Email string
		Token string
	}
)

// Try は成功するか上限回数まで処理を繰り返し行います
func Try(f func() error, limit int) error {
	count := 0
	for {
		err := f()
		if err == nil {
			return nil
		}
		count++
		if count >= limit {
			return err
		}
		log.Warn(err)
	}
}
