package main

import (
	"flag"
	"os"
	"path"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/debug"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	err := config.InitApp(path.Join(os.Getenv("HOME"), ".mypi", "debug"))
	if err != nil {
		panic(err)
	}
	cfg, err := config.GetOrCreateConfigFile("etc/mypi-debug/mypi-debug.yml")
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	err = ginutil.ConfigureSessionCookies(r, cfg, "mypi-debug-session")
	if err != nil {
		panic(err)
	}

	err = r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	r.Use(cors.Default())

	services := debug.NewServices(r)
	services.Run()
}
