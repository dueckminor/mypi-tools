package automation

////////////////////////////////////////////////////////////////////// enum Unit

type Unit int64

const (
	Unit_Number Unit = iota
	Unit_Wh
	Unit_kWh
	Unit_W
	Unit_kW
	Unit_V
	Unit_A
	Unit_Percent
)

func (u Unit) String() string {
	switch u {
	case Unit_Number:
		return "Number"
	case Unit_Wh:
		return "Wh"
	case Unit_kWh:
		return "kWh"
	case Unit_W:
		return "W"
	case Unit_kW:
		return "kW"
	case Unit_V:
		return "V"
	case Unit_A:
		return "A"
	case Unit_Percent:
		return "%"
	default:
		panic("unsupported 'Unit'")
	}
}

//////////////////////////////////////////////////////////////// enum StateClass

type StateClass int64

const (
	StateClass_Measurement StateClass = iota
	StateClass_Total
)

func (sc StateClass) String() string {
	switch sc {
	case StateClass_Measurement:
		return "measurement"
	case StateClass_Total:
		return "total"
	default:
		panic("unsupported 'StateClass'")
	}
}

/////////////////////////////////////////////////////////////// enum DeviceClass

type DeviceClass int64

const (
	DeviceClass_Energy DeviceClass = iota
)

func (dc DeviceClass) String() string {
	switch dc {
	case DeviceClass_Energy:
		return "energy"
	default:
		panic("unsupported 'DeviceClass'")
	}
}

////////////////////////////////////////////////////////////////////// enum Icon

type Icon int64

const (
	Icon_Wh Icon = iota
	Icon_W
	Icon_V
	Icon_A
	Icon_Battery
)

func (i Icon) String() string {
	switch i {
	case Icon_Wh:
		return "mdi:lightning-bolt"
	case Icon_W, Icon_V, Icon_A:
		return "mdi:flash"
	case Icon_Battery:
		return "mdi:battery-10"
	default:
		panic("unsupported 'Icon'")
	}
}

////////////////////////////////////////////////////////////////////////////////

type ObjectTemplate struct {
	name                 string
	separateAvailability bool
}

type SensorTemplate struct {
	ObjectTemplate
	icon        Icon
	unit        Unit
	deviceClass DeviceClass
	stateClass  StateClass
	precision   uint
}

func MakeSensorTemplate(name string) *SensorTemplate {
	return &SensorTemplate{ObjectTemplate: ObjectTemplate{name: name}}
}

func (t *SensorTemplate) SetIcon(icon Icon) *SensorTemplate {
	t.icon = icon
	return t
}

func (t *SensorTemplate) SetUnit(unit Unit) *SensorTemplate {
	t.unit = unit
	return t
}

func (t *SensorTemplate) SetPrecision(precision uint) *SensorTemplate {
	t.precision = precision
	return t
}

func (t *SensorTemplate) SetDeviceClass(deviceClass DeviceClass) *SensorTemplate {
	t.deviceClass = deviceClass
	return t
}

func (t *SensorTemplate) SetStateClass(stateClass StateClass) *SensorTemplate {
	t.stateClass = stateClass
	return t
}
