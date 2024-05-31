package automation

////////////////////////////////////////////////////////////////////// enum Unit

type Unit int64

const (
	Unit_Number Unit = iota
	Unit_Wh
	Unit_kWh
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

////////////////////////////////////////////////////////////////////////////////

type ObjectTemplate struct {
	name                 string
	separateAvailability bool
}

type SensorTemplate struct {
	ObjectTemplate
	unit        Unit
	deviceClass DeviceClass
	stateClass  StateClass
	precision   uint
}

func MakeSensorTemplate(name string) *SensorTemplate {
	return &SensorTemplate{ObjectTemplate: ObjectTemplate{name: name}}
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
