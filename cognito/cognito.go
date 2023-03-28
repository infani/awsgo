package cognito

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/infani/awsgo/config/awsConfig"
	"github.com/lestrrat-go/jwx/jwk"
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

func decodeJWT(jwksURL string, tokenString string) {

	set, err := jwk.Fetch(context.Background(), jwksURL)
	if err != nil {
		panic(err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		key, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		var rsaPublicKey rsa.PublicKey
		if err := key.Raw(&rsaPublicKey); err != nil {
			return nil, fmt.Errorf("failed to create public key: %v", err)
		}

		return &rsaPublicKey, nil
	})

	if err != nil {
		log.Println(err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["sub"])
	} else {
		fmt.Println(err)
	}
}
