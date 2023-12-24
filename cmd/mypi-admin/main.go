package main

import (
	"flag"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/pki"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"

	// provide only the http rest api
	_ "github.com/dueckminor/mypi-tools/go/restapi/http"
)

var (
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		err := config.InitApp(*mypiRoot)
		if err != nil {
			panic(err)
		}
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

	restapi.Run(r)
}
