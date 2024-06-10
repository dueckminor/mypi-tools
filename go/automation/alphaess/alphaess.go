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

	sensors                   []*sensor
	sensorSolarProduction     *sensor
	correctionSolarProduction float64
	sensorToGrid              *sensor
}

type sensor struct {
	automation.Sensor
	Addr   uint16
	Signed bool
	Words  int
	Scale  float64

	Last    float64
	Current float64
}

func (s *scanner) init() {
	s.registry = automation.GetRegistry()
	s.node = s.registry.CreateNode("alphaess")
	// -------------------------------------------------------------------- grid
	s.sensorToGrid =
		s.sensor10Wh(0x0010, "to_grid")
	s.sensor10Wh(0x0012, "from_grid")
	s.sensor1V(0x0014, "grid_voltage_l1")
	s.sensor1V(0x0015, "grid_voltage_l2")
	s.sensor1V(0x0016, "grid_voltage_l3")
	s.sensor100mA(0x0017, "grid_current_l1")
	s.sensor100mA(0x0018, "grid_current_l2")
	s.sensor100mA(0x0019, "grid_current_l3")
	// s.sensorHz(0x001a, "grid_freq")
	s.sensor1W(0x001b, "grid_active_power_l1")
	s.sensor1W(0x001d, "grid_active_power_l2")
	s.sensor1W(0x001f, "grid_active_power_l3")
	s.sensor1W(0x0021, "grid_active_power")
	s.sensor1W(0x0023, "grid_reactive_power_l1")
	s.sensor1W(0x0025, "grid_reactive_power_l2")
	s.sensor1W(0x0027, "grid_reactive_power_l3")
	s.sensor1W(0x0029, "grid_reactive_power")
	s.sensor1W(0x002b, "grid_apparent_power_l1")
	s.sensor1W(0x002d, "grid_apparent_power_l2")
	s.sensor1W(0x002f, "grid_apparent_power_l3")
	s.sensor1W(0x0031, "grid_apparent_power")
	// ---------------------------------------------------------------- pv meter
	s.sensor10Wh(0x0090, "total_to_grid")
	s.sensor10Wh(0x0092, "total_from_grid")
	// ----------------------------------------------------------------- battery
	s.sensor100mV(0x0100, "battery_voltage")
	s.sensor100mA(0x0101, "battery_current")
	s.sensor100Wh(0x0120, "battery_charge")
	s.sensor100Wh(0x0122, "battery_discharge")
	s.sensor100Wh(0x0124, "battery_charge_from_grid")
	s.sensor1W(0x0126, "battery_power")

	s.sensor10Wh(0x0720, "inverter_total_pv_energy")

	s.sensorSolarProduction =
		s.sensor10Wh(0x08D2, "solar_production")
	s.correctionSolarProduction = 0.0
	s.sensorPercent(0x0102, "battery_soc")

}

func (s *scanner) sensor10Wh(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_Wh).
			SetUnit(automation.Unit_Wh).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 2,
		Scale: 10,
	}

	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensor100Wh(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_Wh).
			SetUnit(automation.Unit_Wh).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 2,
		Scale: 100,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensor1W(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_W).
			SetUnit(automation.Unit_W).SetPrecision(0).
			SetStateClass(automation.StateClass_Total).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:   addr,
		Signed: true,
		Words:  2,
		Scale:  1,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensor1V(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_V).
			SetUnit(automation.Unit_V).SetPrecision(0).
			SetStateClass(automation.StateClass_Measurement).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 1,
		Scale: 1,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensor100mV(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_V).
			SetUnit(automation.Unit_V).SetPrecision(1).
			SetStateClass(automation.StateClass_Measurement).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 1,
		Scale: 0.1,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensor100mA(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_A).
			SetUnit(automation.Unit_A).SetPrecision(1).
			SetStateClass(automation.StateClass_Measurement).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:   addr,
		Signed: true,
		Words:  1,
		Scale:  0.1,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
}

func (s *scanner) sensorPercent(addr uint16, name string) *sensor {
	sensor := &sensor{
		Sensor: s.node.CreateSensor(automation.MakeSensorTemplate(name).
			SetIcon(automation.Icon_Battery).
			SetUnit(automation.Unit_Percent).SetPrecision(1).
			SetStateClass(automation.StateClass_Measurement).
			SetDeviceClass(automation.DeviceClass_Energy)),
		Addr:  addr,
		Words: 1,
		Scale: 0.1,
	}
	s.sensors = append(s.sensors, sensor)
	return sensor
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
			var value int64
			if sensor.Words == 1 {
				var value16 uint16
				value16, err = s.client.ReadRegister(sensor.Addr, modbus.HOLDING_REGISTER)
				if sensor.Signed {
					value = int64(int16(value16))
				} else {
					value = int64(value16)
				}
			} else if sensor.Words == 2 {
				var value32 uint32
				value32, err = s.client.ReadUint32(sensor.Addr, modbus.HOLDING_REGISTER)
				if sensor.Signed {
					value = int64(int32(value32))
				} else {
					value = int64(value32)
				}
			}
			if err != nil {
				fmt.Println(err)
			}
			sensor.Last = sensor.Current
			sensor.Current = float64(value) * sensor.Scale
		}

		if s.sensorSolarProduction.Current == s.sensorSolarProduction.Last {
			if s.sensorToGrid.Current > s.sensorToGrid.Last {
				fmt.Println("increasing solar production correction")
				s.correctionSolarProduction +=
					(s.sensorToGrid.Current - s.sensorToGrid.Last)
				fmt.Println("new solar production correction:", s.correctionSolarProduction)
			}
		} else if s.correctionSolarProduction > 0 {
			if s.sensorSolarProduction.Current >= s.sensorSolarProduction.Last {
				fmt.Println("reseting solar production correction")
				s.correctionSolarProduction = 0
			} else {
				fmt.Println("reducing solar production correction")
				s.correctionSolarProduction = s.sensorSolarProduction.Last - s.sensorSolarProduction.Current
				fmt.Println("new solar production correction:", s.correctionSolarProduction)
			}
		}

		s.sensorSolarProduction.Current += s.correctionSolarProduction

		for _, sensor := range s.sensors {
			sensor.SetState(sensor.Current)
		}

		time.Sleep(time.Minute)
	}

}
