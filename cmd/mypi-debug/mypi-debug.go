package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	proxy := ginutil.SingleHostReverseProxy("http://localhost:8081", "external-hostname")

	if *host {
		r.Use(proxy)
	} else {
		r.Use(func(c *gin.Context) {
			if c.IsAborted() {
				return
			}
			fmt.Println("Header: ", c.Request.Header)

			proxy(c)
		})
	}

	panic(r.Run(":" + strconv.Itoa(*port)))
}
