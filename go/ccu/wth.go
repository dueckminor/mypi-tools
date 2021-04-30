package ccu

import (
	"math"
)

type WTH struct {
	deviceImpl
}

func newWTH(ccuc *CcuClient, deviceDesc DeviceDescription) (wth *WTH) {
	wth = new(WTH)
	wth.ccuc = ccuc
	wth.deviceDesc = deviceDesc
	wth.initMaps()
	return wth
}

func (wth *WTH) Refresh() error {
	subdevice, err := wth.GetSubDevice("HEATING_CLIMATECONTROL_TRANSCEIVER")
	if err != nil {
		return err
	}
	return subdevice.SetValue("WINDOW_STATE", 2)
}

func (wth *WTH) GetTemp() (float64, error) {
	subdevice, err := wth.GetSubDevice("HEATING_CLIMATECONTROL_TRANSCEIVER")
	if err != nil {
		return math.NaN(), err
	}
	temp, err := subdevice.GetValue("ACTUAL_TEMPERATURE")
	if err != nil {
		return math.NaN(), err
	}
	return makeFloat64(temp)
}

func (wth *WTH) GetSetpointTemperature() (float64, error) {
	subdevice, err := wth.GetSubDevice("HEATING_CLIMATECONTROL_TRANSCEIVER")
	if err != nil {
		return math.NaN(), err
	}
	temp, err := subdevice.GetValue("SET_POINT_TEMPERATURE")
	if err != nil {
		return math.NaN(), err
	}
	return makeFloat64(temp)
}
