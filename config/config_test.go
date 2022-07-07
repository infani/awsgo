package config

import (
	"testing"
)

func Test_getService(t *testing.T) {
	SetByGoTest()
	type args struct {
		service string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "getService",
			args: args{
				service: "vivoreco",
			},
			want:    map[string]string{"stage": "dev"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getService(tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("getService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got["stage"] == tt.want["stage"]) {
				t.Errorf("getService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	SetByGoTest()
	tests := []struct {
		name string
	}{
		{
			name: "toString",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toString()
		})
	}
}
