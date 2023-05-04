package iotCore

import (
	"reflect"
	"testing"
)

func Test_getEndpointAddress(t *testing.T) {
	got, err := getEndpointAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(got)
}

func TestGetCredentials(t *testing.T) {
	type args struct {
		certificateFiles CertificateFiles
		url              string
		thingName        string
	}
	tests := []struct {
		name    string
		args    args
		want    *Credentials
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				certificateFiles: CertificateFiles{
					CertPath:   "certs/certificate.pem.crt",
					KeyPath:    "certs/private.pem.key",
					CaCertPath: "certs/rootCA.crt",
				},
				url:       "https://ckoauhwbx2s36.credentials.iot.ap-northeast-1.amazonaws.com/role-aliases/hulkAssumeRole/credentials",
				thingName: "000020230213-1683181838271",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCredentials(tt.args.certificateFiles, tt.args.url, tt.args.thingName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
