package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/reactivex/rxgo/v2"
)

type Client interface {
	Publish(topic string, payload interface{}) error
	Subscribe(ctx context.Context, topic string) (rxgo.Observable, error)
	SubscribeReturnMessage(ctx context.Context, topic string) (rxgo.Observable, error)
	IsConnected() bool
	Close()
}

type client struct {
	pahoMqttCli   pahoMqtt.Client
	subscriberMap sync.Map
	mu sync.Mutex
}

type ClientOptions struct {
	Server    string // tcp://host:port
	TLSConfig *ClientOptions_TLSConfig
	ClientID  string
}

type ClientOptions_TLSConfig struct {
	Cert   []byte
	Key    []byte
	RootCa []byte
}

func GetTLSConfigFromFile(certFile, keyFile, rootCaFile string) *ClientOptions_TLSConfig {
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil
	}
	key, _ := ioutil.ReadFile(keyFile)
	rootCa, _ := ioutil.ReadFile(rootCaFile)
	return &ClientOptions_TLSConfig{
		Cert:   cert,
		Key:    key,
		RootCa: rootCa,
	}
}

func NewClient(opts ClientOptions) (Client, error) {
	if opts.TLSConfig == nil {
		return nil, fmt.Errorf("tls config is required")
	}
	cert, err := tls.X509KeyPair(opts.TLSConfig.Cert, opts.TLSConfig.Key)
	if err != nil {
		return nil, fmt.Errorf("X509KeyPair error : %w", err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(opts.TLSConfig.RootCa) {
		return nil, fmt.Errorf("AppendCertsFromPEM error")
	}

	cli := &client{}
	clientOptions := pahoMqtt.NewClientOptions().AddBroker(opts.Server).SetTLSConfig(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}).SetConnectionLostHandler(func(client pahoMqtt.Client, err error) {
		cli.connectionLostHandler()
		// log.Println(err)
	})
	if opts.ClientID != "" {
		clientOptions = clientOptions.SetCleanSession(false).SetClientID(opts.ClientID)
	}
	pahoMqttCli := pahoMqtt.NewClient(clientOptions)

	cli.pahoMqttCli = pahoMqttCli
	token := pahoMqttCli.Connect()
	ok := token.WaitTimeout(time.Second * 5)
	if ok && token.Error() != nil {
		return nil, fmt.Errorf("connect error : %w", token.Error())
	} else if !ok {
		return nil, fmt.Errorf("connect timeout")
	}

	return cli, nil
}

func (cli *client) IsConnected() bool {
	return cli.pahoMqttCli.IsConnectionOpen()
}

func (cli *client) Publish(topic string, payload interface{}) error {
	token := cli.pahoMqttCli.Publish(topic, 1, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("publish error : %w", token.Error())
	}
	return nil
}

func (cli *client) Subscribe(ctx context.Context, topic string) (rxgo.Observable, error) {
	isClosed := false
	_, ok := cli.subscriberMap.Load(topic)
	if ok {
		return nil, fmt.Errorf("topic %s is already subscribed", topic)
	}
	ch := make(chan rxgo.Item)
	token := cli.pahoMqttCli.Subscribe(topic, 1, func(client pahoMqtt.Client, msg pahoMqtt.Message) {
		// fmt.Println(msg.Topic(), string(msg.Payload()))
		cli.mu.Lock()
		if !isClosed {
			ch <- rxgo.Of(msg.Payload())
		}
		cli.mu.Unlock()
	})
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("subscribe error : %w", token.Error())
	}
	cli.subscriberMap.Store(topic, &ch)
	go func() {
		<-ctx.Done()
		cli.mu.Lock()
		cli.subscriberMap.Delete(topic)
		close(ch)
		isClosed = true
		cli.pahoMqttCli.Unsubscribe(topic)
		cli.mu.Unlock()
	}()
	return rxgo.FromChannel(ch), nil
}

type Message struct {
	Topic   string
	Payload []byte
}

func (cli *client) SubscribeReturnMessage(ctx context.Context, topic string) (rxgo.Observable, error) {
	isClosed := false
	_, ok := cli.subscriberMap.Load(topic)
	if ok {
		return nil, fmt.Errorf("topic %s is already subscribed", topic)
	}
	ch := make(chan rxgo.Item)
	token := cli.pahoMqttCli.Subscribe(topic, 1, func(client pahoMqtt.Client, msg pahoMqtt.Message) {
		// fmt.Println(msg.Topic(), string(msg.Payload()))
		cli.mu.Lock()
		if !isClosed {
			ch <- rxgo.Of(Message{Topic: msg.Topic(), Payload: msg.Payload()})
		}
		cli.mu.Unlock()
	})
	if token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("subscribe error : %w", token.Error())
	}
	cli.subscriberMap.Store(topic, &ch)

	go func() {
		<-ctx.Done()
		cli.mu.Lock()
		cli.subscriberMap.Delete(topic)
		close(ch)
		isClosed = true
		cli.pahoMqttCli.Unsubscribe(topic)
		cli.mu.Unlock()
	}()
	return rxgo.FromChannel(ch), nil
}

func (cli *client) connectionLostHandler() {
	cli.mu.Lock()
	cli.subscriberMap.Range(func(key, value interface{}) bool {
		ch, ok := value.(*chan rxgo.Item)
		if !ok {
			return true
		}
		*ch <- rxgo.Error(fmt.Errorf("connection lost"))
		return true
	})
	cli.mu.Unlock()
	cli.subscriberMap = sync.Map{}
}

func (cli *client) Close() {
	cli.pahoMqttCli.Disconnect(250)
	cli.subscriberMap = sync.Map{}
}
