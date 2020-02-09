package event

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wakuwaku3/account-book.api/src/application"
)

type (
	provider struct {
		env application.Env
		sns *sns.SNS
	}
	// Provider は sns へのアクセサを提供します
	Provider interface {
		Initialize() error
		GetClient() *sns.SNS
	}
)

// NewProvider はインスタンスを生成します
func NewProvider(env application.Env) Provider {
	return &provider{env, nil}
}
func (t *provider) Initialize() error {
	awsAccessKey := t.env.GetAwsAccessKey()
	if awsAccessKey == nil || *awsAccessKey == "" {
		return errors.New("AWS_ACCESS_KEY が設定されていません")
	}
	awsSecretAccessKey := t.env.GetAwsSecretAccessKey()
	if awsSecretAccessKey == nil || *awsSecretAccessKey == "" {
		return errors.New("AWS_SECRET_ACCESS_KEY が設定されていません")
	}
	creds := credentials.NewStaticCredentials(*awsAccessKey, *awsSecretAccessKey, "")
	session, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1"),
	})
	if err != nil {
		return err
	}
	t.sns = sns.New(session)
	return nil
}
func (t *provider) GetClient() *sns.SNS { return t.sns }
