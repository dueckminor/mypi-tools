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
	config.InitApp(path.Join(os.Getenv("HOME"), ".mypi", "debug"))
	cfg, err := config.GetOrCreateConfigFile("mypi-debug.yml")
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	err = ginutil.ConfigureSessionCookies(r, cfg)
	if err != nil {
		panic(err)
	}

	r.SetTrustedProxies(nil)
	r.Use(cors.Default())

	services := debug.NewServices(r)
	services.Run()
}
