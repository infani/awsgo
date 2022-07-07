package cloudfront

import (
	"context"
	"time"
	"github.com/infani/awsgo/config/awsConfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func RemoveCache(distributionId string, path string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return err
	}
	input := awsCloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionId),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(time.Now().String()),
			Paths: &types.Paths{
				Items:    []string{path},
				Quantity: aws.Int32(1),
			},
		},
	}

	client := awsCloudfront.NewFromConfig(cfg)
	ctx := context.Background()
	_, err = client.CreateInvalidation(ctx, &input)
	return err
}
