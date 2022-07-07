package cloudfront

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/infani/awsgo/config"
)

func getKeyReader() *strings.Reader {
	privateKey := ""
	if config.PrivateKey != "" {
		privateKey = config.PrivateKey
	}

	return strings.NewReader(privateKey)
}

func SignUrl(url string, expire time.Time) string {
	keyReader := getKeyReader()
	if keyReader == nil {
		log.Println("Failed to get key reader")
		return ""
	}
	privateKey, err := sign.LoadPEMPrivKey(keyReader)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	signer := sign.NewURLSigner(config.KeyID, privateKey)
	signedURL, err := signer.Sign(url, expire)
	if err != nil {
		log.Fatalf("Failed to sign url, err: %s\n", err.Error())
		return ""
	}
	return signedURL
}

func SignCookie(url string, expire time.Time) ([]*http.Cookie, error) {
	keyReader := getKeyReader()
	privateKey, err := sign.LoadPEMPrivKey(keyReader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	signer := sign.NewCookieSigner(config.KeyID, privateKey)
	signedCookie, err := signer.Sign(url, expire)
	if err != nil {
		log.Fatalf("Failed to sign cookie, err: %s\n", err.Error())
		return nil, err
	}
	return signedCookie, nil
}

func GenCookie(url string, expire time.Time) (string, error) {
	got, err := SignCookie(url, time.Now().Add(time.Hour*24))
	if err != nil {
		return "", err
	}
	cookie := fmt.Sprintf("%s=%s;%s=%s;%s=%s;", got[0].Name, got[0].Value, got[1].Name, got[1].Value, got[2].Name, got[2].Value)
	return cookie, err
}
