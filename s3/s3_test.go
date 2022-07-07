package s3

import (
	"bytes"
	"strings"
	"testing"
)

const (
	bucket = "vivoreco-go-test-bucket"
)

func TestGetFile(t *testing.T) {
	const key = "thingName/recbackup/2022/01/19/02/playback.m3u8"
	PutFile(bucket, key, bytes.NewReader([]byte("#EXTM3U")))
	type args struct {
		bucket string
		key    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "test", args: args{bucket: bucket, key: key}, want: "#EXTM3U", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFile(tt.args.bucket, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(string(got), tt.want) {
				t.Errorf("GetFile() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestPutFile(t *testing.T) {
	type args struct {
		bucket string
		key    string
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test TestPutFile", args: args{bucket: bucket, key: "TestPutFile", data: []byte("TestPutFile")}, wantErr: false},
		{name: "test TestPutFile.m3u8", args: args{bucket: bucket, key: "TestPutFile.m3u8", data: []byte("TestPutFile.m3u8")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutFile(tt.args.bucket, tt.args.key, bytes.NewReader(tt.args.data)); (err != nil) != tt.wantErr {
				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	type args struct {
		bucket string
		key    string
		newKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test", args: args{bucket: bucket, key: "thingName/recbackup/2022/01/25/09/48733.ts", newKey: "thingName/archive/TestGenM3u8FromBackup/2022/01/25/09/48733.ts"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFile(tt.args.bucket, tt.args.key, tt.args.newKey); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
