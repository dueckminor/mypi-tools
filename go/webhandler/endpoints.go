package webhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dueckminor/mypi-tools/go/users"

	"github.com/docker/docker/api/types"
	"github.com/dueckminor/mypi-tools/go/docker"

	// "github.com/docker/docker/api/types/container"
	// "github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func (wh *WebHandler) postLogin(c *gin.Context) {
	var params struct {
		Username string
		Password string
	}
	err := c.BindJSON(&params)

	cookie, _ := c.Cookie("token")
	fmt.Println("Cookie:", cookie)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !users.CheckPassword(params.Username, params.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var response struct {
		Username string
	}
	response.Username = params.Username

	c.SetCookie("token", params.Username, 3600, "/", "rpi.fritz.box", true, false)

	c.JSON(http.StatusOK, response)
}

type UserInfo struct {
	Name string `json:"text,omitempty"`
	Icon string `json:"icon,omitempty"`
}

func (wh *WebHandler) getUsers(c *gin.Context) {
	userCfg, err := users.ReadUserCfg()
	if err != nil {
		fmt.Println(err)
	}
	users, err := userCfg.GetUsers()

	userInfos := make([]UserInfo, len(users), len(users))

	for i, user := range users {
		userInfos[i].Name = user.Name
		userInfos[i].Icon = "mdi-account" // "mdi-account-star"
	}

	data, err := json.Marshal(userInfos)
	c.Data(200, "application/json", data)
}

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

// WebHandler handles all http endpoints
type WebHandler struct {
	dockerCLI *client.Client
}

// SetupEndpoints registers the http endpoints
func (wh *WebHandler) SetupEndpoints(r *gin.Engine) (err error) {
	wh.dockerCLI, err = client.NewEnvClient()
	if err != nil {
		return err
	}

	r.GET("api/users", wh.getUsers)
	r.POST("api/login", wh.postLogin)
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
