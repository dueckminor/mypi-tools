package alphaess

import (
	"fmt"
	"log"
	"time"

	"github.com/dueckminor/mypi-tools/go/homeassistant"
	"github.com/dueckminor/mypi-tools/go/influxdb"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/simonvetter/modbus"
)

type SensorInfo struct {
	Addr  uint16
	Type  string
	Name  string
	Unit  string
	Scale float64
	Icon  string
}

func SensorWh(addr uint16, name string) SensorInfo {
	return SensorInfo{Addr: addr, Type: "uint16", Name: name, Unit: "Wh", Scale: 10}
}

var sensorInfos = []SensorInfo{
	SensorWh(0x0010, "to_grid"),
	SensorWh(0x0012, "from_grid"),
	SensorWh(0x0090, "total_to_grid"),
	SensorWh(0x0092, "total_from_grid"),
	SensorWh(0x0120, "battery_charge"),
	SensorWh(0x0122, "battery_discharge"),
	SensorWh(0x0124, "battery_charge_from_grid"),
	SensorWh(0x0720, "inverter_total_pv_energy"),
	SensorWh(0x08D2, "solar_production"),
}

func RegisterSensors(ha homeassistant.HomeAssistantMqtt) {
	for _, sensorInfo := range sensorInfos {
		config := homeassistant.Config{
			DeviceClass:       "energy",
			StateClass:        "total",
			Name:              sensorInfo.Name,
			StateTopic:        fmt.Sprintf("alphaess/sensor/%s/state", sensorInfo.Name),
			UnitOfMeasurement: sensorInfo.Unit,
			Icon:              "mdi:lightning-bolt",
			UniqueId:          fmt.Sprintf("alphaess.%s", sensorInfo.Name),
			Device: &homeassistant.DeviceConfig{
				Identifiers:  []string{"alphaess_sensor"},
				Name:         "Alpha ESS",
				Model:        "Alpha ESS",
				Manufacturer: "Alpha ESS",
			},
		}
		err := ha.AddSensorConfig("alphaess", sensorInfo.Name, config)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Run(uri string, mqttClient mqtt.Client, influx influxdb.Client) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://alpha-ess:502",
		Timeout: 1 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}

	err = client.SetUnitId(0x55)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for _, sensorInfo := range sensorInfos {
			value, err := client.ReadUint32(sensorInfo.Addr, modbus.HOLDING_REGISTER)
			if err != nil {
				fmt.Println(err)
			}
			scaledValue := float64(value) * sensorInfo.Scale
			strValue := fmt.Sprintf("%f", scaledValue)
			if strValue == "0.000000" {
				strValue = "0"
			}
			mqttClient.Publish("alphaess/sensor/"+sensorInfo.Name+"/state", 0, false, strValue)

			if influx != nil {
				influx.SendMetric(sensorInfo.Name, scaledValue)
			}
		}
		time.Sleep(time.Minute)
	}

}
