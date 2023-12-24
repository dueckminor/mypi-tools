package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/dueckminor/mypi-tools/go/auth"
	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/dueckminor/mypi-tools/go/util/network"
	"github.com/dueckminor/mypi-tools/go/util/panic"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var (
	target              string
	portHTTP            int
	portHTTPS           int
	gatewayInternalName string
	store               memstore.Store

	localhostOnly = flag.Bool("localhost-only", false, "Listen on localhost only")
	mypiRoot      = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.StringVar(&target, "target", "", "the target (<host>:<port>)")
	flag.IntVar(&portHTTPS, "https-port", 8443, "the listening port for https")
	flag.IntVar(&portHTTP, "http-port", -1, "the listening port for http")
	flag.StringVar(&gatewayInternalName, "router-name", "", "the (internal) name of the router (fritz.box, 192.168.0.1, ...)")
}

type CertConfig struct {
	CertFile string `yaml:"cert"`
	KeyFile  string `yaml:"key"`
	cert     tls.Certificate
}

type HostConfig struct {
	Name    string   `yaml:"name"`
	Target  string   `yaml:"target"`
	Mode    string   `yaml:"mode"`
	Options []string `yaml:"options"`
}

func (h *HostConfig) hasOption(option string) bool {
	for _, o := range h.Options {
		if o == option {
			return true
		}
	}
	return false
}

// /////////////////////////////////////////////////////////////////////////////
type HostImpl interface {
	String() string
	HandleConnection(conn net.Conn)
}

// /////////////////////////////////////////////////////////////////////////////

type HostImplBase struct {
	HostConfig
}

func (h *HostImplBase) String() string {
	return h.Target
}

// /////////////////////////////////////////////////////////////////////////////

type HostImplSocket struct {
	HostImplBase
}

func (h *HostImplSocket) HandleConnection(conn net.Conn) {
	panic.Panic("HandleConnectionSocket expected")
}

func (h *HostImplSocket) HandleConnectionSocket(conn net.Conn, buf []byte) {
	defer conn.Close()
	targetConn, err := net.Dial("tcp", h.Target)
	if err != nil {
		return
	}
	_, err = targetConn.Write(buf)
	if err != nil {
		return
	}
	forwardConnect(conn, targetConn)
}

func NewHostImplSocket(hostConfig *HostConfig) *HostImplSocket {
	h := new(HostImplSocket)
	h.HostConfig = *hostConfig
	return h
}

// /////////////////////////////////////////////////////////////////////////////

type HostImplTLS struct {
	HostImplBase
}

func (h *HostImplTLS) HandleConnection(conn net.Conn) {
	defer conn.Close()
	conf := &tls.Config{
		InsecureSkipVerify: h.hasOption("insecure"),
		ServerName:         h.Target,
	}
	targetConn, err := tls.Dial("tcp", h.Target, conf)
	if err != nil {
		fmt.Println("Dial Err:", err)
		return
	}
	forwardConnect(conn, targetConn)
}

func NewHostImplTLS(hostConfig *HostConfig) *HostImplTLS {
	h := new(HostImplTLS)
	h.HostConfig = *hostConfig
	return h
}

// /////////////////////////////////////////////////////////////////////////////

type HostImplPort struct {
	HostImplBase
}

func (h *HostImplPort) HandleConnection(conn net.Conn) {
	defer conn.Close()
	targetConn, err := net.Dial("tcp", h.Target)
	if err != nil {
		fmt.Println("Dial Err:", err)
		return
	}
	forwardConnect(conn, targetConn)
}

func NewHostImplPort(hostConfig *HostConfig) *HostImplPort {
	h := new(HostImplPort)
	h.HostConfig = *hostConfig
	return h
}

// /////////////////////////////////////////////////////////////////////////////

type HostImplReverseProxy struct {
	HostImplBase
	listener *Listener
	r        *gin.Engine
}

func (h *HostImplReverseProxy) HandleConnection(conn net.Conn) {
	// no need to do this here: defer conn.Close()
	// the connection will be closed by the gin.Engine
	h.listener.Connections <- conn
}

func NewHostImplReverseProxy(hostConfig *HostConfig, uri string, ac *auth.AuthClient) *HostImplReverseProxy {
	h := new(HostImplReverseProxy)
	h.HostConfig = *hostConfig
	h.r = gin.Default()

	if hostConfig.hasOption("external-hostname") {
		h.Options = []string{"external-hostname"}
	}

	if ac != nil {
		h.r.Use(sessions.Sessions("MYPI_ROUTER_SESSION", store))
		ac.RegisterHandler(h.r)
	}

	h.listener = MakeListener()
	go h.r.RunListener(h.listener) // nolint:errcheck
	h.r.Use(ginutil.SingleHostReverseProxy(uri, h.Options...))
	return h
}

func NewHostImplHTTP(hostConfig *HostConfig, ac *auth.AuthClient) *HostImplReverseProxy {
	return NewHostImplReverseProxy(hostConfig, "http://"+hostConfig.Target, ac)
}

func NewHostImplHTTPS(hostConfig *HostConfig, ac *auth.AuthClient) *HostImplReverseProxy {
	return NewHostImplReverseProxy(hostConfig, "https://"+hostConfig.Target, ac)
}

// /////////////////////////////////////////////////////////////////////////////

type AuthConfig struct {
	URI          string `yaml:"uri"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	ServerKey    string `yaml:"server_key"`
}

// /////////////////////////////////////////////////////////////////////////////

type GatewayConfig struct {
	Certs      []*CertConfig `yaml:"certs"`
	Hosts      []*HostConfig `yaml:"hosts"`
	Auth       AuthConfig    `yaml:"auth"`
	certByName map[string]*CertConfig
	hostByName map[string]HostImpl
	configFile string
	mutex      sync.RWMutex
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

func forwardConnect(client, server net.Conn) {
	done := make(chan bool, 2)

	go func() {
		// when the server closes the connection,
		// it's no longer necessary to send something
		// -> lets close the client connection
		defer client.Close()

		io.Copy(client, server) // nolint: errcheck
		done <- true
	}()

	go func() {
		io.Copy(server, client) // nolint: errcheck
		done <- true
	}()

	<-done
	<-done
}

func (gateway *GatewayConfig) createHostMap() map[string]HostImpl {
	ac := gateway.GetAuthClient()
	hostByName := make(map[string]HostImpl)
	for _, hostConfig := range gateway.Hosts {
		var impl HostImpl
		switch hostConfig.Mode {
		case "http":
			if util.StringsContains(hostConfig.Options, "auth") {
				impl = NewHostImplHTTP(hostConfig, ac)
			} else {
				impl = NewHostImplHTTP(hostConfig, nil)
			}
		case "https":
			if util.StringsContains(hostConfig.Options, "auth") {
				impl = NewHostImplHTTPS(hostConfig, ac)
			} else {
				impl = NewHostImplHTTPS(hostConfig, nil)
			}
		case "socket":
			impl = NewHostImplSocket(hostConfig)
		case "tls":
			impl = NewHostImplTLS(hostConfig)
		default:
			impl = NewHostImplPort(hostConfig)
		}

		hostByName[hostConfig.Name] = impl
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
		panic.OnError(err)
		x509Cert, err := x509.ParseCertificate(certConfig.cert.Certificate[0])
		panic.OnError(err)
		for _, dnsName := range x509Cert.DNSNames {
			certByName[dnsName] = certConfig
			fmt.Println("  ", dnsName)
		}
	}
	return certByName
}

func (gateway *GatewayConfig) updateMaps() {
	gateway.setCertMap(gateway.createCertMap())
	gateway.setHostMap(gateway.createHostMap())
}

func (gateway *GatewayConfig) setCertMap(certByName map[string]*CertConfig) {
	gateway.mutex.Lock()
	defer gateway.mutex.Unlock()
	gateway.certByName = certByName
}

func (gateway *GatewayConfig) setHostMap(hostByName map[string]HostImpl) {
	gateway.mutex.Lock()
	defer gateway.mutex.Unlock()
	gateway.hostByName = hostByName
}

func (gateway *GatewayConfig) startWatcher() (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Name == gateway.configFile {
					err = gateway.loadConfig()
					if err != nil {
						return
					}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					gateway.setCertMap(gateway.createCertMap())
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _, certConfig := range gateway.Certs {
		err = watcher.Add(certConfig.CertFile)
		if err != nil {
			return err
		}
	}

	if len(gateway.configFile) > 0 {
		err = watcher.Add(gateway.configFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gateway *GatewayConfig) loadConfig() (err error) {
	if len(gateway.configFile) > 0 {
		var newConfig *GatewayConfig
		err = config.ReadYAML(&newConfig, gateway.configFile)
		if err != nil {
			return err
		}
		if newConfig != nil {
			gateway.Certs = newConfig.Certs
			gateway.Hosts = newConfig.Hosts
			gateway.Auth = newConfig.Auth
		}
	}

	gateway.updateMaps()
	return nil
}

func (gateway *GatewayConfig) getHostImpl(serverName string) HostImpl {
	gateway.mutex.RLock()
	defer gateway.mutex.RUnlock()

	if hostImpl, ok := gateway.hostByName[serverName]; ok {
		return hostImpl
	}
	if strings.HasPrefix(serverName, "*.") {
		return nil
	}
	serverNameParts := strings.SplitN(serverName, ".", 2)
	if len(serverNameParts) != 2 {
		return nil
	}
	if hostImpl, ok := gateway.hostByName["*."+serverNameParts[1]]; ok {
		return hostImpl
	}
	return nil
}

func (gateway *GatewayConfig) getCertConfig(serverName string) *CertConfig {
	gateway.mutex.RLock()
	defer gateway.mutex.RUnlock()

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

type Listener struct {
	Connections chan net.Conn
}

func (l *Listener) Accept() (net.Conn, error) {
	conn := <-l.Connections
	return conn, nil
}
func (l *Listener) Close() error {
	return nil
}
func (l *Listener) Addr() net.Addr {
	return nil
}

func MakeListener() *Listener {
	l := &Listener{}
	l.Connections = make(chan net.Conn)
	return l
}

func (gateway *GatewayConfig) handleConnection(client net.Conn) {
	clientWrapper := &connWrapper{conn: client, cacheRead: true}

	closeConn := true
	defer func() {
		if closeConn {
			client.Close()
		}
	}()

	var serverName string
	var hostImpl HostImpl
	var hostImplSocket *HostImplSocket

	tlsConn := tls.Server(clientWrapper, &tls.Config{GetConfigForClient: func(clientHelloInfo *tls.ClientHelloInfo) (*tls.Config, error) {
		clientWrapper.cacheRead = false
		serverName = clientHelloInfo.ServerName
		fmt.Println("ServerName:", serverName)

		hostImpl = gateway.getHostImpl(serverName)
		if nil == hostImpl {
			fmt.Println("-> dropped")
			return nil, os.ErrInvalid
		}
		var ok bool
		if hostImplSocket, ok = hostImpl.(*HostImplSocket); ok {
			// from now on the connection is handled by hostImplSocket
			clientWrapper.conn = nil
			return nil, os.ErrInvalid
		}

		certConfig := gateway.getCertConfig(serverName)
		if nil == certConfig {
			fmt.Println("-> dropped (have no cert")
			return nil, os.ErrInvalid
		}

		fmt.Println("->", hostImpl.String())
		return &tls.Config{
			Certificates: []tls.Certificate{certConfig.cert},
		}, nil
		// return nil, os.ErrInvalid
	}})

	err := tlsConn.Handshake()

	if nil == hostImpl {
		fmt.Println("ServerName:", serverName, "rejected")
		return
	}

	if err != nil && hostImplSocket == nil {
		fmt.Println("Handshake Err:", err)
		return
	}

	fmt.Println("ServerName:", serverName)

	// from now on hostImpl is responsible to close the connection
	closeConn = false

	if hostImplSocket != nil {
		hostImplSocket.HandleConnectionSocket(client, clientWrapper.buff)
	} else {
		hostImpl.HandleConnection(tlsConn)
	}
}

func (c *GatewayConfig) GetAuthClient() *auth.AuthClient {
	if len(c.Auth.URI) == 0 {
		return nil
	}

	ac := new(auth.AuthClient)
	ac.AuthURI = c.Auth.URI
	ac.ClientID = c.Auth.ClientID
	ac.ClientSecret = c.Auth.ClientSecret

	ServerKey := c.Auth.ServerKey

	if len(ServerKey) == 0 || len(ac.ClientSecret) == 0 {
		clientConfig, err := config.ReadConfigFile("etc/mypi-auth/clients", ac.ClientID+".yml")
		panic.OnError(err)
		if len(ac.ClientSecret) == 0 {
			ac.ClientSecret = clientConfig.GetString("client_secret")
			panic.OnCond(ac.ClientSecret == "", "No client secret specified")
		}
		ServerKey = clientConfig.GetString("server_key")
	}

	var err error
	ac.ServerKey, err = config.StringToRSAPublicKey(ServerKey)
	panic.OnError(err)

	return ac
}

func main() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		err := config.InitApp(*mypiRoot)
		panic.OnError(err)
	}

	if len(gatewayInternalName) > 0 {
		network.SetRouterInternalName(gatewayInternalName)
	}

	fmt.Println("Root-Dir:", config.GetRoot())

	ip, _ := network.GetRouterInternalIP()
	fmt.Println("Router internal IP:", ip)
	ip, _ = network.GetRouterExternalIP()
	fmt.Println("Router external IP:", ip)

	key := make([]byte, 64)
	_, err := rand.Read(key)
	panic.OnError(err)
	store = memstore.NewStore([]byte(key))

	var gatewayConfig *GatewayConfig

	nArgs := flag.CommandLine.NArg()
	panic.OnCond(nArgs > 1, "To many args specified")
	gatewayConfig = &GatewayConfig{}
	configFile := config.GetFilename("etc/mypi-router/router.yml")

	fmt.Println("Using config file:", configFile)

	if nArgs == 1 {
		configFile = flag.CommandLine.Arg(0)
	}
	gatewayConfig.configFile = configFile

	err = gatewayConfig.loadConfig()
	panic.OnError(err)
	err = gatewayConfig.startWatcher()
	panic.OnError(err)

	signals := make(chan os.Signal, 1)
	stop := make(chan bool)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for range signals {
			fmt.Println("\nReceived an interrupt, stopping...")
			stop <- true
		}
	}()

	host := ""
	if *localhostOnly {
		host = "localhost"
	}

	if portHTTP > 0 {
		http.HandleFunc("/.well-known/acme-challenge/", func(w http.ResponseWriter, r *http.Request) {
			_, token := path.Split(r.URL.Path)
			if !util.FileIsSafe(r.Host) || !util.FileIsSafe(token) {
				return
			}
			acmeChallenge := path.Join("/etc/letsencrypt/acme-challenge", r.Host, token)
			if util.FileExists(acmeChallenge) {
				if stream, err := os.Open(acmeChallenge); err == nil {
					io.Copy(w, stream) // nolint:errcheck
				}
			}
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})
		go func() {
			http.ListenAndServe(fmt.Sprintf("%s:%d", host, portHTTP), nil) // nolint:errcheck
		}()
	}

	incoming, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, portHTTPS))
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
				remoteAddr := client.RemoteAddr()
				fmt.Printf("client '%v' connected!\n", remoteAddr)
				gatewayConfig.handleConnection(client)
				fmt.Printf("client '%v' disconnected!\n", remoteAddr)
			}()
		}
	}()

	<-stop
}
