package debug

import (
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

func (comp *componentDebug) Start() error {
	tty, err := comp.GetTTY()
	if err != nil {
		return err
	}
	pty, err := tty.CreatePTY()
	if err != nil {
		return err
	}

	comp.connector, err = StartSSHConnector("ssh://pi@mypi:2022", 8443, pty)
	if err != nil {
		return err
	}
	return nil
}

func (service *serviceDebug) handler(c *gin.Context) {
	p := c.Request.URL.Path
	if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "api") {
		// we don't want to handle the API here
		return
	}

	handler := service.proxyHandler
	//if handler == nil {
	handler = service.fileHandler
	//}
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
	comp.info.Name = "ssh"
	comp.info.Service = svc.Name()

	svc.AddComponent(comp)

	err := comp.Start()
	if err != nil {
		panic(err)
	}

	svcs.Subscribe("mypi-router/golang/port", func(topic string, value any) {
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

	svc.Subscribe("web/port", func(topic string, value any) {
		svc.webPort = value.(int)
		if svc.webPort <= 0 {
			svc.proxyHandler = nil
		} else {
			svc.proxyHandler = ginutil.SingleHostReverseProxy(
				fmt.Sprintf("http://localhost:%d", svc.webPort))
		}
	})

	return svc
}
