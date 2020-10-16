package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dueckminor/mypi-tools/go/fdisk"
	"github.com/dueckminor/mypi-tools/go/gotty/cachedcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/localcommand"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdmakesd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdsetup"
)

func forwardToMypiHost(c *gin.Context, f func(c *gin.Context)) {
	host := c.Params.ByName("host")
	req := c.Request

	local := (host == "localhost") || (host == "127.0.0.1")
	if !local {
		if osHost, err := os.Hostname(); err == nil {
			local = osHost == host
		}
	}
	if !local {
		urlHostParts := strings.Split(req.Host, ":")
		if len(urlHostParts) > 0 {
			local = urlHostParts[0] == host
		}
	}

	if local {
		f(c)
	} else {
		newHost := host + ":8080"
		targetURI, _ := url.ParseRequestURI("http://" + newHost)
		proxy := httputil.NewSingleHostReverseProxy(targetURI)
		proxy.ServeHTTP(c.Writer, req)
	}
}

func makeForwardToMypiHost(f func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		forwardToMypiHost(c, f)
	}
}

func flashLED() {
	const led0 = "/sys/class/leds/led0/brightness"
	const led1 = "/sys/class/leds/led1/brightness"
	if !util.FileExists(led0) {
		return
	}
	for {
		ioutil.WriteFile(led0, []byte("0"), os.ModePerm)
		ioutil.WriteFile(led1, []byte("1"), os.ModePerm)
		time.Sleep(time.Second)
		ioutil.WriteFile(led0, []byte("1"), os.ModePerm)
		ioutil.WriteFile(led1, []byte("0"), os.ModePerm)
		time.Sleep(time.Second)
	}
}

func main() {
	if len(os.Args) > 1 {
		if command, err := cmd.GetCommand(os.Args[1]); err == nil {
			var parsedArgs interface{}
			if len(os.Args) == 3 && os.Args[2] == "@" {
				var data []byte
				data, err = ioutil.ReadAll(os.Stdin)
				parsedArgs, err = command.UnmarshalArgs(data)
			} else {
				parsedArgs, err = command.ParseArgs(os.Args[2:])
			}
			if err == nil {
				err = command.Execute(parsedArgs)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	if util.FileExists("/sbin/setup-alpine") {
		go flashLED()
		if _, err := cmd.GetCommand("setup"); err == nil {
			c := exec.Command(os.Args[0], "setup", "@")
			err = cachedcommand.AttachProcess("setup", c)
			c.Stdin = bytes.NewReader([]byte("{}"))
			go func() {
				c.Start()
				c.Wait()
			}()
		}
	}

	flag.Parse()

	r := gin.Default()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	r.GET("/api/certificates", func(c *gin.Context) {

	})

	r.GET("/api/hosts", func(c *gin.Context) {

	})

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

	r.GET("/api/hosts/:host/actions/:action/webtty", makeForwardToMypiHost(func(c *gin.Context) {
		action := c.Params.ByName("action")
		if _, err := cmd.GetCommand(action); err == nil {
			factory, err := cachedcommand.NewFactory(action)
			if err == nil {
				ginhandler.Handler(c, factory)
			}
		} else {
			fmt.Fprintln(os.Stderr, "action '"+action+"' not found")
		}
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
			fmt.Println(parsedArgs)

			data, err = json.Marshal(parsedArgs)

			c := exec.Command(os.Args[0], action, "@")
			err = cachedcommand.AttachProcess(action, c)
			c.Stdin = bytes.NewReader(data)

			go func() {
				c.Start()
				c.Wait()
			}()
		}
	}))

	if err != nil {
		panic(err)
	}

	restapi.Run(r)
}
