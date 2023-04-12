package ssm

import (
	"context"
	"log"

	"github.com/infani/awsgo/config/awsConfig"

	awsssm "github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetParameter(name *string) (string, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Println(err)
		panic("configuration error, " + err.Error())
	}

	client := awsssm.NewFromConfig(cfg)
	b := true
	out, err := client.GetParameter(context.Background(), &awsssm.GetParameterInput{Name: name, WithDecryption: &b})
	if err != nil {
		return "", err
	}
	return *out.Parameter.Value, nil
}
