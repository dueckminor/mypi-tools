package alphaess

import (
	"fmt"
	"time"

	"github.com/dueckminor/mypi-tools/go/automation"
	"github.com/simonvetter/modbus"
)

type Config struct {
	URI string `yaml:"uri"`
}

type SensorInfo struct {
	Addr     uint16
	Scale    float64
	Template automation.SensorTemplate
}

func SensorWh(addr uint16, name string) SensorInfo {
	return SensorInfo{
		Addr:  addr,
		Scale: 10,
		Template: *automation.MakeSensorTemplate(name).
			SetUnit(automation.Unit_Wh).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy),
	}
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

type scanner struct {
	client *modbus.ModbusClient

	registry automation.Registry
	node     automation.Node

	sensors []sensor
}

type sensor struct {
	SensorInfo
	sensor automation.Sensor
}

func (s *scanner) init() {
	s.registry = automation.GetRegistry()
	s.node = s.registry.CreateNode("alphaess")

	for _, sensorInfo := range sensorInfos {
		s.sensors = append(s.sensors, sensor{
			SensorInfo: sensorInfo,
			sensor:     s.node.CreateSensor(&sensorInfo.Template),
		})
	}
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

func Run(uri string) (err error) {
	s := &scanner{}

	s.client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     uri,
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return err
	}

	s.init()

	go func() {
		for {
			s.handleModbus()
			time.Sleep(time.Minute)
		}
	}()

	return nil
}

func (s *scanner) handleModbus() {
	defer func() {
		s.client.Close()
		s.node.Disconnect()
		if err := recover(); err != nil {
			fmt.Println("crash in handleModbus")
		}
	}()

	err := s.modbusConnect()
	if err != nil {
		fmt.Println("modbusConnect failed: ", err)
		return
	}

	s.node.Connect()

	for {
		for _, sensor := range s.sensors {
			value, err := s.client.ReadUint32(sensor.Addr, modbus.HOLDING_REGISTER)
			if err != nil {
				fmt.Println(err)
			}
			scaledValue := float64(value) * sensor.Scale
			sensor.sensor.SetState(scaledValue)
		}
		time.Sleep(time.Minute)
	}
}
