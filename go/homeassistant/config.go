package homeassistant

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceConfig struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Model        string   `json:"model,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
}

type Config struct {
	DeviceClass       string        `json:"device_class"`
	StateClass        string        `json:"state_class"`
	Name              string        `json:"name"`
	StateTopic        string        `json:"state_topic,omitempty"`
	UnitOfMeasurement string        `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string        `json:"value_template,omitempty"`
	UniqueId          string        `json:"unique_id,omitempty"`
	AvailabilityTopic string        `json:"availability_topic,omitempty"`
	Icon              string        `json:"icon,omitempty"`
	Device            *DeviceConfig `json:"device,omitempty"`
}

type HomeAssistantMqtt interface {
	AddSensorConfig(node string, sensor string, config Config) error
}

type homeAssistantMqtt struct {
	mqttClient mqtt.Client
}

func NewHomeAssistantMqtt(mqttClient mqtt.Client) HomeAssistantMqtt {
	return &homeAssistantMqtt{
		mqttClient: mqttClient,
	}
}

func (h *homeAssistantMqtt) AddSensorConfig(node string, sensor string, config Config) (err error) {
	topic := fmt.Sprintf("homeassistant/sensor/%s/%s/config", node, sensor)

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	h.mqttClient.Publish(topic, 0, true, configBytes)
	return nil
}
