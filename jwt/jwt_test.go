package jwt

import "testing"

func TestCreateJWT(t *testing.T) {
	jwt, err := CreateJWT("./certs/private.pem")
	if err != nil {
		t.Error(err)
	}
	t.Log(jwt)
}
