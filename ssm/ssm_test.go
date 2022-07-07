package ssm

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestGetParameter(t *testing.T) {
	type args struct {
		name *string
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
				name: aws.String("TestGetParameter"),
			},
			want:    "TestGetParameter",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetParameter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("GetParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
