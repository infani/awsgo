package cloudmap

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
	"github.com/infani/awsgo/config/awsConfig"
)

func DiscoverInstances(namespace string, service string) ([]types.HttpInstanceSummary, error) {

	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := servicediscovery.NewFromConfig(cfg)

	input := &servicediscovery.DiscoverInstancesInput{
		NamespaceName: aws.String(namespace),
		ServiceName:   aws.String(service),
	}

	ctx := context.Background()
	output, err := client.DiscoverInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.Instances, nil
}
