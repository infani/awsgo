package signer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	url := fmt.Sprintf("https://iotconfigdevice.dev.vortexcloud.com/things/%s/vsaas/system/general", thingName)
	payload := fmt.Sprintf("{\"info\":{ \"name\":\"%s\", \"city\":\"%s\" } }", "ben", "AsiaTaipei")
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
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
			if err := Sign(tt.args.req, tt.args.credentials, "ap-northeast-1"); (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
			client := &http.Client{}
			res, err := client.Do(tt.args.req)

			if err != nil {
				log.Println(err)
				return
			}
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			log.Println(string(body))
			log.Println("response Status:", res.Status)
		})
	}
}
