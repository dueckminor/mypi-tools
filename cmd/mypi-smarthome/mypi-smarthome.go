package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	// authURI   string
	port       = flag.Int("port", 8080, "The port")
	mypiRoot   = flag.String("mypi-root", "", "The root of the mypi filesystem")
	mqttBroker = flag.String("mqtt", "mqtt", "The MQTT Broker")
	// targetURI string

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

	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}
	config.GetConfig()
}

type DeviceInfo struct {
	Address    string `yaml:"address"`
	Name       string `yaml:"name"`
	TopicCool  string `yaml:"topic_cool"`
	TopicHeat  string `yaml:"topic_heat"`
	TempShould float64
	TempIs     float64
}

type Config struct {
	Climate []*DeviceInfo `yaml:"climate"`
}

type ClimateControl struct {
	mqttClient mqtt.Client
	devices    map[string]*DeviceInfo
}

func NewClimateControl(mqttClient mqtt.Client) *ClimateControl {
	c := &ClimateControl{
		mqttClient: mqttClient,
		devices:    make(map[string]*DeviceInfo),
	}
	return c
}

func (c *ClimateControl) AddDevices(devices ...*DeviceInfo) {
	for _, device := range devices {
		c.devices[device.Address] = device
		device.TempIs = math.NaN()
		device.TempShould = math.NaN()
	}
}

func (c *ClimateControl) HandleMQTT(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	topicParts := strings.Split(msg.Topic(), "/")
	addr := topicParts[1]
	what := topicParts[2]

	temp, _ := strconv.ParseFloat(payload, 64)

	deviceInfo := c.devices[addr]
	if nil == deviceInfo {
		return
	}

	switch what {
	case "ACTUAL_TEMPERATURE":
		deviceInfo.TempIs = temp
		fmt.Println(deviceInfo.Address, deviceInfo.Name, "is", temp)
	case "SET_POINT_TEMPERATURE":
		deviceInfo.TempShould = temp
		fmt.Println(deviceInfo.Address, deviceInfo.Name, "should", temp)
	}
}

func (c *ClimateControl) AdjustAll() {
	for _, device := range c.devices {
		if math.IsNaN(device.TempIs) || math.IsNaN(device.TempShould) {
			continue
		}
		tempDiff := device.TempIs - device.TempShould
		fmt.Println(device.Name, "diff:", tempDiff)

		topic := ""
		if tempDiff > 0.1 {
			topic = device.TopicCool
		} else if tempDiff < -0.1 {
			tempDiff = -tempDiff
			topic = device.TopicHeat
		} else {
			continue
		}
		amount := int64(math.Round(tempDiff * 10))
		if amount > 30 {
			amount = 30
		}
		value := fmt.Sprintf("%ds", amount)
		fmt.Println(topic, value)
		c.mqttClient.Publish(topic, 2, false, value)
	}
}

func main() {

	var cfg Config
	if len(flag.Args()) == 1 {
		config.ReadYAML(&cfg, flag.Arg(0))
	}

	tlsconfig := tlsconfig.NewTLSConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://" + *mqttBroker + ":8883")
	opts.SetTLSConfig(tlsconfig)

	// Start the connection
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	c := NewClimateControl(mqttClient)
	c.AddDevices(cfg.Climate...)

	mqttClient.Subscribe("hm/+/ACTUAL_TEMPERATURE", 2, c.HandleMQTT)
	mqttClient.Subscribe("hm/+/SET_POINT_TEMPERATURE", 2, c.HandleMQTT)

	for {
		time.Sleep(time.Minute * 2)
		c.AdjustAll()
	}
}
