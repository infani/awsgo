package iotCore

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/infani/awsgo/config/awsConfig"
)

// aws iot describe-endpoint --endpoint-type iot:CredentialProvider
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

type CertificateFiles struct {
	CertPath   string
	KeyPath    string
	CaCertPath string
}

type Credentials struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	Expiration      string `json:"expiration"`
}

// https://docs.aws.amazon.com/zh_cn/iot/latest/developerguide/authorizing-direct-aws.html
func GetCredentials(certificateFiles CertificateFiles, url string, thingName string) (*Credentials, error) {
	cert, err := tls.LoadX509KeyPair(certificateFiles.CertPath, certificateFiles.KeyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(certificateFiles.CaCertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			},
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-amzn-iot-thingname", thingName)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d", res.StatusCode)
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	m := struct {
		Credentials Credentials `json:"credentials"`
	}{}
	json.Unmarshal(bodyBytes, &m)
	return &m.Credentials, nil
}
