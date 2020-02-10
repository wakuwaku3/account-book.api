package event

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tampopos/dijct"
	"github.com/wakuwaku3/account-book.api/src/application"
	"github.com/wakuwaku3/account-book.api/src/enterprise/core"
)

type (
	subscriber struct {
		provider Provider
		router   Router
		env      application.Env
	}
	// Subscriber は Queue イベントを購読します
	Subscriber interface {
		Subscribe(container dijct.Container) error
	}
)

// NewSubscriber はインスタンスを生成します
func NewSubscriber(provider Provider, router Router, env application.Env) Subscriber {
	return &subscriber{provider, router, env}
}
func (t *subscriber) Subscribe(container dijct.Container) error {
	semaphore := make(chan bool, 10)
	queues := *t.env.GetAwsQueues()
	client := t.provider.GetSQSClient()
	return t.router.Each(func(name core.QueueName, handler Handler) error {
		url, ok := queues[name]
		if !ok {
			return fmt.Errorf("指定されたキュー(%s)のURLが見つかりません", name)
		}
		go t.receiveMessage(name, url, client, semaphore, container, handler)
		return nil
	})
}

func (t *subscriber) receiveMessage(
	name core.QueueName,
	url application.AwsQueueURL,
	client *sqs.SQS,
	semaphore chan bool,
	container dijct.Container,
	handler Handler,
) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(string(url)),
		// 一度に取得する最大メッセージ数。最大でも10まで。
		MaxNumberOfMessages: aws.Int64(10),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}
	for {
		resp, err := client.ReceiveMessage(params)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			continue
		}

		if len(resp.Messages) > 0 {
			fmt.Fprintf(os.Stdout, "ReceiveMessage name:%s, url:%s\n", name, url)
		}

		for _, m := range resp.Messages {
			semaphore <- true
			go t.handleMessage(url, semaphore, container, handler, m)
		}
	}
}

func (t *subscriber) handleMessage(
	url application.AwsQueueURL,
	semaphore chan bool,
	container dijct.Container,
	handler Handler,
	msg *sqs.Message,
) {
	defer t.postHandleMessage(url, semaphore, msg)

	if err := handler(container, msg.Body); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		t.sendFailure(url, msg.ReceiptHandle)
	} else if err := t.sendSuccess(url, msg.ReceiptHandle); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func (t *subscriber) postHandleMessage(
	url application.AwsQueueURL,
	semaphore chan bool,
	msg *sqs.Message,
) {
	err := recover()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		t.sendFailure(url, msg.ReceiptHandle)
	}
	<-semaphore
}

var visibilityTimeoutZero int64 = 0

func (t *subscriber) sendFailure(url application.AwsQueueURL, receiptHandle *string) {
	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(string(url)),
		ReceiptHandle:     receiptHandle,
		VisibilityTimeout: &visibilityTimeoutZero,
	}
	client := t.provider.GetSQSClient()

	if _, err := client.ChangeMessageVisibility(params); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
}

func (t *subscriber) sendSuccess(url application.AwsQueueURL, receiptHandle *string) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(string(url)),
		ReceiptHandle: receiptHandle,
	}
	client := t.provider.GetSQSClient()
	_, err := client.DeleteMessage(params)
	return err
}
