package main

import (
	"flag"
	"os"
	"os/exec"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	// authURI   string
	webpackDebug = flag.String("webpack-debug", "", "The debug URI")
	port         = flag.Int("port", 8080, "The port")
	execDebug    = flag.String("exec", "", "start process")
)

func init() {
	flag.Parse()
}

func main() {
	if len(*execDebug) > 0 {
		cmd := exec.Command(*execDebug)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		defer cmd.Wait()
	}

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
