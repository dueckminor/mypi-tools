package automation

import (
	"encoding/json"
	"fmt"

	"github.com/dueckminor/mypi-tools/go/automation/homeassistant"
	"github.com/dueckminor/mypi-tools/go/protocols/influxdb"
	"github.com/dueckminor/mypi-tools/go/protocols/mqtt"
)

// I use a naming scheme inspired by Home Assistant. Nodes and objects.
// A node is a container for multiple objects.
//
// For example:
//
// The 'alphaess' node is responsible for communication with the Alpha ESS
// modbus device. It has multiple sensor objects
//
// The 'homematic' node talks with a homematic-ccu using the XMLRPC protocol
//
// If 'homeassistant' is enabled, the mqtt discovery topic will be send.
// https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery
//
// Otherwise only the measurement itself and the up-status will be sent to
// mqtt.
//
// If 'influxdb' is enabled, the state of some objects will be sent to
// the influxdb. In contrast to 'homeassistant', the state will be sent
// even if it hasn't changed. This ensures that there are no gaps in the
// graphs (but there will be gaps in case of an outage)

type Registry interface {
	EnableMqtt(broker mqtt.Broker)
	EnableInfluxDB(influx influxdb.Client)
	EnableHomeAssistant()
	CreateNode(name string) Node
}

var theRegistry *registry

func GetRegistry() Registry {
	if theRegistry == nil {
		theRegistry = &registry{
			nodes: make(map[string]Node),
		}
	}
	return theRegistry
}

type Node interface {
	Connect() error
	Disconnect() error
	CreateSensor(template *SensorTemplate) Sensor
	CreateClimate(template *ObjectTemplate) Climate
}

type Object interface {
}

type Sensor interface {
	Object
	SetState(state any)
	Unit() Unit
}

type Climate interface {
}

////////////////////////////////////////////////////////////////////////////////

type registry struct {
	nodes                 map[string]Node
	broker                mqtt.Broker
	influx                influxdb.Client
	homeAssistant         bool
	homeAssistantMqttConn mqtt.Conn
}

func (r *registry) CreateNode(name string) Node {
	if node, ok := r.nodes[name]; ok {
		return node
	}
	node := &node{
		name:     name,
		registry: r,
		objects:  make(map[string]Object),
	}
	r.nodes[name] = node
	return node
}

func (r *registry) EnableMqtt(broker mqtt.Broker) {
	r.broker = broker
}
func (r *registry) EnableInfluxDB(influx influxdb.Client) {
	r.influx = influx
}
func (r *registry) EnableHomeAssistant() {
	r.homeAssistant = true
}

func (r *registry) publishHomeAssistantConfig(node string, objectType string, config homeassistant.Config) (err error) {
	if !r.homeAssistant {
		return nil
	}

	topic := fmt.Sprintf("homeassistant/%s/%s/%s/config", objectType, node, config.Name)

	if r.homeAssistantMqttConn == nil {
		r.homeAssistantMqttConn, err = r.broker.Dial("mypi-mqtt-bridge", "")
		if err != nil {
			return err
		}
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	r.homeAssistantMqttConn.PublishRetain(topic, string(configBytes))
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type node struct {
	registry *registry
	mqtt     mqtt.Conn
	name     string
	objects  map[string]Object
}

func (node *node) Connect() (err error) {
	if node.mqtt != nil {
		return nil
	}
	node.mqtt, err = node.registry.broker.Dial(node.name, fmt.Sprintf("%s/status", node.name))
	return err
}

func (node *node) Disconnect() (err error) {
	if node.mqtt == nil {
		return nil
	}
	node.mqtt.Close()
	node.mqtt = nil
	return nil
}

func (node *node) Publish(topic string, payload string) {
	err := node.Connect()
	if err != nil {
		return
	}
	node.mqtt.Publish(topic, payload)
}

func (node *node) CreateSensor(template *SensorTemplate) Sensor {
	if object, ok := node.objects[template.name]; ok {
		if sensor, ok := object.(Sensor); ok {
			return sensor
		}
		return nil
	}

	sensor := &sensor{}
	sensor.node = node
	sensor.influx = node.registry.influx
	sensor.template = *template

	sensor.object.stateTopic = fmt.Sprintf("%s/sensor/%s/state", node.name, template.name)
	if template.separateAvailability {
		sensor.object.availabilityTopic = fmt.Sprintf("%s/sensor/%s/status", node.name, template.name)
	} else {
		sensor.object.availabilityTopic = fmt.Sprintf("%s/status", node.name)
	}

	node.objects[template.name] = sensor

	if node.registry.homeAssistant {
		config := homeassistant.Config{
			Name:              template.name,
			DeviceClass:       template.deviceClass.String(),
			StateClass:        template.stateClass.String(),
			StateTopic:        sensor.stateTopic,
			AvailabilityTopic: sensor.availabilityTopic,
			UnitOfMeasurement: sensor.template.unit.String(),
			UniqueId:          fmt.Sprintf("%s.%s", node.name, template.name),
			Icon:              "mdi:lightning-bolt",
			Device: &homeassistant.DeviceConfig{
				Identifiers:  []string{fmt.Sprintf("%s_sensor", node.name)},
				Name:         "Alpha ESS",
				Model:        "Alpha ESS",
				Manufacturer: "Alpha ESS",
			},
		}
		node.registry.publishHomeAssistantConfig(node.name, "sensor", config)

	}

	return sensor
}

func (node *node) CreateClimate(template *ObjectTemplate) Climate {
	if object, ok := node.objects[template.name]; ok {
		return object
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type object struct {
	node              *node
	influx            influxdb.Client
	stateTopic        string
	availabilityTopic string
}

////////////////////////////////////////////////////////////////////////////////

type sensor struct {
	object
	state    any
	template SensorTemplate
}

func (sensor *sensor) SetState(state any) {
	sensor.state = state
	sensor.publish()
}

func (sensor *sensor) Unit() Unit {
	return sensor.template.unit
}

func (sensor *sensor) publish() {
	value := ""
	switch v := sensor.state.(type) {
	case string:
		value = v
	case int, int64, int32, int16, uint, uint16, uint32, uint64:
		value = fmt.Sprintf("%d", v)
	case float32:
		value = sensor.float2string(float64(v))
		if sensor.influx != nil {
			sensor.influx.SendMetric(sensor.template.name, float64(v))
		}
	case float64:
		value = sensor.float2string(float64(v))
		if sensor.influx != nil {
			sensor.influx.SendMetric(sensor.template.name, v)
		}
	}

	sensor.node.Publish(sensor.stateTopic, value)
}

func (sensor *sensor) float2string(value float64) string {
	switch sensor.template.precision {
	case 0:
		return fmt.Sprintf("%d", int64(value))
	case 1:
		return fmt.Sprintf("%0.1f", value)
	case 2:
		return fmt.Sprintf("%0.2f", value)
	case 3:
		return fmt.Sprintf("%0.3f", value)
	default:
		return fmt.Sprintf("%f", value)
	}
}

////////////////////////////////////////////////////////////////////////////////

type climate struct {
	object
}
