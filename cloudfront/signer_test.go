package cloudfront

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/infani/awsgo/s3"
	"github.com/infani/awsgo/secretsmanager"
)

const bucket = "hulk-devices-push-storage-the-greate-one-vsaas-vortex-dev"
const cloudfrontDomain = "cloudfront.dev.vortexcloud.com"

func getCloudfront() (*cloudfront, error) {
	keyID := "K8ASH4B6QLZS9"
	privateKey, err := secretsmanager.GetSecretValue("/vortex/infra/cloudfront/privatekey")
	if err != nil {
		return nil, err
	}
	return New(keyID, privateKey)
}

func TestSignUrl(t *testing.T) {
	c, err := getCloudfront()
	if err != nil {
		t.Error(err)
		return
	}
	const key = "playback.m3u8"
	s3.PutFile(bucket, key, bytes.NewReader([]byte("#EXTM3U")))
	rawURL := url.URL{
		Scheme: "https",
		Host:   cloudfrontDomain,
		Path:   key,
	}
	type args struct {
		url    string
		expire time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				url:    rawURL.String(),
				expire: time.Now().Add(time.Hour),
			},
			want: "Signature",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.SignUrl(tt.args.url, tt.args.expire)
			if !strings.Contains(got, tt.want) {
				t.Errorf("SignUrl() = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}

// curl -H "Origin: https://config.CloudfrontDomain" https://config.CloudfrontDomain/playback.m3u8 --cookie "CloudFront-Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cHM6Ly9jbG91ZGZyb250Lmh1bGsudnNhYXMudml2b2Nsb3VkLmNvbS8qIiwiQ29uZGl0aW9uIjp7IkRhdGVMZXNzVGhhbiI6eyJBV1M6RXBvY2hUaW1lIjoxNjc0NTMxMDUwfX19XX0_;CloudFront-Signature=rYndCbE1bZTbC-dK~d~6Mc2-ZCITs8ETBnhnW1aEsgUY2nSft0L5atIF4p9jDgB9Q8erq1ygsBaW5R2UALHuGAgygf50DcpxylAU-BHhHzoLPL5YCojyRzDF3VEB2lxO0oS4-0qpMSWJ2IxR5d4ZEgUqm6v9DDT9jhvqSIkDW42ZJXktUv1gm1u7KcfCM7MmOxBt4kDJVYSOwprdL3P997ZQZYYCUF8qeuSY76j0B1Cuv0VqnFLaWXEicctK8NtIY~ENBPGwTVVdxKVoIYQt3QUXmdn6Us-U1DXeb~rEhLBAd8jhY4wmAmpAxMmuy~zSObEMCIl8hg2Y~ek7l-knDQ__;CloudFront-Key-Pair-Id=K1TADGJ66R8YEU;"
func TestSignCookie(t *testing.T) {
	c, err := getCloudfront()
	if err != nil {
		t.Error(err)
		return
	}
	rawURL := url.URL{
		Scheme: "https",
		Host:   cloudfrontDomain,
		Path:   "*",
	}
	type args struct {
		url    string
		expire time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				url:    rawURL.String(),
				expire: time.Now().Add(time.Hour * 24 * 365),
			},
			want:    "CloudFront-Policy",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.SignCookie(tt.args.url, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got[0].Name, tt.want) {
				t.Errorf("SignCookie() = %v, want %v", got, tt.want)
			}
			t.Log("Sign cookie:", rawURL.String())
			fmt.Printf("cookies: \"%s=%s;%s=%s;%s=%s;\"\n", got[0].Name, got[0].Value, got[1].Name, got[1].Value, got[2].Name, got[2].Value)
			fmt.Printf("document.cookie=\"%s=%s;domain=.vivocloud.com;samesite=none;secure\"\n", got[0].Name, got[0].Value)
			fmt.Printf("document.cookie=\"%s=%s;domain=.vivocloud.com;samesite=none;secure\"\n", got[1].Name, got[1].Value)
			fmt.Printf("document.cookie=\"%s=%s;domain=.vivocloud.com;samesite=none;secure\"\n", got[2].Name, got[2].Value)
		})
	}
}
