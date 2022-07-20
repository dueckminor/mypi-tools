package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/debug"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmddebug"
)

var (
	// authURI   string
	port     = flag.Int("port", 8080, "The port")
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
	host     = flag.Bool("host", false, "enable the host mode (Used in the mypi-debug container on the mypi host)")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}
}

func main() {
	if len(os.Args) > 1 {
		done, err := cmd.UnmarshalAndExecute(os.Args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if done {
			os.Exit(0)
		}
	}

	arch, err := debug.RunOnMypiCaptureOutput("uname", "-m")
	if err != nil {
		panic(err)
	}
	arch = strings.Trim(arch, "\r\n")
	goarch := ""
	switch arch {
	case "aarch64":
		goarch = "arm64"
	default:
		panic("arch not supported")
	}
	fmt.Println("Running on:", goarch)
	debug.RunLocal("")

	r := gin.Default()
	proxy := ginutil.SingleHostReverseProxy("http://localhost:8081", "external-hostname")

	if *host {
		r.Use(proxy)
	} else {
		proxyAuth := ginutil.SingleHostReverseProxy("http://localhost:9100")
		proxyDebugWeb := ginutil.SingleHostReverseProxy("http://localhost:8081")

		r.Use(func(c *gin.Context) {
			if c.IsAborted() {
				return
			}
			fmt.Println("Host: ", c.Request.Host)
			if c.Request.Host == "auth.rh94-dev.dueckminor.de" {
				proxyAuth(c)
			} else if c.Request.Host == "debug.rh94-dev.dueckminor.de" {
				proxyDebugWeb(c)
			} else {
				proxy(c)
			}
		})
	}

	panic(r.Run(":" + strconv.Itoa(*port)))
}
