package awsConfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var cfg aws.Config

func init() {
	LoadAWSDefaultConfig()
}

func LoadAWSDefaultConfig() (aws.Config, error) {
	if cfg.Credentials != nil {
		credentials, err := cfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			return cfg, err
		}
		if credentials.Expired() {
			cfg, err = config.LoadDefaultConfig(context.TODO())
			if err != nil {
				return cfg, err
			}
		}
	} else {
		var err error
		cfg, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}
