package shadow

import (
	"context"
	"encoding/json"
	"log"

	"github.com/infani/awsgo/config/awsConfig"
	"github.com/infani/awsgo/iotCore"

	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

type iotCoreClient struct {
	shadowName string
}

func New(shadowName string) iotCore.Client {
	return &iotCoreClient{shadowName: shadowName}
}

func (c *iotCoreClient) Update(thing string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
		return err
	}
	// log.Println(msg)

	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := iotdataplane.NewFromConfig(cfg)
	_, err = client.UpdateThingShadow(context.Background(),
		&iotdataplane.UpdateThingShadowInput{ThingName: &thing, ShadowName: &c.shadowName, Payload: payload})
	return err
}

func (c *iotCoreClient) Get(thing string) ([]byte, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := iotdataplane.NewFromConfig(cfg)
	out, err := client.GetThingShadow(context.Background(), &iotdataplane.GetThingShadowInput{ThingName: &thing, ShadowName: &c.shadowName})

	if err != nil {
		return nil, err
	}

	return out.Payload, err
}

func (c *iotCoreClient) Publish(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
		return err
	}
	return iotCore.Publish(topic, payload)
}
