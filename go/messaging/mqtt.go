package messaging

import (
	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateClient() {
	tlsconfig := tlsconfig.NewTLSConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://rpi:8883")
	opts.SetClientID(hostname).SetTLSConfig(tlsconfig)

	mqttClient := mqtt.NewClient(opts)
}
