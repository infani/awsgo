package cloudmap

import (
	"reflect"
	"testing"
)

func Test_DiscoverInstances(t *testing.T) {
	type args struct {
		region    string
		namespace string
		service   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "discoverInstances",
			args: args{
				region:    "ap-northeast-1",
				namespace: "vivoreco",
				service:   "vivoreco",
			},
			want:    "vivoreco",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiscoverInstances(tt.args.namespace, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got[0].NamespaceName, tt.want) {
				t.Errorf("discoverInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}
