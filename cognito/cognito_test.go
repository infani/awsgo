package cognito

import (
	"testing"
)

func TestInitiateAuth(t *testing.T) {
	type args struct {
		clientId string
		username string
		password string
	}
	tests := []struct {
		name            string
		args            args
		wantAccessToken string
		wantErr         bool
	}{
		{
			name: "initiateAuth",
			args: args{
				clientId: "",
				username: "",
				password: "",
			},
			wantAccessToken: "",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, err := InitiateAuth(tt.args.clientId, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitiateAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccessToken != tt.wantAccessToken {
				t.Errorf("InitiateAuth() = %v, want %v", gotAccessToken, tt.wantAccessToken)
			}
		})
	}
}

func Test_decodeJWT(t *testing.T) {
	type args struct {
		jwksURL     string
		tokenString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "decodeJWT",
			args: args{
				jwksURL:     "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_BxG4FgrFm/.well-known/jwks.json",
				tokenString: "eyJraWQiOiJ3ZWh0TkFBQlJuR0RxbDVaT0E5dzZuU2JNU1FBemJSbDZhMTRZNFhXS1IwPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI2ZGQ3YzAxZC04YzVmLTRiZDItOTA1Yy00MjBmMWY5Yjk2YTAiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuYXAtbm9ydGhlYXN0LTEuYW1hem9uYXdzLmNvbVwvYXAtbm9ydGhlYXN0LTFfQnhHNEZnckZtIiwiY2xpZW50X2lkIjoiM25jZnVhM29oM3YxM2NqdDFrM25qa25jMHMiLCJvcmlnaW5fanRpIjoiZTljZTUzYmMtMTU3Ny00OGI2LTllZGQtOTQxODM2ZDAwYzcyIiwiZXZlbnRfaWQiOiIwMTM0MDkzMi0wZWU5LTRhY2MtOTMzNy04YTI1M2Q0YWVmNjIiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjY5OTcwNDI2LCJleHAiOjE2Njk5NzQwMjYsImlhdCI6MTY2OTk3MDQyNiwianRpIjoiMDllYjJhZGMtOTFiYy00N2U0LTllNDQtZWJiOGNkYTk0YTk1IiwidXNlcm5hbWUiOiI2ZGQ3YzAxZC04YzVmLTRiZDItOTA1Yy00MjBmMWY5Yjk2YTAifQ.FqmOwum-XCcJ5n6Wtw5m7oMs337fX_nJrG-2kjiOerNFx2yvYGisr5cgPf1g2p-NHzuAZo91SQyc55oFhGqfmOM8R-vrYAnSVCfZgpWHr8_eeRw348UBaeUgf_uvTCBRsEOOlaWu0IV5_iBYS8AyJwr70WLhNr63HTjIFXZpXO8v9R_mDVdqzeqLVTgn4odFmkYkpPvtnsT_OPBxHEjeX9HbOHZh5kFM9JsGuaXU7W0K-OvKaeolSe8Q5rUwTBwmyMlH0CqmYYDEIHidX77FmVQodCKYtqFaKGVN3VlKc6mRhfUL3Y1z5_Q0DikHXPcFZkaqCHfv-8nYVVjJJYviwg",
				},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeJWT(tt.args.jwksURL, tt.args.tokenString)
		})
	}
}
