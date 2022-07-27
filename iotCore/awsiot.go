package iotCore

import (
	"context"

	"github.com/infani/awsgo/config/awsConfig"

	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

type Client interface {
	Update(thing string, msg interface{}) error
	Get(thing string) ([]byte, error)
	Publish(topic string, msg interface{}) error
}

func Publish(topic string, payload []byte) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return err
	}

	in := &iotdataplane.PublishInput{Topic: &topic, Payload: payload}
	client := iotdataplane.NewFromConfig(cfg)
	_, err = client.Publish(context.Background(), in)
	return err
}
