package main

import (
	"fmt"

	"github.com/dueckminor/mypi-admin/go/config"
	"github.com/dueckminor/mypi-admin/go/debug"
	"github.com/dueckminor/mypi-admin/go/webhandler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.GetString("config", "root"))

	r := gin.Default()
	webhandler.SetupEndpoints(r)

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	debug.SetupProxy(r, "http://localhost:8081")

	r.Run()
}
