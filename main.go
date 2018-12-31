package main

import (
	"github.com/dueckminor/mypi-api/go/debug"
	"github.com/dueckminor/mypi-api/go/pki"
	"github.com/dueckminor/mypi-api/go/webhandler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	pki.Setup()

	r := gin.Default()
	webhandler.SetupEndpoints(r)

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	debug.SetupProxy(r, "http://localhost:8081")

	r.Run()
}
