package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/dueckminor/mypi-tools/cmd/mypi-setup/impl"
	"github.com/dueckminor/mypi-tools/go/gotty/cachedcommand"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dueckminor/mypi-tools/go/cmd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdmakesd"
	_ "github.com/dueckminor/mypi-tools/go/cmd/cmdsetup"
)

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
	r.Use(cors.Default())

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	r.GET("/api/certificates", impl.GetCertificates)
	r.GET("/api/hosts", impl.GetHosts)

	r.GET("/api/hosts/:host/terminal/webtty", impl.MakeForwardToHost(impl.GetWebTTY))
	r.GET("/api/hosts/:host/disks", impl.MakeForwardToHost(impl.GetDisks))

	r.GET("/api/hosts/:host/actions/:action/webtty", impl.MakeForwardToHost(impl.GetActionWebTTY))

	r.POST("/api/hosts/:host/actions/:action", impl.MakeForwardToHost(impl.GetAction))

	if err != nil {
		panic(err)
	}

	restapi.Run(r)
}
