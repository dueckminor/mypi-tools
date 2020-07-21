package main

import (
	"crypto/rand"
	"flag"
	"path"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/auth"
	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Hostname string `yaml:"hostname"`
	Target   string `yaml:"target"`
	Insecure bool   `yaml:"insecure"`
}
type RedirectConfig struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type RouterConfig struct {
	Routes    []RouteConfig    `yaml:"routes,omitempty"`
	Redirects []RedirectConfig `yaml:"redirects,omitempty"`
}

var (
	port           int
	targetURI      string
	redirectSource string
	redirectTarget string
	mypiRoot       string
	staticRoot     string
	ac             *auth.AuthClient
	routerConfig   RouterConfig
)

func init() {
	ac = new(auth.AuthClient)

	flag.StringVar(&ac.AuthURI, "auth", "", "The URI of the authentication server")
	flag.IntVar(&port, "port", 8080, "The port")
	flag.StringVar(&targetURI, "target", "", "The target URI")
	flag.StringVar(&staticRoot, "static-root", "", "The target URI")
	flag.StringVar(&redirectSource, "redirect-source", "", "The redirect source")
	flag.StringVar(&redirectTarget, "redirect-target", "", "The redirect target")
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
}

func main() {
	r := gin.Default()

	key := make([]byte, 64)
	rand.Read(key)

	store := memstore.NewStore([]byte(key))
	r.Use(sessions.Sessions("MYPI_ROUTER_SESSION", store))

	ac.RegisterHandler(r)

	if len(staticRoot) > 0 {
		r.Use(static.ServeRoot("/", staticRoot))
	}

	for _, redirect := range routerConfig.Redirects {
		r.GET(redirect.Source, ginutil.Redirect(redirect.Target))
	}

	if len(redirectSource) > 0 {
		r.GET(redirectSource, ginutil.Redirect(redirectTarget))
	}

	for _, route := range routerConfig.Routes {
		if len(route.Hostname) == 0 || route.Hostname == "*" {
			if (len(targetURI)) > 0 {
				panic("more than one default route speciefed")
			}
			targetURI = route.Target
			continue
		}
		if route.Insecure {
			r.Use(ginutil.OnHostname(route.Hostname, ginutil.SingleHostReverseProxy(route.Target, "insecure")))
		} else {
			r.Use(ginutil.OnHostname(route.Hostname, ginutil.SingleHostReverseProxy(route.Target)))
		}

	}

	if len(targetURI) > 0 {
		r.Use(ginutil.SingleHostReverseProxy(targetURI))
	}

	panic(r.Run(":" + strconv.Itoa(port)))
}
