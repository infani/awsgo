package mqtt

import (
	"context"
	"reflect"
	"testing"
	"time"
)

// https://ap-northeast-1.console.aws.amazon.com/iot/home?region=ap-northeast-1#/thing/mqtt_test
const (
	server     = "tcps://a2vbbdorkkxq7n-ats.iot.ap-northeast-1.amazonaws.com:8883"
	certFile   = "./certs/mqtt_test.cert.pem"
	keyFile    = "./certs/mqtt_test.private.key"
	rootCaFile = "./certs/rootCA.crt"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts ClientOptions
	}
	tests := []struct {
		name            string
		args            args
		wantIsConnected bool
		wantErr         bool
	}{
		{
			name: "wrong host",
			args: args{
				opts: ClientOptions{
					Server:    "tcps://nohost:8883",
					TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
				},
			},
			wantErr: true,
		},
		{
			name: "wrong port",
			args: args{
				opts: ClientOptions{
					Server:    server[:len(server)-1],
					TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				opts: ClientOptions{
					Server:    server,
					TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
				},
			},
			wantIsConnected: true,
			wantErr:         false,
		},
		{
			name: "withClientID",
			args: args{
				opts: ClientOptions{
					Server:    server,
					TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
					ClientID:  "clientID",
				},
			},
			wantIsConnected: true,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && !reflect.DeepEqual(got.IsConnected(), tt.wantIsConnected) {
				t.Errorf("IsConnected() = %v, want %v", got.IsConnected(), tt.wantIsConnected)
			}
		})
	}
}

func Test_client_Publish(t *testing.T) {
	type args struct {
		topic   string
		payload interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				topic:   "Device/Test_client_Publish/sub/peopletrack/webrtc/",
				payload: "Test_client_Publish",
			},
			wantErr: false,
		}, {
			name: "noAuth",
			args: args{
				topic:   "noAuth",
				payload: "Test_client_Publish",
			},
			wantErr: false,
		},
	}
	cli, err := NewClient(ClientOptions{
		Server:    server,
		TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
	})
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.Publish(tt.args.topic, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("client.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if !cli.IsConnected() {
		t.Log("cli.IsConnected()", cli.IsConnected())
	}
	time.Sleep(1 * time.Second)
	if cli.IsConnected() {
		t.Log("cli.IsConnected()", cli.IsConnected())
	}
}

func Test_client_Subscribe(t *testing.T) {
	ctx := context.Background()
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				topic: "Device/Test_client_Subscribe/sub/peopletrack/webrtc/",
			},
			wantErr: false,
		}, {
			name: "noAuth",
			args: args{
				topic: "noAuth",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := NewClient(ClientOptions{
				Server:    server,
				TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
			})
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
				return
			}
			obs, err := cli.Subscribe(ctx, tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			cli.Publish(tt.args.topic, "Test_client_Subscribe")
			stream := obs.First().Observe()
			msg := <-stream
			if msg.Error() {
				t.Errorf("client.Subscribe() error = %v", msg.E)
			} else {
				t.Log(string(msg.V.([]byte)))
			}
			cli.Close()
		})
	}
}

func Test_client_SubscribeReturnMessage(t *testing.T) {
	ctx := context.Background()
	type args struct {
		topic string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				topic: "Device/Test_client_Subscribe/sub/peopletrack/webrtc/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := NewClient(ClientOptions{
				Server:    server,
				TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
			})
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
				return
			}
			obs, err := cli.SubscribeReturnMessage(ctx, tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			cli.Publish(tt.args.topic, "Test_client_Subscribe")
			stream := obs.First().Observe()
			msg := <-stream
			if msg.Error() {
				t.Errorf("client.Subscribe() error = %v", msg.E)
			} else {
				t.Log(msg.V)
			}
			cli.Close()
		})
	}
}

// use "netstat -tunp|grep 8883" and "sudo tcpkill host ip"
func Test_Connection_Lost(t *testing.T) {
	ctx := context.Background()
	topic := "Device/Test_Connection_Lost/sub/peopletrack/webrtc/"
	cli, err := NewClient(ClientOptions{
		Server:    server,
		TLSConfig: GetTLSConfigFromFile(certFile, keyFile, rootCaFile),
	})
	if err != nil {
		t.Errorf("NewClient() error = %v", err)
		return
	}
	obs, err := cli.Subscribe(ctx, topic)
	if err != nil {
		t.Errorf("client.Subscribe() error = %v", err)
		return
	}
	cli.Publish(topic, "Test_Connection_Lost")
	stream := obs.Observe()
	go func() {
		time.Sleep(time.Second * 5)
		cli.Close()
	}()
	for msg := range stream {
		if msg.Error() {
			t.Errorf("client.Subscribe() error = %v", msg.E)
		} else {
			t.Log(string(msg.V.([]byte)))
		}
	}
}
