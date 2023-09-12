package secretsmanager

import (
	"context"
	"encoding/json"
	"log"

	"github.com/infani/awsgo/config/awsConfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretValue(secretName string) (string, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Println(err)
		panic("configuration error, " + err.Error())
	}

	client := secretsmanager.NewFromConfig(cfg)
	out, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	})
	if err != nil {
		return "", err
	}
	return *out.SecretString, nil
}

func GetSecretValueToMap(secretName string) (secrets map[string]string, err error) {
	v, err := GetSecretValue(secretName)
	if err != nil {
		return secrets, err
	}
	err = json.Unmarshal([]byte(v), &secrets)
	return secrets, err
}
