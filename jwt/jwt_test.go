package jwt

import "testing"

func TestCreateJWT(t *testing.T) {
	token, err := CreateJWT("./certs/private.pem.key")
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
	err = VerifyJWT("./certs/public.pem.key", token)
	if err != nil {
		t.Error(err)
	}
}
