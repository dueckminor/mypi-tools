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

type scanner struct {
	client *modbus.ModbusClient

	registry automation.Registry
	node     automation.Node

	sensors []sensor
}

type sensor struct {
	automation.Sensor
	Addr  uint16
	Words int
	Scale float64
}

func (s *scanner) init() {
	s.registry = automation.GetRegistry()
	s.node = s.registry.CreateNode("alphaess")
	s.sensorWh(0x0010, "to_grid")
	s.sensorWh(0x0012, "from_grid")
	s.sensorWh(0x0090, "total_to_grid")
	s.sensorWh(0x0092, "total_from_grid")
	s.sensorWh(0x0120, "battery_charge")
	s.sensorWh(0x0122, "battery_discharge")
	s.sensorWh(0x0124, "battery_charge_from_grid")
	s.sensorWh(0x0720, "inverter_total_pv_energy")
	s.sensorWh(0x08D2, "solar_production")
	s.sensorPercent(0x0102, "battery_soc")
	s.sensorW(0x001b, "active_power_l1")
	s.sensorW(0x001d, "active_power_l2")
	s.sensorW(0x001f, "active_power_l3")
	s.sensorW(0x0021, "active_power")
	s.sensorW(0x0023, "reactive_power_l1")
	s.sensorW(0x0025, "reactive_power_l2")
	s.sensorW(0x0027, "reactive_power_l3")
	s.sensorW(0x0029, "reactive_power")
}

func (s *scanner) sensorWh(addr uint16, name string) {
	s.sensors = append(s.sensors, sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_Wh).
			SetUnit(automation.Unit_Wh).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 2,
		Scale: 10,
	})
}

func (s *scanner) sensorW(addr uint16, name string) {
	s.sensors = append(s.sensors, sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_W).
			SetUnit(automation.Unit_W).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 2,
		Scale: 10,
	})
}

func (s *scanner) sensorPercent(addr uint16, name string) {
	s.sensors = append(s.sensors, sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_Battery).
			SetUnit(automation.Unit_Percent).SetPrecision(1).
			SetStateClass(automation.StateClass_Measurement).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 1,
		Scale: 0.1,
	})
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
			var value uint32
			if sensor.Words == 1 {
				var value16 uint16
				value16, err = s.client.ReadRegister(sensor.Addr, modbus.HOLDING_REGISTER)
				value = uint32(value16)
			} else if sensor.Words == 2 {
				value, err = s.client.ReadUint32(sensor.Addr, modbus.HOLDING_REGISTER)
			}
			if err != nil {
				fmt.Println(err)
			}
			scaledValue := float64(value) * sensor.Scale
			sensor.SetState(scaledValue)
		}
		time.Sleep(time.Minute)
	}
}
