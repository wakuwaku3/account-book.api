package event

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	publisher struct {
		provider Provider
		env      application.Env
	}
)

// NewPublisher はインスタンスを生成します
func NewPublisher(provider Provider, env application.Env) core.Publisher {
	return &publisher{provider, env}
}
func (t *publisher) Publish(name core.EventName, jsonMessage core.JSONMessage) error {
	messageBytes, err := json.Marshal(jsonMessage)
	if err != nil {
		return err
	}
	message := string(messageBytes)

	topics := t.env.GetAwsTopics()
	if topics == nil {
		return fmt.Errorf("指定されたイベント(%s)は見つかりません(message:%s)", name, message)
	}
	arn, ok := (*topics)[name]
	if !ok {
		return fmt.Errorf("指定されたイベント(%s)は見つかりません(message:%s)", name, message)
	}

	client := t.provider.GetSNSClient()
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(string(arn)),
	}

	result, _ := client.ListTopics(nil)
	for _, t := range result.Topics {
		fmt.Println(*t.TopicArn)
	}

	_, err = client.Publish(input)
	return err
}
