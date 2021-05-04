package ccu

import "net/url"

type CcuMultiClientImpl struct {
	clients []*CcuClientImpl
}

func NewCcuClient(uri string) (ccuc CcuClient, err error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	// port 2001: HM
	// port 2010: HM-IP

	port := parsedURI.Port()
	if len(port) > 0 {
		return newCcuClient(uri)
	}

	host := parsedURI.Host
	parsedURI.Host = host + ":2001"

	client, err := newCcuClient(parsedURI.String())
	if err != nil {
		return nil, err
	}

	ccumc := &CcuMultiClientImpl{}
	ccumc.clients = append(ccumc.clients, client)

	parsedURI.Host = host + ":2010"
	client, err = newCcuClient(parsedURI.String())
	if err != nil {
		return nil, err
	}

	ccumc.clients = append(ccumc.clients, client)
	return ccumc, err
}

func (ccumc *CcuMultiClientImpl) GetVersion() (version string, err error) {
	return ccumc.clients[0].GetVersion()
}
func (ccumc *CcuMultiClientImpl) SetCallback(cb CcuCallback) {
	for _, ccuc := range ccumc.clients {
		ccuc.SetCallback(cb)
	}
}
func (ccumc *CcuMultiClientImpl) StartCallbackHandler() error {
	for _, ccuc := range ccumc.clients {
		err := ccuc.StartCallbackHandler()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ccumc *CcuMultiClientImpl) GetOwnIP() string {
	return ccumc.clients[0].GetOwnIP()
}

func (ccumc *CcuMultiClientImpl) GetDevices() (devices []Device, err error) {
	for _, ccuc := range ccumc.clients {
		devs, err := ccuc.GetDevices()
		if err != nil {
			return nil, err
		}
		devices = append(devices, devs...)
	}
	return devices, nil
}

func (ccumc *CcuMultiClientImpl) GetDevice(addr string) (device Device, err error) {
	for _, ccuc := range ccumc.clients {
		dev, _ := ccuc.GetDevice(addr)
		if dev != nil {
			return dev, nil
		}
	}
	return nil, nil
}
