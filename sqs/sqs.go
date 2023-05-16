package sqs

import (
	"context"

	"github.com/infani/awsgo/config/awsConfig"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func SendMessage(queueUrl string, message string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awssqs.NewFromConfig(cfg)
	_, err = client.SendMessage(context.Background(), &awssqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &message,
	})
	return err
}

// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs#Client.ReceiveMessage
func ReceiveMessage(queueUrl string, visibilityTimeout int32) (*types.Message, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := awssqs.NewFromConfig(cfg)
	out, err := client.ReceiveMessage(context.Background(), &awssqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   visibilityTimeout,
	})
	if err != nil {
		return nil, err
	}
	if len(out.Messages) == 0 {
		return nil, nil
	}
	return &out.Messages[0], nil
}

func DeleteMessage(queueUrl string, receiptHandle *string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := awssqs.NewFromConfig(cfg)
	_, err = client.DeleteMessage(context.Background(), &awssqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: receiptHandle,
	})
	return err
}
