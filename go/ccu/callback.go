package ccu

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/dueckminor/mypi-tools/go/xmlrpc"
)

type HttpHandler struct {
	ccuc *CcuClient
}

func resolveHostIp() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			return networkIp.IP.String()
		}
	}
	return ""
}

func (h *HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	response, err := xmlrpc.HandleMethodCall(body, h.ccuc)
	if err == nil {
		res.Write(response)
	}
}

func (ccuc *CcuClient) GetOwnIP() string {
	if len(ccuc.ownIP) == 0 {
		ccuc.ownIP = resolveHostIp()
	}
	return ccuc.ownIP
}

func (ccuc *CcuClient) StartCallbackHandler() error {
	httpHandler := &HttpHandler{}
	httpHandler.ccuc = ccuc
	ln, err := net.Listen("tcp", ":2000")
	if err != nil {
		return err
	}
	go http.Serve(ln, httpHandler)
	return ccuc.Init("http://"+ccuc.GetOwnIP()+":2000", "TESTCLIENT")
}
