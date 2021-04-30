package ccu

import (
	"strings"

	"github.com/dueckminor/mypi-tools/go/xmlrpc"
	"golang.org/x/net/html/charset"
)

type CcuCallback func(dev Device, valueKey string, value interface{})

type CcuClient struct {
	xmlrpcClient *xmlrpc.Client
	devices      map[string]deviceInt
	callbacks    []CcuCallback
}

func NewCcuClient(uri string) (ccuc *CcuClient, err error) {
	xmlrpc.CharsetReader = charset.NewReaderLabel
	ccuc = &CcuClient{}

	ccuc.xmlrpcClient, err = xmlrpc.NewClient(uri, nil)
	return ccuc, err
}

func (ccuc *CcuClient) ParseMethodCall(methodName string, cb xmlrpc.MethodCallParserCB) (err error) {
	//fmt.Println(methodName)
	switch methodName {
	case "event":
		var interfaceID string
		var address string
		var valueKey string
		var value interface{}
		err = cb.GetCallParam(&interfaceID)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&address)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&valueKey)
		if err != nil {
			return err
		}
		err = cb.GetCallParam(&value)
		if err != nil {
			return err
		}
		//fmt.Println(interfaceID, address, valueKey, value)
		if dev, ok := ccuc.devices[address]; ok {
			dev.putValue(valueKey, value)
		}
		cb.PutResult(nil)
		return nil
	case "system.listMethods":
		cb.IgnoreParams()
		cb.PutResult([]string{})
		return nil
	default:
		cb.IgnoreParams()
		cb.PutResult(nil)
		return nil
	}
}

func (ccuc *CcuClient) SetCallback(cb CcuCallback) {
	ccuc.callbacks = append(ccuc.callbacks, cb)
}

func (ccuc *CcuClient) Init(url string, interfaceID string) (err error) {
	return ccuc.xmlrpcClient.Call("init", []interface{}{url, interfaceID}, nil)
}

func (ccuc *CcuClient) GetVersion() (version string, err error) {
	err = ccuc.xmlrpcClient.Call("getVersion", nil, &version)
	return version, err
}
func (ccuc *CcuClient) ListMethods() (methods []string, err error) {
	err = ccuc.xmlrpcClient.Call("system.listMethods", nil, &methods)
	return methods, err
}

func (ccuc *CcuClient) ListDevices() (devices []DeviceDescription, err error) {
	err = ccuc.xmlrpcClient.Call("listDevices", nil, &devices)
	return devices, err
}

func (ccuc *CcuClient) GetDeviceDescription(address string) (device *DeviceDescription, err error) {
	err = ccuc.xmlrpcClient.Call("getDeviceDescription", []interface{}{address}, &device)
	return device, err
}

func (ccuc *CcuClient) GetParamsetDescription(address string, paramsetType string) (paramsetDescription ParamsetDescription, err error) {
	err = ccuc.xmlrpcClient.Call("getParamsetDescription", []interface{}{
		address, paramsetType,
	}, &paramsetDescription)
	return paramsetDescription, err
}

func (ccuc *CcuClient) GetMasterDescription(address string) (paramsetDescription ParamsetDescription, err error) {
	return ccuc.GetParamsetDescription(address, "MASTER")
}

func (ccuc *CcuClient) GetValueDescription(address string) (paramsetDescription ParamsetDescription, err error) {
	return ccuc.GetParamsetDescription(address, "VALUES")
}

func (ccuc *CcuClient) GetLinkDescription(address string) (paramsetDescription ParamsetDescription, err error) {
	return ccuc.GetParamsetDescription(address, "LINK")
}

func (ccuc *CcuClient) GetValue(address, valueKey string) (value interface{}, err error) {
	err = ccuc.xmlrpcClient.Call("getValue", []interface{}{address, valueKey}, &value)
	return value, err
}

func (ccuc *CcuClient) GetParamsetID(address, paramsetType string) (value string, err error) {
	err = ccuc.xmlrpcClient.Call("getParamsetId", []interface{}{address, paramsetType}, &value)
	return value, err
}

func (ccuc *CcuClient) GetParamset(address, paramsetKey string) (value map[string]interface{}, err error) {
	err = ccuc.xmlrpcClient.Call("getParamset", []interface{}{address, paramsetKey}, &value)
	return value, err
}

func (ccuc *CcuClient) SetValue(address, valueKey string, value interface{}) (err error) {
	return ccuc.xmlrpcClient.Call("setValue", []interface{}{address, valueKey, value}, nil)
}

func (ccuc *CcuClient) getDevices() (err error) {
	devs, err := ccuc.ListDevices()
	if err != nil {
		return
	}

	newDevices := make(map[string]deviceInt)

	for _, dev := range devs {
		if impl, ok := ccuc.devices[dev.Address]; ok {
			newDevices[dev.Address] = impl
			continue
		}

		var device deviceInt

		switch dev.Type {
		case `HmIPW-WTH`, `HmIP-WTH-2`:
			device = newWTH(ccuc, dev)
		default:
			device = newDevice(ccuc, dev)
		}
		newDevices[dev.Address] = device

		addrParts := strings.Split(dev.Address, ":")
		if len(addrParts) > 1 {
			devImpl := newDevices[addrParts[0]].(deviceInt)
			devImpl.setSubdevice(device)
		}

	}

	ccuc.devices = newDevices
	return nil
}

func (ccuc *CcuClient) GetDevices() (devices []Device, err error) {
	err = ccuc.getDevices()
	if err != nil {
		return nil, err
	}
	for _, dev := range ccuc.devices {
		devices = append(devices, dev)
	}
	return devices, nil
}

func (ccuc *CcuClient) GetWTHs() (wths []*WTH, err error) {
	err = ccuc.getDevices()
	if err != nil {
		return nil, err
	}
	for _, dev := range ccuc.devices {
		if wth, ok := dev.(*WTH); ok {
			wths = append(wths, wth)
		}
	}
	return wths, nil
}
