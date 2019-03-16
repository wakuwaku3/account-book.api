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
		CreateRefreshToken(claims *JwtRefreshClaims) (*string, error)
		ParseRefreshToken(refreshToken *string) (*JwtRefreshClaims, error)
	}
	// JwtClaims はJwtTokenにうめこまれます
	JwtClaims struct {
		UserID       string
		UserName     string
		Email        string
		Culture      string
		UseStartDate time.Time
	}
	// JwtRefreshClaims はRefreshTokenにうめこまれます
	JwtRefreshClaims struct {
		UserID       string
		Email        string
		AccountToken string
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
		CreatePasswordResetToken(model *models.PasswordResetToken) (*string, error)
		CleanUp() error
		CleanUpByEmail(email *string) error
		GetPasswordResetToken(passwordResetToken *string) (*models.PasswordResetToken, error)
		SetPassword(email *string, hashedPassword *string) error
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
	// TransactionsRepository は取引のリポジトリです
	TransactionsRepository interface {
		Get(id *string) (*models.Transaction, error)
		GetByMonth(month *time.Time) (*[]models.Transaction, error)
		Create(model *models.Transaction) (*string, error)
		Update(id *string, model *models.Transaction) error
		Delete(id *string) error
	}
	// PlansRepository は計画のリポジトリです
	PlansRepository interface {
		Get() (*[]models.Plan, error)
		GetByID(id *string) (*models.Plan, error)
		Create(model *models.Plan) (*string, error)
		Update(id *string, model *models.Plan) error
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
