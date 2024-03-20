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

func VerifyJWT(publicKey string, tokenString string) (err error) {
	keyBytes, err := os.ReadFile(publicKey)
	if err != nil {
		return err
	}
	// Parse the public key
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	// Parse the token
	token, err := jwtv5.Parse(tokenString, func(token *jwtv5.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return err
	}

	// Verify the token
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// Token is valid
	return nil
}

func VerifyJWTByCert(cert string, tokenString string) (err error) {
	keyBytes, err := os.ReadFile(cert)
	if err != nil {
		return err
	}
	// Parse the public key
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode PEM block containing public key")
	}
	certData, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	// Parse the token
	token, err := jwtv5.Parse(tokenString, func(token *jwtv5.Token) (interface{}, error) {
		return certData.PublicKey, nil
	})
	if err != nil {
		return err
	}

	// Verify the token
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}