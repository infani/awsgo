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
