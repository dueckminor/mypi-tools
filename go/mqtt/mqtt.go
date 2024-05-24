package mqtt

import (
	"crypto/tls"
	"io"

	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClientConfig struct {
	URI      string `yaml:"uri"`
	ClientID string `yaml:"client_id"`
}

type Broker interface {
	Dial(clientId string, statusTopic string) (Conn, error)
}

type Conn interface {
	io.Closer
	Publish(topic string, payload string)
	PublishRetain(topic string, payload string)
	Subscribe(topic string, cb func(topic string, payload string))
}

type broker struct {
	tlsConfig *tls.Config
	uri       string
}

func (b *broker) Dial(clientId string, statusTopic string) (Conn, error) {
	c := &conn{}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(b.uri)
	opts.SetClientID(clientId).SetTLSConfig(b.tlsConfig)
	if statusTopic != "" {
		opts.SetWill(statusTopic, "offline", 0, true)
	}
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if statusTopic != "" {
		c.PublishRetain(statusTopic, "online")
	}
	return c, nil
}

type conn struct {
	client mqtt.Client
}

func (c *conn) Close() error {
	c.client.Disconnect(500)
	return nil
}

func (c *conn) Publish(topic string, payload string) {
	c.client.Publish(topic, 0, false, payload)
}

func (c *conn) PublishRetain(topic string, payload string) {
	c.client.Publish(topic, 0, true, payload)
}

func (c *conn) Subscribe(topic string, cb func(topic string, payload string)) {
	c.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		cb(msg.Topic(), string(msg.Payload()))
	})
}

func NewBroker(uri string) Broker {
	b := &broker{}
	b.tlsConfig = tlsconfig.NewTLSConfig()
	b.uri = uri
	return b
}
