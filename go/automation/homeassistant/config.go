package homeassistant

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
