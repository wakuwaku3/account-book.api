package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	env struct {
		isProduction        bool
		credentialsFilePath *string
		jwtSecret           *[]byte
		passwordHashedKey   *[]byte
		sendGridAPIKey      *string
		frontEndURL         *string
	}
)

// NewEnv は環境変数管理用のインスタンスを生成します
func NewEnv() application.Env {
	return &env{}
}
func (env *env) Initialize() error {
	slice := make([]string, 0)
	if exists("./.user.env") {
		slice = append(slice, "./.user.env")
	}
	if exists("./.env") {
		slice = append(slice, "./.env")
	}
	err := godotenv.Load(slice...)
	env.isProduction = os.Getenv("APPLICATION_MODE") == "PRODUCTION"
	credentialsFilePath := os.Getenv("CREDENTIALS_FILE_PATH")
	env.credentialsFilePath = &credentialsFilePath
	passwordHashedKey := []byte(os.Getenv("PASSWORD_HASHED_KEY"))
	env.passwordHashedKey = &passwordHashedKey
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	env.jwtSecret = &jwtSecret
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
	env.sendGridAPIKey = &sendGridAPIKey
	frontEndURL := os.Getenv("FRONT_END_URL")
	env.frontEndURL = &frontEndURL
	return err
}
func (env *env) GetCredentialsFilePath() *string {
	return env.credentialsFilePath
}
func (env *env) GetSendGridAPIKey() *string {
	return env.sendGridAPIKey
}
func (env *env) GetFrontEndURL() *string {
	return env.frontEndURL
}
func (env *env) GetPasswordHashedKey() *[]byte {
	return env.passwordHashedKey
}
func (env *env) GetJwtSecret() *[]byte {
	return env.jwtSecret
}
func (env *env) IsProduction() bool {
	return env.isProduction
}
func (env *env) GetAllowOrigins() *[]string {
	if env.isProduction {
		return &[]string{*env.frontEndURL}
	}
	return &[]string{"*"}
}
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
