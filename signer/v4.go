package signer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

//https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/sig-v4-header-based-auth.html
//https://docs.aws.amazon.com/zh_cn/apigateway/latest/developerguide/api-gateway-iam-policy-examples-for-api-execution.html
func Sign(req *http.Request, credentials aws.Credentials) error {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
	}

	signer := v4.NewSigner()
	h := sha256.New()
	body, _ := req.GetBody()
	io.Copy(h, body)
	payloadHash := hex.EncodeToString(h.Sum(nil))
	err = signer.SignHTTP(context.Background(), credentials, req, payloadHash, "execute-api", cfg.Region, time.Now())
	// log.Println(req.Header)
	if err != nil {
		log.Println(err)
	}
	return nil
}
