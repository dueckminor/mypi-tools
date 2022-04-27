package main

import (
	"flag"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	port     = flag.Int("port", 8080, "The port")
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}
}

func main() {
	r := gin.Default()

	r.Use(ginutil.SingleHostReverseProxy("http://localhost:8081"))

	panic(r.Run(":" + strconv.Itoa(*port)))
}
