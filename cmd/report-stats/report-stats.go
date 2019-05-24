package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func getOpenConnections(filterCBs ...func(local, remote string) (use, drop bool)) (openConnections uint64) {
	out, err := exec.Command("netstat", "-an").Output()

	if err != nil {
		return
	}

	bytesReader := bytes.NewReader(out)
	reader := bufio.NewReader(bytesReader)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fields := strings.Fields(line)
		if len(fields) != 6 {
			continue
		}
		if fields[0] == "tcp" && fields[5] == "ESTABLISHED" {
			use := len(filterCBs) == 0
			drop := false
			for _, filterCB := range filterCBs {
				useThis, dropThis := filterCB(fields[3], fields[4])
				use = use || useThis
				drop = drop || dropThis
			}
			if use || !drop {
				openConnections++
			}
		}
	}
	return
}

type MemInfo struct {
	MemTotal  uint64
	MemFree   uint64
	SwapTotal uint64
	SwapFree  uint64
}

type Report struct {
	MemInfo
	CPU  uint64
	Conn uint64
}

func getMemInfo() (memInfo MemInfo) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		if fields[0] == "MemTotal:" {
			memInfo.MemTotal, _ = strconv.ParseUint(fields[1], 10, 64)
			memInfo.MemTotal *= 1024
		}
		if fields[0] == "MemFree:" {
			memInfo.MemFree, _ = strconv.ParseUint(fields[1], 10, 64)
			memInfo.MemFree *= 1024
		}
		if fields[0] == "SwapTotal:" {
			memInfo.SwapTotal, _ = strconv.ParseUint(fields[1], 10, 64)
			memInfo.SwapTotal *= 1024
		}
		if fields[0] == "SwapFree:" {
			memInfo.SwapFree, _ = strconv.ParseUint(fields[1], 10, 64)
			memInfo.SwapFree *= 1024
		}
	}
	return
}

var (
	hostname string
	prefix   string
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}
	prefix = "/computer/" + hostname + "/"
}

func main() {

	tlsconfig := tlsconfig.NewTLSConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://rpi:8883")
	opts.SetClientID(hostname).SetTLSConfig(tlsconfig)

	ignoreConnections, _ := net.LookupHost("rpi")
	for i, hostname := range ignoreConnections {
		ignoreConnections[i] = hostname + ":8883"
	}
	fmt.Println(ignoreConnections)

	// Start the connection
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	idle0, total0 := getCPUSample()
	for range time.Tick(time.Duration(2 * time.Second)) {
		idle1, total1 := getCPUSample()

		idleTicks := idle1 - idle0
		totalTicks := total1 - total0

		idle0, total0 = idle1, total1

		report := Report{
			MemInfo: getMemInfo(),
			CPU:     1000 * (totalTicks - idleTicks) / totalTicks,
			Conn: getOpenConnections(func(local, remote string) (use, drop bool) {
				for _, ignoreConnection := range ignoreConnections {
					if remote == ignoreConnection {
						return false, true
					}
				}
				return false, false
			}),
		}

		reportJSON, _ := json.Marshal(report)

		c.Publish(prefix+"stat", 0, false, string(reportJSON))
	}
	c.Disconnect(250)
}
