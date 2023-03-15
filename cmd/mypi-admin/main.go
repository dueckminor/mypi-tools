package main

import (
	"flag"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/pki"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	webpackDebug = flag.String("webpack-debug", "", "The debug URI")
	port         = flag.Int("port", 8080, "The port")
	mypiRoot     = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}
}

func main() {
	pki.Setup()

	r := gin.Default()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)
	if err != nil {
		panic(err)
	}

	if len(*webpackDebug) > 0 {
		r.Use(ginutil.SingleHostReverseProxy(*webpackDebug))
	} else {
		r.Use(static.ServeRoot("/", "./dist"))
	}

	panic(r.Run(":" + strconv.Itoa(*port)))
}
