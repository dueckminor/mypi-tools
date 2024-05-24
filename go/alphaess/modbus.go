package alphaess

import (
	"fmt"
	"time"

	"github.com/dueckminor/mypi-tools/go/homeassistant"
	"github.com/dueckminor/mypi-tools/go/influxdb"
	"github.com/dueckminor/mypi-tools/go/mqtt"
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
			AvailabilityTopic: "alphaess/status",
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

type scanner struct {
	client *modbus.ModbusClient
	broker mqtt.Broker
	conn   mqtt.Conn
	influx influxdb.Client
}

func (s *scanner) mqttConnect() (err error) {
	if s.conn != nil {
		return nil
	}
	s.conn, err = s.broker.Dial("alphaess", "alphaess/status")
	return err
}

func (s *scanner) modbusConnect() (err error) {
	err = s.client.Open()
	if err != nil {
		return err
	}
	err = s.client.SetUnitId(0x55)
	if err != nil {
		return err
	}
	return nil
}

func Run(uri string, broker mqtt.Broker, influx influxdb.Client) (err error) {
	s := &scanner{
		broker: broker,
		influx: influx,
	}

	s.client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     uri,
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return err
	}

	err = s.mqttConnect()
	if err != nil {
		return err
	}

	err = s.modbusConnect()
	if err != nil {
		return err
	}

	go func() {
		for {
			for _, sensorInfo := range sensorInfos {
				value, err := s.client.ReadUint32(sensorInfo.Addr, modbus.HOLDING_REGISTER)
				if err != nil {
					fmt.Println(err)
				}
				var strValue string
				scaledValue := float64(value) * sensorInfo.Scale
				if sensorInfo.Scale >= 1.0 {
					strValue = fmt.Sprintf("%d", int64(scaledValue))
				} else {
					strValue = fmt.Sprintf("%f", scaledValue)
					if strValue == "0.000000" {
						strValue = "0"
					}
				}
				s.conn.Publish("alphaess/sensor/"+sensorInfo.Name+"/state", strValue)

				if s.influx != nil {
					s.influx.SendMetric(sensorInfo.Name, scaledValue)
				}
			}
			time.Sleep(time.Minute)
		}
	}()

	return nil
}
