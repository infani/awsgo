package iotCore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/infani/awsgo/config/awsConfig"
)

func getEndpointAddress() (string, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return "", err
	}

	client := iot.NewFromConfig(cfg)
	out, err := client.DescribeEndpoint(context.Background(), &iot.DescribeEndpointInput{
		EndpointType: aws.String("iot:CredentialProvider"),
	})
	if err != nil {
		return "", err
	}
	return *out.EndpointAddress, nil
}
