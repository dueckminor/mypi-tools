package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/util/network"
)

var (
	target              string
	portHTTP            int
	portHTTPS           int
	identityCert        string
	identityKey         string
	gatewayInternalName string
)

func init() {
	flag.StringVar(&target, "target", "", "the target (<host>:<port>)")
	flag.IntVar(&portHTTPS, "https-port", 8443, "the listening port for https")
	flag.IntVar(&portHTTP, "http-port", -1, "the listening port for http")
	flag.StringVar(&identityCert, "identity-cert", "", "the tls server certificate")
	flag.StringVar(&identityKey, "identity-key", "", "the tls server certificate")
	flag.StringVar(&gatewayInternalName, "router-name", "", "the (internal) name of the router (fritz.box, 192.168.0.1, ...)")
}

type CertConfig struct {
	CertFile    string `yaml:"cert"`
	KeyFile     string `yaml:"key"`
	cert        tls.Certificate
	hostNames   []string
	domainNames []string
}

type HostConfig struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
	Mode   string `yaml:"mode"`
}

type GatewayConfig struct {
	Certs      []*CertConfig `yaml:"certs"`
	Hosts      []*HostConfig `yaml:"hosts"`
	certByName map[string]*CertConfig
	hostByName map[string]*HostConfig
}

type connWrapper struct {
	conn      net.Conn
	cacheRead bool
	buff      []byte
}

func (w *connWrapper) Read(b []byte) (n int, err error) {
	if w.conn == nil {
		return 0, nil
	}
	n, err = w.conn.Read(b)
	if w.cacheRead && n > 0 {
		w.buff = append(w.buff, b[0:n]...)
	}
	return
}
func (w *connWrapper) Write(b []byte) (n int, err error) {
	if w.conn == nil {
		return len(b), nil
	}
	return w.conn.Write(b)
}
func (w *connWrapper) Close() error {
	if w.conn == nil {
		return nil
	}
	return w.conn.Close()
}
func (w *connWrapper) LocalAddr() net.Addr {
	if w.conn == nil {
		return nil
	}
	return w.conn.LocalAddr()
}
func (w *connWrapper) RemoteAddr() net.Addr {
	if w.conn == nil {
		return nil
	}
	return w.conn.RemoteAddr()
}
func (w *connWrapper) SetDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetDeadline(t)
}
func (w *connWrapper) SetReadDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetReadDeadline(t)
}
func (w *connWrapper) SetWriteDeadline(t time.Time) error {
	if w.conn == nil {
		return nil
	}
	return w.conn.SetWriteDeadline(t)
}

func forwardConnect(a, b net.Conn) {
	done := make(chan bool, 2)

	go func() { io.Copy(a, b); done <- true }()
	go func() { io.Copy(b, a); done <- true }()

	<-done
	<-done
}

func (gateway *GatewayConfig) createHostMap() map[string]*HostConfig {
	hostByName := make(map[string]*HostConfig)
	for _, hostConfig := range gateway.Hosts {
		hostByName[hostConfig.Name] = hostConfig
	}
	return hostByName
}

func (gatewayConfig *GatewayConfig) createCertMap() map[string]*CertConfig {
	fmt.Println("Checking certificates...")

	certByName := make(map[string]*CertConfig)

	var err error
	for _, certConfig := range gatewayConfig.Certs {
		fmt.Println(certConfig.CertFile)
		certConfig.cert, err = tls.LoadX509KeyPair(certConfig.CertFile, certConfig.KeyFile)
		if err != nil {
			panic(err)
		}
		x509Cert, err := x509.ParseCertificate(certConfig.cert.Certificate[0])
		if err != nil {
			panic(err)
		}
		for _, dnsName := range x509Cert.DNSNames {
			certByName[dnsName] = certConfig
			fmt.Println("  ", dnsName)
		}
	}
	return certByName
}

func (gateway *GatewayConfig) updateMaps() {
	gateway.certByName = gateway.createCertMap()
	gateway.hostByName = gateway.createHostMap()
}

func (gateway *GatewayConfig) getHostConfig(serverName string) *HostConfig {
	if hostConfig, ok := gateway.hostByName[serverName]; ok {
		return hostConfig
	}
	if strings.HasPrefix(serverName, "*.") {
		return nil
	}
	serverNameParts := strings.SplitN(serverName, ".", 2)
	if len(serverNameParts) != 2 {
		return nil
	}
	if hostConfig, ok := gateway.hostByName["*."+serverNameParts[1]]; ok {
		return hostConfig
	}
	return nil
}

func (gateway *GatewayConfig) getCertConfig(serverName string) *CertConfig {
	if certConfig, ok := gateway.certByName[serverName]; ok {
		return certConfig
	}
	if strings.HasPrefix(serverName, "*.") {
		return nil
	}
	serverNameParts := strings.SplitN(serverName, ".", 2)
	if len(serverNameParts) != 2 {
		return nil
	}
	if certConfig, ok := gateway.certByName["*."+serverNameParts[1]]; ok {
		return certConfig
	}
	return nil
}

func (gateway *GatewayConfig) getCertificate(serverName string) (cert tls.Certificate, err error) {
	return tls.Certificate{}, nil
}

func (gateway *GatewayConfig) handleConnection(client net.Conn) {
	clientWrapper := &connWrapper{conn: client, cacheRead: true}

	var serverName string
	var hostConfig *HostConfig
	ignoreHandshakeError := false

	tlsConn := tls.Server(clientWrapper, &tls.Config{GetConfigForClient: func(clientHelloInfo *tls.ClientHelloInfo) (*tls.Config, error) {
		clientWrapper.cacheRead = false
		serverName = clientHelloInfo.ServerName
		fmt.Println("ServerName:", serverName)

		hostConfig = gateway.getHostConfig(serverName)
		if nil == hostConfig {
			fmt.Println("-> dropped")
			return nil, os.ErrInvalid
		}
		if hostConfig.Mode == "socket" {
			clientWrapper.conn = nil
			ignoreHandshakeError = true
			return nil, os.ErrInvalid
		}

		certConfig := gateway.getCertConfig(serverName)
		if nil == certConfig {
			fmt.Println("-> dropped (have no cert")
			return nil, os.ErrInvalid
		}

		fmt.Println("->", hostConfig.Target)
		return &tls.Config{
			Certificates: []tls.Certificate{certConfig.cert},
		}, nil
		// return nil, os.ErrInvalid
	}})

	err := tlsConn.Handshake()

	if nil == hostConfig {
		fmt.Println("ServerName:", serverName, "rejected")
		return
	}

	if err != nil && !ignoreHandshakeError {
		fmt.Println("Handshake Err:", err)
		return
	}

	fmt.Println("ServerName:", serverName)

	var targetConn net.Conn
	var sourceConn net.Conn = tlsConn

	switch hostConfig.Mode {
	case "socket":
		targetConn, err = net.Dial("tcp", hostConfig.Target)
		if err == nil {
			_, err = targetConn.Write(clientWrapper.buff)
		}
		sourceConn = client
	case "tls":
		conf := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         serverName,
		}
		targetConn, err = tls.Dial("tcp", hostConfig.Target, conf)
	default:
		targetConn, err = net.Dial("tcp", hostConfig.Target)
	}

	if err != nil {
		fmt.Println("Dial Err:", err)
		return
	}

	forwardConnect(sourceConn, targetConn)
}

func main() {
	flag.Parse()

	if len(gatewayInternalName) > 0 {
		network.SetRouterInternalName(gatewayInternalName)
	}

	ip, _ := network.GetRouterInternalIP()
	fmt.Println("Router internal IP:", ip)
	ip, _ = network.GetRouterExternalIP()
	fmt.Println("Router external IP:", ip)

	var gatewayConfig *GatewayConfig

	nArgs := flag.CommandLine.NArg()
	if nArgs > 1 {
		panic("To many args specified")
	}
	if nArgs == 1 {
		config.ReadYAML(flag.CommandLine.Arg(0), &gatewayConfig)
	} else {
		gatewayConfig = &GatewayConfig{}
	}

	if len(identityCert) > 0 && len(identityKey) > 0 {
		gatewayConfig.Certs = append(gatewayConfig.Certs, &CertConfig{
			CertFile: identityCert,
			KeyFile:  identityKey,
		})
	}

	gatewayConfig.updateMaps()

	signals := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for range signals {
			fmt.Println("\nReceived an interrupt, stopping...")
			stop <- true
		}
	}()

	if portHTTP > 0 {
		http.HandleFunc("/.well-known/acme-challenge", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "ACME-CHALLENGE: %q", html.EscapeString(r.Host))
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		go func() {
			http.ListenAndServe(fmt.Sprintf(":%d", portHTTP), nil)
		}()
	}

	incoming, err := net.Listen("tcp", fmt.Sprintf(":%d", portHTTPS))
	if err != nil {
		log.Fatalf("could not start server on %d: %v", portHTTPS, err)
	}
	fmt.Printf("server running on %d\n", portHTTPS)

	go func() {
		for {
			client, err := incoming.Accept()
			if err != nil {
				log.Fatal("could not accept client connection", err)
			}
			go func() {
				defer client.Close()
				fmt.Printf("client '%v' connected!\n", client.RemoteAddr())
				gatewayConfig.handleConnection(client)
				fmt.Printf("client '%v' disconnected!\n", client.RemoteAddr())
			}()
		}
	}()

	<-stop
}
