package env

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	env struct {
		isProduction        bool
		credentialsFilePath *string
		jwtSecret           *[]byte
		passwordHashedKey   *[]byte
		sendGridAPIKey      *string
		frontEndURL         *string
		awsAccessKey        *string
		awsSecretAccessKey  *string
		awsTopics           *map[core.EventName]string
	}
	awsTopics struct {
		Name string `json:"name"`
		Arn  string `json:"arn"`
	}
)

// NewEnv は環境変数管理用のインスタンスを生成します
func NewEnv() application.Env {
	return &env{}
}
func (t *env) Initialize() error {
	slice := make([]string, 0)
	if exists("./.user.env") {
		slice = append(slice, "./.user.env")
	}
	if exists("./.env") {
		slice = append(slice, "./.env")
	}
	if err := godotenv.Load(slice...); err != nil {
		return err
	}
	t.isProduction = os.Getenv("APPLICATION_MODE") == "PRODUCTION"
	credentialsFilePath := os.Getenv("CREDENTIALS_FILE_PATH")
	t.credentialsFilePath = &credentialsFilePath
	passwordHashedKey := []byte(os.Getenv("PASSWORD_HASHED_KEY"))
	t.passwordHashedKey = &passwordHashedKey
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	t.jwtSecret = &jwtSecret
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
	t.sendGridAPIKey = &sendGridAPIKey
	frontEndURL := os.Getenv("FRONT_END_URL")
	t.frontEndURL = &frontEndURL
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	t.awsAccessKey = &awsAccessKey
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	t.awsSecretAccessKey = &awsSecretAccessKey
	return t.setAwsTopics()
}
func (env *env) setAwsTopics() error {
	awsTopicsByte := []byte(os.Getenv("AWS_TOPICS"))
	var awsTopics []awsTopics
	if err := json.Unmarshal(awsTopicsByte, &awsTopics); err != nil {
		return err
	}
	awsTopicsMap := map[core.EventName]string{}
	for _, awsTopic := range awsTopics {
		awsTopicsMap[core.EventName(awsTopic.Name)] = awsTopic.Arn
	}
	env.awsTopics = &awsTopicsMap
	return nil
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
func (env *env) GetAwsAccessKey() *string {
	return env.awsAccessKey
}
func (env *env) GetAwsSecretAccessKey() *string {
	return env.awsSecretAccessKey
}
func (env *env) GetAwsTopics() *map[core.EventName]string {
	return env.awsTopics
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
