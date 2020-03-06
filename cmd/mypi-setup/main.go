package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/gotty/localcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdmakesd"
)

func forwardToMypiHost(c *gin.Context, f func(c *gin.Context)) {
	host := c.Params.ByName("host")
	req := c.Request

	local := (host == "localhost") || (host == "127.0.0.1")
	if !local {
		urlHostParts := strings.Split(req.Host, ":")
		if len(urlHostParts) > 0 {
			local = urlHostParts[0] == host
		}
	}

	if local {
		f(c)
	} else {
		targetURI, _ := url.ParseRequestURI("http://" + host + ":8080")
		proxy := httputil.NewSingleHostReverseProxy(targetURI)
		proxy.ServeHTTP(c.Writer, req)
	}
}

func makeForwardToMypiHost(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		forwardToMypiHost(c, f)
	}
}

func main() {
	if len(os.Args) > 1 {
		if cmd.IsAvailable(os.Args[1]) {
			cmd.ExecuteCmdline(os.Args[1], os.Args[2:]...)
			return
		}
	}

	flag.Parse()

	r := gin.Default()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	r.GET("/api/hosts/:host/terminal/webtty", makeForwardToMypiHost(func(c *gin.Context) {
		sh := "zsh"
		if runtime.GOOS == "linux" && runtime.GOARCH == "arm64" {
			sh = "ash"
		}
		factory, err := localcommand.NewFactory(sh, []string{"-l"}, &localcommand.Options{})
		if err == nil {
			ginhandler.Handler(c, factory)
		}
	}))

	r.GET("/api/hosts/:host/disks", makeForwardToMypiHost(func(c *gin.Context) {
		type DiskInfo struct {
			Name string `json:"name,omitempty"`
			Size int64  `json:"size"`
		}
		diskInfos := []DiskInfo{}
		disks, err := fdisk.GetDisks()
		if err == nil {
			for _, disk := range disks {
				if disk.IsRemovable() {
					diskInfos = append(diskInfos, DiskInfo{
						Name: disk.GetDeviceName(),
						Size: disk.GetSize(),
					})
				}
			}
		}
		data, err := json.Marshal(diskInfos)
		c.Data(200, "application/json", data)
	}))

	r.POST("/api/hosts/:host/actions/:action", makeForwardToMypiHost(func(c *gin.Context) {
		action := c.Params.ByName("action")
		if command, err := cmd.GetCommand(action); err == nil {
			data, err := c.GetRawData()
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
			parsedArgs, err := command.UnmarshalArgs(data)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
			command.Execute(parsedArgs)
		}
	}))

	if err != nil {
		panic(err)
	}

	restapi.Run(r)
}
