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

func Sign(req *http.Request, credentials aws.Credentials) error {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
	}

	signer := v4.NewSigner()
	// hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	//hash body
	h := sha256.New()
	body, _ := req.GetBody()
	io.Copy(h, body)
	payloadHash := hex.EncodeToString(h.Sum(nil))
	err = signer.SignHTTP(context.Background(), credentials, req, payloadHash, req.Host, cfg.Region, time.Now())
	if err != nil {
		log.Println(err)
	}
	return nil
}
