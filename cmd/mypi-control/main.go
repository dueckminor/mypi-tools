package main

import (
	"flag"

	"github.com/dueckminor/mypi-tools/go/gotty/webtty"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/webhandler"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    webtty.Protocols,
}

func main() {
	flag.Parse()

	r := gin.Default()

	wh := &webhandler.WebHandler{}
	err := wh.SetupEndpoints(r)

	if err != nil {
		panic(err)
	}

	r.GET("/ws", func(c *gin.Context) {
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithError(400, err)
		}
		writer, err := conn.NextWriter(1)
		if err != nil {
			c.AbortWithError(400, err)
		}
		defer writer.Close()
	})

	restapi.Run(r)
}
