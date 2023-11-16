package iot

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/infani/awsgo/config/awsConfig"
)

func ListThings(ctx context.Context) (thingNames []string, err error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	cli := iot.NewFromConfig(cfg)
	var nextToken *string
	for {
		out, err := cli.ListThings(ctx, &iot.ListThingsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, thing := range out.Things {
			thingNames = append(thingNames, *thing.ThingName)
		}
		if out.NextToken == nil {
			break
		} else {
			nextToken = out.NextToken
		}
	}

	return thingNames, nil
}
