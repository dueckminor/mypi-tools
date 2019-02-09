package webhandler

import (
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/dueckminor/mypi-api/go/docker"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func (wh *WebHandler) getDynDNS(c *gin.Context) {
	c.Data(200, "text/plain", []byte("foo"))
}

func (wh *WebHandler) putDynDNS(c *gin.Context) {

}

func (wh *WebHandler) getContainers(c *gin.Context) {
	containers, err := wh.dockerCLI.ContainerList(c, types.ContainerListOptions{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, containers)
}

func (wh *WebHandler) postContainersStop(c *gin.Context) {
	id := c.Param("id")
	wh.dockerCLI.ContainerStop(c, id, nil)
}

func (wh *WebHandler) postServicesStart(c *gin.Context) {
	//	id := c.Param("id")
	var err error
	result := struct {
		ID string
	}{}

	result.ID, err = docker.StartAction(c, "mqtt", "start")
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}

type WebHandler struct {
	dockerCLI *client.Client
}

// SetupEnpoints registers the http endpoints
func (wh *WebHandler) SetupEndpoints(r *gin.Engine) (err error) {
	wh.dockerCLI, err = client.NewEnvClient()
	if err != nil {
		return err
	}

	r.GET("api/dyndns", wh.getDynDNS)
	r.PUT("api/dyndns", wh.putDynDNS)
	r.GET("api/containers", wh.getContainers)
	r.POST("api/containers/:id/stop", wh.postContainersStop)
	r.GET("api/services", wh.getServices)
	r.POST("api/services/:id/start", wh.postServicesStart)
	// r.GET("api/services:id/logs/ws", func(c *gin.Context) {
	// 	handler := websocket.Handler(EchoServer)
	// 	handler.ServeHTTP(c.Writer, c.Request)
	// })

	return nil
}
