package ccu

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dueckminor/mypi-tools/go/xmlrpc"
)

type HttpHandler struct {
	ccuc *CcuClientImpl
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

func (ccuc *CcuClientImpl) GetOwnIP() string {
	if len(ccuc.ownIP) == 0 {
		ccuc.ownIP = resolveHostIp()
	}
	return ccuc.ownIP
}

func (ccuc *CcuClientImpl) StartCallbackHandler() error {
	httpHandler := &HttpHandler{}
	httpHandler.ccuc = ccuc

	port := "0"
	parsedURI, err := url.Parse(ccuc.uri)
	if err == nil {
		port = parsedURI.Port()
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	port = strconv.FormatInt(int64(ln.Addr().(*net.TCPAddr).Port), 10)

	ownURL := "http://" + ccuc.GetOwnIP() + ":" + port
	ownID := "MYPI-" + ccuc.GetOwnIP() + "-" + port

	go http.Serve(ln, httpHandler)
	err = ccuc.Init(ownURL, ownID)
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(15 * time.Minute)
			fmt.Println("Init again...")
			ccuc.Init(ownURL, ownID)
		}
	}()

	return err
}
