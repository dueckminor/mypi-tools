package main

import (
	"os"
	"strings"

	"github.com/dueckminor/mypi-tools/go/automation"
	"github.com/dueckminor/mypi-tools/go/automation/alphaess"
	"github.com/dueckminor/mypi-tools/go/protocols/influxdb"
	"github.com/dueckminor/mypi-tools/go/protocols/mqtt"
	"github.com/dueckminor/mypi-tools/go/util"
	"gopkg.in/yaml.v3"
)

type HomematicClientConfig struct {
	URI      string `yaml:"uri"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type HomeassistantConfig struct {
	Enabled bool `yaml:"enabled"`
}

type Config struct {
	MQTT         mqtt.MQTTClientConfig `yaml:"mqtt"`
	Homematic    HomematicClientConfig `yaml:"homematic"`
	Homeassisant HomeassistantConfig   `yaml:"homeassistant"`
	AlphaEss     alphaess.Config       `yaml:"alphaess"`
	InfluxDB     influxdb.Config       `yaml:"influxdb"`
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

	registry := automation.GetRegistry()

	broker := mqtt.NewBroker(cfg.MQTT.URI)
	registry.EnableMqtt(broker)

	if cfg.Homeassisant.Enabled {
		registry.EnableHomeAssistant()
	}

	var influx influxdb.Client
	if cfg.InfluxDB.Uri != "" {
		influx = influxdb.NewClient(cfg.InfluxDB)
		registry.EnableInfluxDB(influx)
	}

	alphaEssURI := cfg.AlphaEss.URI
	if alphaEssURI != "" {
		alphaess.Run(alphaEssURI)
	}

	// uri := cfg.CCU.URI
	// if uri != "" {

	// 	if len(cfg.CCU.Username) > 0 {
	// 		parsedURI, err := url.Parse(uri)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		parsedURI.User = url.UserPassword(cfg.CCU.Username, cfg.CCU.Password)
	// 		uri = parsedURI.String()
	// 	}

	// 	ccuc, err := ccu.NewCcuClient(uri)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	ccuc.SetCallback(func(dev ccu.Device, valueKey string, value interface{}) {
	// 		topic := "hm/" + dev.Address() + "/" + valueKey

	// 		payload, _ := json.Marshal(value)

	// 		mqttClient.Publish(topic, 2, false, string(payload))
	// 		fmt.Println("<-", topic, string(payload))
	// 	})

	// 	devices, _ := ccuc.GetDevices()

	// 	for _, device := range devices {
	// 		if _, err := device.GetValues(); err != nil {
	// 			continue
	// 		}
	// 		topic := "hm/" + device.Address() + "/@TYPE"
	// 		payload := device.Type()
	// 		mqttClient.Publish(topic, 2, true, payload)
	// 		fmt.Println("<-", topic, payload)
	// 	}

	// 	mqttClient.Subscribe("hm/#", 2, func(client mqtt.Client, msg mqtt.Message) {
	// 		topic := msg.Topic()
	// 		topicParts := strings.Split(topic, "/")
	// 		addr := topicParts[1]
	// 		valueName := topicParts[2]

	// 		device, err := ccuc.GetDevice(addr)

	// 		if valueName == "_TYPE_" || (nil == device && msg.Retained()) {
	// 			mqttClient.Publish(topic, 2, true, "")
	// 			return
	// 		}
	// 		if len(valueName) == 0 || valueName[0] == '@' {
	// 			return
	// 		}

	// 		if device != nil && err == nil {
	// 			var value interface{}
	// 			err = json.Unmarshal(msg.Payload(), &value)
	// 			if err != nil {
	// 				return
	// 			}
	// 			changed, _ := device.SetValueIfChanged(valueName, value)
	// 			if changed {
	// 				fmt.Println("->", topic, value)
	// 			}
	// 		}
	// 	})

	// 	err = ccuc.StartCallbackHandler()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	done := make(chan bool)

	<-done
}
