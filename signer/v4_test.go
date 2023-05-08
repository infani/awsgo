package signer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestSign(t *testing.T) {
	url := "https://foo.amazonaws.com/bar"
	data := map[string]string{"firstname": "John", "lastname": "Doe"}
	jsonValue, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	type args struct {
		req         *http.Request
		credentials aws.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				req: req,
				credentials: aws.Credentials{
					AccessKeyID:     "ASIARCNICF4ETRDHRZVL",
					SecretAccessKey: "eNrW9c52JIvbeSArWw/t0efL5pICKQ+ufJdWCZbQ",
					SessionToken:    "XXX",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Sign(tt.args.req, tt.args.credentials); (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
			client := &http.Client{}
			res ,err := client.Do(tt.args.req)

			if err != nil {
				log.Println(err)
				return
			}
			defer res.Body.Close()
		
			log.Println("response Status:", res.Status)
		})
	}
}
