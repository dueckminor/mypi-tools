package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/tlsconfig"
	"github.com/dueckminor/mypi-tools/go/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	// authURI   string
	port     = flag.Int("port", 8080, "The port")
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
	// targetURI string

	hostname string
	autoOpen = false
	nearHome = true
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		err = config.InitApp(*mypiRoot)
		if err != nil {
			panic(err)
		}
	}
	config.GetConfig()
}

type OwntracksMsg struct {
	InRegions []string `json:"inregions"`
}

func mqttOwnTracks(client mqtt.Client, msg mqtt.Message) {
	ownTracksMsg := OwntracksMsg{}
	err := json.Unmarshal(msg.Payload(), &ownTracksMsg)
	if err != nil {
		return
	}
	if util.StringsContains(ownTracksMsg.InRegions, "Work") {
		autoOpen = true
		client.Publish("tor/autoopen", 2, true, `true`)
	}
	if util.StringsContains(ownTracksMsg.InRegions, "NearHome") {
		if autoOpen {
			autoOpen = false
			fmt.Println("Open gate")
			client.Publish("tor/open", 2, false, `1`)
		}
		client.Publish("tor/autoopen", 2, true, `false`)
		nearHome = true
	} else {
		nearHome = false
	}
}

func mqttAutoOpen(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	shouldAutoOpen := (payload == "true") || (payload == "1")
	if nearHome && shouldAutoOpen {
		fmt.Println("reject autoopen request")
		client.Publish("tor/autoopen", 2, true, `false`)
		return
	}
	autoOpen = shouldAutoOpen
}

func handleQR(c *gin.Context) {
	c.Header("Content-Type", "image/png")
	png, _ := qrcode.Encode("todo", qrcode.Medium, 256)
	_, err := c.Writer.Write(png)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err) // nolint: errcheck
	}
}

func main() {

	r := gin.Default()

	tlsconfig := tlsconfig.NewTLSConfig()
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://mqtt-int:8883")
	opts.SetClientID(hostname).SetTLSConfig(tlsconfig)

	// Start the connection
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient.Subscribe("tor/autoopen", 2, mqttAutoOpen)
	mqttClient.Subscribe("owntracks/#", 2, mqttOwnTracks)

	r.Use(static.ServeRoot("/config", "/opt/owntracks/config"))
	r.GET("/qr", handleQR)

	panic(r.Run(":" + strconv.Itoa(*port)))
}
