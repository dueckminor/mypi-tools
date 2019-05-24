package main

import (
	"crypto/rand"
	"flag"
	"path"
	"strconv"

	"github.com/dueckminor/mypi-api/go/auth"
	"github.com/dueckminor/mypi-api/go/config"
	"github.com/dueckminor/mypi-api/go/ginutil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Hostname string `yaml: "hostname"`
	Target   string `yaml: "target"`
}

type RouterConfig struct {
	Routes []RouteConfig `yaml: "routes"`
}

var (
	port         int
	targetURI    string
	mypiRoot     string
	ac           *auth.AuthClient
	routerConfig RouterConfig
)

func init() {
	ac = new(auth.AuthClient)

	flag.StringVar(&ac.AuthURI, "auth", "", "The URI of the authentication server")
	flag.IntVar(&port, "port", 8080, "The port")
	flag.StringVar(&targetURI, "target", "", "The target URI")
	flag.StringVar(&ac.ClientID, "client-id", "", "The client ID")
	flag.StringVar(&ac.ClientSecret, "client-secret", "", "The client secret")
	flag.StringVar(&mypiRoot, "mypi-root", "", "The root of the mypi filesystem")

	flag.Parse()

	nArgs := flag.CommandLine.NArg()
	if nArgs > 1 {
		panic("To many args specified")
	}
	if nArgs == 1 {
		config.ReadYAML(flag.CommandLine.Arg(0), &routerConfig)
	}

	if len(mypiRoot) > 0 {
		config.InitApp(mypiRoot)
	}

	if len(ac.ClientID) == 0 {
		panic("No client id specified")
	}

	if len(ac.ClientSecret) == 0 {
		clientConfig, err := config.ReadConfigFile(path.Join("etc/auth/clients", ac.ClientID+".yml"))
		if err != nil {
			panic(err)
		}
		ac.ClientSecret = clientConfig.GetString("client_secret")
		if len(ac.ClientSecret) == 0 {
			panic("No client secret specified")
		}
		ac.ServerKey, err = config.StringToRSAPublicKey(clientConfig.GetString("server_key"))
		if err != nil {
			panic(err)
		}
	}

	if len(targetURI) > 0 && len(routerConfig.Routes) > 0 {
		panic("It's not allowed to specify a config file with routes and a target")
	}

}

func main() {
	r := gin.Default()

	key := make([]byte, 64)
	rand.Read(key)

	store := memstore.NewStore([]byte(key))
	r.Use(sessions.Sessions("MYPI_ROUTER_SESSION", store))

	ac.RegisterHandler(r)

	if len(targetURI) > 0 {
		r.Use(ginutil.SingleHostReverseProxy(targetURI))
	}

	for _, route := range routerConfig.Routes {
		r.Use(ginutil.OnHostname(route.Hostname, ginutil.SingleHostReverseProxy(route.Target)))
	}

	panic(r.Run(":" + strconv.Itoa(port)))
}
