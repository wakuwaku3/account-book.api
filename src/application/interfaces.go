package application

import (
	"time"

	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
	"github.com/wakuwaku3/account-book.api/src/enterprise/models"
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
		GetAwsAccessKey() *string
		GetAwsSecretAccessKey() *string
		GetAwsTopics() *map[core.EventName]AwsTopicArn
		GetAwsQueues() *map[core.QueueName]AwsQueueURL
	}
	// AwsTopicArn は Topic の Arn です
	AwsTopicArn string
	// AwsQueueURL は Queue の URL です
	AwsQueueURL string
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
	// UsersRepository は新ユーザーのリポジトリです
	UsersRepository interface {
		Get(userID *string) (*models.User, error)
		GetByAuth() (*models.User, error)
	}
	// AccountsRepository はアカウントのリポジトリです
	AccountsRepository interface {
		Get(email *string) (*models.Account, error)
		CreatePasswordResetToken(model *models.PasswordResetToken) (*string, error)
		CleanUpPasswordResetToken() error
		CleanUpPasswordResetTokenByEmail(email *string) error
		GetPasswordResetToken(passwordResetToken *string) (*models.PasswordResetToken, error)
		SetPassword(email *string, hashedPassword *string) error
		CreateSignUpToken(model *models.SignUpToken) (*string, error)
		CleanUpSignUpToken() error
		GetSignUpToken(signUpToken *string) (*models.SignUpToken, error)
		CreateUserAndAccount(user *models.User, account *models.Account) (*models.User, *models.Account, error)
		GetByAuth() (*models.Account, error)
		Delete() error
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
	// UserExistingMail はアカウント作成済みお知らせメール送信サービスです
	UserExistingMail interface {
		Send(args *UserExistingMailSendArgs) error
	}
	// UserExistingMailSendArgs はアカウント作成済みお知らせメール送信用パラメータです
	UserExistingMailSendArgs struct {
		Email string
		Token string
	}
	// UserCreationMail はアカウント作成メール送信サービスです
	UserCreationMail interface {
		Send(args *UserCreationMailSendArgs) error
	}
	// UserCreationMailSendArgs はアカウント作成メール送信用パラメータです
	UserCreationMailSendArgs struct {
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
		GetByMonth(month *time.Time) (*[]models.Plan, error)
		GetByID(id *string) (*models.Plan, error)
		Create(model *models.Plan) (*string, error)
		Update(id *string, model *models.Plan) error
	}
	// DashboardRepository はダッシュボードのリポジトリです
	DashboardRepository interface {
		GetByID(id *string) (*models.Dashboard, error)
		ExistsClosedNext(id *string) error
		GetLatestClosedDashboard() (*models.Dashboard, error)
		GetOldestOpenDashboard() (*models.Dashboard, error)
		GetByMonth(month *time.Time) (*models.Dashboard, error)
		Create(month *time.Time) (*string, error)
		Approve(model *models.Dashboard) error
		CancelApprove(model *models.Dashboard) error
		GetActual(dashboardID *string, id *string) (*models.Actual, error)
		ExistsActual(dashboardID *string, planID *string) (*string, error)
		CreateActual(dashboardID *string, model *models.Actual) (*string, error)
		UpdateActual(dashboardID *string, id *string, model *models.Actual) error
		AdjustBalance(id *string, balance int) error
	}
)
