package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	"github.com/dueckminor/mypi-tools/go/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	hostname string
	autoOpen = false
	nearHome = true
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
}

type OwntracksMsg struct {
	InRegions []string `json:"inregions"`
}

func mqttOwnTracks(client mqtt.Client, msg mqtt.Message) {
	ownTracksMsg := OwntracksMsg{}
	json.Unmarshal(msg.Payload(), &ownTracksMsg)
	if util.StringsContains(ownTracksMsg.InRegions, "Work") {
		autoOpen = true
		client.Publish("tor/autoopen", 2, true, `true`)
	}
	if util.StringsContains(ownTracksMsg.InRegions, "NearHome") {
		if autoOpen {
			autoOpen = false
			fmt.Println("Open gate")
			client.Publish("tor/open", 2, true, `1`)
		}
		client.Publish("tor/autoopen", 2, true, `false`)
		nearHome = true
	} else {
		nearHome = false
	}
}

func mqttAutoOpen(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	shouldAutoOpen := (payload == "true") || (payload == "1")
	if nearHome && shouldAutoOpen {
		fmt.Println("reject autoopen request")
		client.Publish("tor/autoopen", 2, true, `false`)
		return
	}
	autoOpen = shouldAutoOpen
}

func main() {

	tlsconfig := tlsconfig.NewTLSConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://rpi:8883")
	opts.SetClientID(hostname).SetTLSConfig(tlsconfig)

	// Start the connection
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient.Subscribe("tor/autoopen", 2, mqttAutoOpen)
	mqttClient.Subscribe("owntracks/#", 2, mqttOwnTracks)

	quit := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		mqttClient.Disconnect(250)
		fmt.Println("[MQTT] Disconnected")

		quit <- struct{}{}
	}()
	<-quit
}
