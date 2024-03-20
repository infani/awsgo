package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func CreateJWT(privateKey string) (jwt string, err error) {
	keyBytes, err := os.ReadFile(privateKey)
	if err != nil {
		return "", err
	}

	// Parse the private key
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Create the JWT token
	token := jwtv5.New(jwtv5.SigningMethodRS256)
	token.Claims = jwtv5.MapClaims{
		"sub": "your-subject",
		"aud": "your-audience",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// Sign the token with the private key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
