package cloudfront

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
)

type cloudfront struct {
	keyID      string
	privateKey *rsa.PrivateKey
}

func New(keyID, privateKey string) (*cloudfront, error) {
	keyReader := strings.NewReader(privateKey)
	key, err := sign.LoadPEMPrivKey(keyReader)
	if err != nil {
		return nil, err
	}
	return &cloudfront{keyID: keyID, privateKey: key}, nil
}

func (h *cloudfront) SignUrl(url string, expire time.Time) string {
	signer := sign.NewURLSigner(h.keyID, h.privateKey)
	signedURL, err := signer.Sign(url, expire)
	if err != nil {
		log.Fatalf("Failed to sign url, err: %s\n", err.Error())
		return ""
	}
	return signedURL
}

func (h *cloudfront) SignCookie(url string, expire time.Time) ([]*http.Cookie, error) {
	signer := sign.NewCookieSigner(h.keyID, h.privateKey)
	signedCookie, err := signer.Sign(url, expire)
	if err != nil {
		log.Fatalf("Failed to sign cookie, err: %s\n", err.Error())
		return nil, err
	}
	return signedCookie, nil
}

func (h *cloudfront) GenCookie(url string, expire time.Time) (string, error) {
	got, err := h.SignCookie(url, time.Now().Add(time.Hour*24))
	if err != nil {
		return "", err
	}
	cookie := fmt.Sprintf("%s=%s;%s=%s;%s=%s;", got[0].Name, got[0].Value, got[1].Name, got[1].Value, got[2].Name, got[2].Value)
	return cookie, err
}
