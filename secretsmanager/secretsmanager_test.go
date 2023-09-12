package secretsmanager

import (
	"testing"
)

func TestGetSecretValue(t *testing.T) {
	secretName := "vortex/cloudfront"
	v, err := GetSecretValue(secretName)
	if err != nil {
		t.Error(err)
	}
	t.Log(v)
}

func TestGetSecretValueToMap(t *testing.T) {
	secretName := "vortex/cloudfront"
	v, err := GetSecretValueToMap(secretName)
	if err != nil {
		t.Error(err)
	}
	t.Log(v)
}
