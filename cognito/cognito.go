package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/infani/awsgo/config/awsConfig"
)

func InitiateAuth(clientId string, username string, password string) (accessToken string, err error) {

	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(clientId),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	ctx := context.Background()
	output, err := client.InitiateAuth(ctx, input)
	if err != nil {
		return "", err
	}

	return *output.AuthenticationResult.AccessToken, nil
}
