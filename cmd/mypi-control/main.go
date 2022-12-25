package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/dueckminor/mypi-tools/go/cmd"
	"github.com/dueckminor/mypi-tools/go/gotty/cachedcommand"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/util"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

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
		done, err := cmd.UnmarshalAndExecute(os.Args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if done {
			os.Exit(0)
		}
	}

	flag.Parse()

	if util.IsRunningOnMypi() {
		fmt.Println("Running on MYPI...")
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
	} else {
		fmt.Println("Not running on MYPI...")
	}

	r := gin.Default()

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())

		go func() {
			for {
				cpuUsage, err := cpu.Percent(0, false)
				if err == nil {
					s.Emit("stats/rpi/cpu", strconv.FormatFloat(cpuUsage[0], 'f', 2, 64))
				}
				v, err := mem.VirtualMemory()
				if err == nil {
					s.Emit("stats/rpi/mem", strconv.FormatFloat(v.UsedPercent, 'f', 2, 64))
					if v.SwapTotal > 0 {
						s.Emit("stats/rpi/swap", strconv.FormatFloat(100.0-100.0*float64(v.SwapFree)/float64(v.SwapTotal), 'f', 2, 64))
					}
				}
				time.Sleep(time.Second * 1)
			}
		}()

		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("stats/rpi/cpu", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("stats", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	if err != nil {
		panic(err)
	}

	r.GET("/ws/*any", gin.WrapH(server))
	r.POST("/ws/*any", gin.WrapH(server))

	r.GET("/api/hosts/:host/terminal/webtty", webhandler.MakeForwardToHost(webhandler.GetWebTTY))
	r.GET("/api/hosts/:host/actions", webhandler.MakeForwardToHost(webhandler.GetActions))
	r.GET("/api/hosts/:host/actions/:action/webtty", webhandler.MakeForwardToHost(webhandler.GetActionWebTTY))
	r.POST("/api/hosts/:host/actions/:action", webhandler.MakeForwardToHost(webhandler.GetAction))

	restapi.RunBoth(r, "priv.pem", "cert.pem")
}
