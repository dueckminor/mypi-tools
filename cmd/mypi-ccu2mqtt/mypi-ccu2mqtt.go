package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/dueckminor/mypi-tools/go/automation/homematic"
	"github.com/dueckminor/mypi-tools/go/protocols/mqtt"
	"github.com/dueckminor/mypi-tools/go/util"
	"gopkg.in/yaml.v3"
)

type CCUClientConfig struct {
	URI      string `yaml:"uri"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	MQTT mqtt.MQTTClientConfig `yaml:"mqtt"`
	CCU  CCUClientConfig       `yaml:"ccu"`
}

func main() {
	var cfg Config

	if (len(os.Args) == 2) && !strings.HasPrefix(os.Args[1], "-") && util.FileExists(os.Args[1]) {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			panic(err)
		}
	}

	broker := mqtt.NewBroker(cfg.MQTT.URI)
	mqttClient, err := broker.Dial(cfg.MQTT.ClientID, "")

	uri := cfg.CCU.URI
	if len(cfg.CCU.Username) > 0 {
		parsedURI, err := url.Parse(uri)
		if err != nil {
			panic(err)
		}
		parsedURI.User = url.UserPassword(cfg.CCU.Username, cfg.CCU.Password)
		uri = parsedURI.String()
	}

	ccuc, err := homematic.NewCcuClient(uri)
	if err != nil {
		panic(err)
	}

	ccuc.SetCallback(func(dev homematic.Device, valueKey string, value interface{}) {
		topic := "hm/" + dev.Address() + "/" + valueKey

		payload, _ := json.Marshal(value)

		mqttClient.Publish(topic, string(payload))
		fmt.Println("<-", topic, string(payload))
	})

	devices, _ := ccuc.GetDevices()

	for _, device := range devices {
		if _, err := device.GetValues(); err != nil {
			continue
		}
		topic := "hm/" + device.Address() + "/@TYPE"
		payload := device.Type()
		mqttClient.PublishRetain(topic, payload)
		fmt.Println("<-", topic, payload)
	}

	mqttClient.Subscribe("hm/#", func(topic string, payload string) {
		topicParts := strings.Split(topic, "/")
		addr := topicParts[1]
		valueName := topicParts[2]

		device, err := ccuc.GetDevice(addr)

		if valueName == "_TYPE_" || (nil == device) {
			mqttClient.PublishRetain(topic, "")
			return
		}
		if len(valueName) == 0 || valueName[0] == '@' {
			return
		}

		if device != nil && err == nil {
			var value interface{}
			err = json.Unmarshal([]byte(payload), &value)
			if err != nil {
				return
			}
			changed, _ := device.SetValueIfChanged(valueName, value)
			if changed {
				fmt.Println("->", topic, value)
			}
		}
	})

	err = ccuc.StartCallbackHandler()
	if err != nil {
		panic(err)
	}

	done := make(chan bool)

	<-done

}
