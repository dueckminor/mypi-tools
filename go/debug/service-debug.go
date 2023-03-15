package debug

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type ServiceDebug interface {
	Service
	Run(r *gin.Engine)
}

type componentDebug struct {
	component
	connector SSHConnector
}

type serviceDebug struct {
	service
	mutex      sync.Mutex
	webPort    int
	distFolder string

	fileHandler  gin.HandlerFunc
	proxyHandler gin.HandlerFunc
}

func (comp *componentDebug) startSSH(ctx context.Context) (result chan error, err error) {
	tty, err := comp.GetTTY()
	if err != nil {
		return nil, err
	}
	pty, err := tty.CreatePTY()
	if err != nil {
		return nil, err
	}

	if nil == comp.connector {
		comp.connector = new(sshConnector)
	}

	result, err = comp.connector.Run(ctx, "ssh://pi@mypi:2022", 8443, pty)
	if err != nil {
		return nil, err
	}
	comp.SetState("running")
	return result, nil
}

func (service *serviceDebug) handler(c *gin.Context) {
	p := c.Request.URL.Path
	if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "api") {
		// we don't want to handle the API here
		return
	}

	handler := service.proxyHandler
	if handler == nil {
		handler = service.fileHandler
	}
	handler(c)

	//if !c.IsAborted() {
	//	c.File(path.Join(service.distFolder, "index.html"))
	//}
}

func (service *serviceDebug) Run(r *gin.Engine) {
	restapi.LocalhostOnly()
	r.Use(service.handler)
	panic(r.Run("localhost:8080"))
}

func newServiceDebug(svcs *services, rgAPI *gin.RouterGroup) ServiceDebug {
	svc := &serviceDebug{}
	svc.name = "mypi-debug"
	svc.svcs = svcs

	svcs.AddService(svc)

	comp := &componentDebug{}
	comp.startFunc = comp.startSSH
	comp.info.Name = "ssh"
	comp.info.Service = svc.Name()
	comp.info.Actions = []ActionInfo{
		ActionInfo{
			Name: "restart",
		},
	}

	svc.AddComponent(comp)

	err := comp.Start()
	if err != nil {
		comp.SetState("failed")
	}

	svcs.Subscribe("mypi-router/go/port", func(topic string, value any) {
		comp.connector.SetLocalRouterPort(value.(int))
	})

	sshFS, err := comp.connector.GetHttpFS()
	if err == nil {
		rgAPI.StaticFS("fs/", sshFS)
	}

	ccNodejs := newComponent(&svc.service, "web")

	workspace := GetWorkspaceRoot()
	svc.distFolder = path.Join(workspace, "web", "mypi-debug", "dist")
	svc.fileHandler = static.ServeRoot("/", svc.distFolder)

	// if util.FileExists(path.Join(service.distFolder, "index.html")) {
	// } else {
	// }
	ccNodejs.Start()

	compWeb := svc.GetComponent("web")

	svc.Subscribe("web/state", func(topic string, value any) {
		if compWeb.GetInfo().State == "running" {
			svc.webPort = compWeb.GetInfo().Port
			svc.proxyHandler = ginutil.SingleHostReverseProxy(
				fmt.Sprintf("http://localhost:%d", svc.webPort))
		} else {
			svc.proxyHandler = nil
		}
	})

	return svc
}
