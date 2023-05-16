package sqs

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gookit/goutil/dump"
)

const queueUrl = "https://sqs.ap-northeast-1.amazonaws.com/073903779593/awsgo"

func TestReceiveMessage(t *testing.T) {
	type args struct {
		queueUrl          string
		visibilityTimeout int32
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Message
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				queueUrl:          queueUrl,
				visibilityTimeout: 1,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReceiveMessage(tt.args.queueUrl, tt.args.visibilityTimeout)
			dump.P(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReceiveMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	err := SendMessage(queueUrl, "TestDeleteMessage")
	if err != nil {
		t.Error(err)
		return
	}
	msg, err := ReceiveMessage(queueUrl, 1)
	if err != nil {
		t.Error(err)
		return
	}
	dump.P(msg)
	err = DeleteMessage(queueUrl, msg.ReceiptHandle)
	if err != nil {
		t.Error(err)
		return
	}
}
