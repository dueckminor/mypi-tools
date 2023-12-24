package main

import (
	"flag"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/gin-contrib/static"
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
		err := config.InitApp(*mypiRoot)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	r := gin.Default()

	r.Use(static.ServeRoot("/", "/opt/www"))

	panic(r.Run(":" + strconv.Itoa(*port)))
}
