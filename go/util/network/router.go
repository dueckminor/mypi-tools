package network

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

type cbIPAddr func() (*net.IPAddr, error)

var (
	routerExternalIP    *net.IPAddr
	routerInternalIP    *net.IPAddr
	cbRouterInternalIPs []cbIPAddr
)

func SetRouterInternalIP(ip *net.IPAddr) {
	routerInternalIP = ip
}
func SetRouterInternalName(name string) {
	routerInternalIP, _ = net.ResolveIPAddr("", name)
}

func GetRouterInternalIP() (routerIP *net.IPAddr, err error) {
	if routerInternalIP != nil {
		return routerInternalIP, nil
	}
	if nil == cbRouterInternalIPs {
		return nil, os.ErrNotExist
	}
	for _, cbRouterInternalIP := range cbRouterInternalIPs {
		routerIP, err = cbRouterInternalIP()
		if err == nil {
			routerInternalIP = routerIP
			return
		}
	}

	return nil, os.ErrNotExist
}

func registerCbRouterInternalIP(cb cbIPAddr) {
	cbRouterInternalIPs = append(cbRouterInternalIPs, cb)
}

// /////////////////////////////////////////////////////////////////////////////

func GetRouterExternalIP() (routerIP *net.IPAddr, err error) {
	if routerExternalIP != nil {
		return routerExternalIP, nil
	}
	routerExternalIP, err = getRouterExternalIP()
	return routerExternalIP, err
}

func getRouterExternalIP() (routerIP *net.IPAddr, err error) {
	routerIP, err = GetRouterExternalIPUPNP()
	if err == nil {
		return routerIP, err
	}
	return GetRouterExternalIPIpify()
}

// /////////////////////////////////////////////////////////////////////////////

func GetRouterExternalIPIpify() (publicIP *net.IPAddr, err error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("", string(ip))
}

type getExternalIPAddressEnvelope struct {
	Body getExternalIPAddressBody `xml:"Body"`
}

type getExternalIPAddressBody struct {
	Response getExternalIPAddressResponse `xml:"GetExternalIPAddressResponse"`
}

type getExternalIPAddressResponse struct {
	NewExternalIPAddress string `xml:"NewExternalIPAddress"`
}

func GetRouterExternalIPUPNP() (publicIP *net.IPAddr, err error) {
	internalIP, _ := GetRouterInternalIP()
	req, err := http.NewRequest("POST", "http://"+internalIP.String()+":49000/igdupnp/control/WANIPConn1", bytes.NewBufferString(`<?xml version='1.0' encoding='utf-8'?> 
<s:Envelope s:encodingStyle='http://schemas.xmlsoap.org/soap/encoding/' 
	xmlns:s='http://schemas.xmlsoap.org/soap/envelope/'> 
	<s:Body>
		<u:GetExternalIPAddress xmlns:u='urn:schemas-upnp-org:service:WANIPConnection:1' /> 
	</s:Body>
</s:Envelope>`))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type",
		"text/xml; charset=\"utf-8\"")
	req.Header.Add("SoapAction",
		"urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	env := getExternalIPAddressEnvelope{}
	err = xml.Unmarshal(data, &env)
	if err != nil {
		return nil, err
	}
	return net.ResolveIPAddr("", env.Body.Response.NewExternalIPAddress)
}
