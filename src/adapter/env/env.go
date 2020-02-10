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
		awsTopics           *map[core.EventName]application.AwsTopicArn
		awsQueues           *map[core.QueueName]application.AwsQueueURL
	}
	awsTopic struct {
		Name string `json:"name"`
		Arn  string `json:"arn"`
	}
	awsQueue struct {
		Name string `json:"name"`
		URL  string `json:"url"`
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
	if err := t.setAwsTopics(); err != nil {
		return err
	}
	return t.setAwsQueues()
}
func (t *env) setAwsTopics() error {
	awsTopicsByte := []byte(os.Getenv("AWS_TOPICS"))
	var awsTopics []awsTopic
	if err := json.Unmarshal(awsTopicsByte, &awsTopics); err != nil {
		return err
	}
	awsTopicsMap := map[core.EventName]application.AwsTopicArn{}
	for _, awsTopic := range awsTopics {
		awsTopicsMap[core.EventName(awsTopic.Name)] = application.AwsTopicArn(awsTopic.Arn)
	}
	t.awsTopics = &awsTopicsMap
	return nil
}
func (t *env) setAwsQueues() error {
	awsTopicsByte := []byte(os.Getenv("AWS_QUEUES"))
	var awsQueues []awsQueue
	if err := json.Unmarshal(awsTopicsByte, &awsQueues); err != nil {
		return err
	}
	awsQueuesMap := map[core.QueueName]application.AwsQueueURL{}
	for _, awsQueue := range awsQueues {
		awsQueuesMap[core.QueueName(awsQueue.Name)] = application.AwsQueueURL(awsQueue.URL)
	}
	t.awsQueues = &awsQueuesMap
	return nil
}
func (t *env) GetCredentialsFilePath() *string {
	return t.credentialsFilePath
}
func (t *env) GetSendGridAPIKey() *string {
	return t.sendGridAPIKey
}
func (t *env) GetFrontEndURL() *string {
	return t.frontEndURL
}
func (t *env) GetPasswordHashedKey() *[]byte {
	return t.passwordHashedKey
}
func (t *env) GetJwtSecret() *[]byte {
	return t.jwtSecret
}
func (t *env) IsProduction() bool {
	return t.isProduction
}
func (t *env) GetAwsAccessKey() *string {
	return t.awsAccessKey
}
func (t *env) GetAwsSecretAccessKey() *string {
	return t.awsSecretAccessKey
}
func (t *env) GetAwsTopics() *map[core.EventName]application.AwsTopicArn {
	return t.awsTopics
}
func (t *env) GetAwsQueues() *map[core.QueueName]application.AwsQueueURL {
	return t.awsQueues
}
func (t *env) GetAllowOrigins() *[]string {
	if t.isProduction {
		return &[]string{*t.frontEndURL}
	}
	return &[]string{"*"}
}
func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
