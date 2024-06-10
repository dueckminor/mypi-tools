package main

import (
	"bytes"
	"flag"
	"fmt"
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

	_ "github.com/dueckminor/mypi-tools/go/restapi/http"
	//"github.com/dueckminor/mypi-tools/go/restapi/https"
)

type led struct {
	file     string
	inverted bool
}

func (l led) exists() bool {
	return util.FileExists(l.file)
}

func (l led) set(state bool) {
	if l.inverted {
		state = !state
	}
	var err error
	if !state {
		err = os.WriteFile(l.file, []byte("0"), os.ModePerm)
	} else {
		err = os.WriteFile(l.file, []byte("1"), os.ModePerm)
	}
	if err != nil {
		fmt.Println("Failed to flash led:", err)
	}
}

func flashLED() {
	led0 := led{file: "/sys/class/leds/led0/brightness"}
	led1 := led{file: "/sys/class/leds/led1/brightness"}

	if !led0.exists() || !led1.exists() {
		led0 = led{file: "/sys/class/leds/ACT/brightness", inverted: true}
		led1 = led{file: "/sys/class/leds/PWR/brightness"}
		if !led0.exists() || !led1.exists() {
			return
		}
	}
	state := false
	for {
		led0.set(state)
		state = !state
		led1.set(state)
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
			if err != nil {
				panic("failed to start setup")
			}
			c.Stdin = bytes.NewReader([]byte("{}"))
			go func() {
				c.Start() // nolint: errcheck
				c.Wait()  // nolint: errcheck
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
	go server.Serve() // nolint: errcheck
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

	//https.SetKeyFiles("priv.pem", "cert.pem")
	restapi.Run(r)
}
