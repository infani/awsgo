package signer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/infani/awsgo/iotCore"
)

func TestSign(t *testing.T) {

	certificateFiles := iotCore.CertificateFiles{
		CertPath:   "../certs/certificate.pem.crt",
		KeyPath:    "../certs/private.pem.key",
		CaCertPath: "../certs/rootCA.crt",
	}
	certsUrl := "https://ckoauhwbx2s36.credentials.iot.ap-northeast-1.amazonaws.com/role-aliases/hulkAssumeRole/credentials"
	thingName := "000020230213-1683181838271"
	certs, _ := iotCore.GetCredentials(certificateFiles, certsUrl, thingName)
	t.Log(certs)

	url := fmt.Sprintf("https://iotconfig.dev.vortexcloud.com/things/%s/vsaas/system/general", thingName)
	t.Log(url)
	data := map[string]string{"name": "ben1"}
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
					AccessKeyID:     certs.AccessKeyId,
					SecretAccessKey: certs.SecretAccessKey,
					SessionToken:    certs.SessionToken,
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
			res, err := client.Do(tt.args.req)

			if err != nil {
				log.Println(err)
				return
			}
			defer res.Body.Close()

			log.Println("response Status:", res.Status)
		})
	}
}
