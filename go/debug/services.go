package debug

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"time"

	"github.com/dueckminor/mypi-tools/go/auth"
	"github.com/dueckminor/mypi-tools/go/gotty/server/ginhandler"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type Services interface {
	MessageHost
	AddService(service Service)
	GetServices() []Service
	GetService(name string) Service
	GetComponent(svcName string, compName string) Component

	AddGenericService(name string)

	Run()
}

//##############################################################################

type services struct {
	messageHost
	services     map[string]Service
	r            *gin.Engine
	rgAPI        *gin.RouterGroup
	serviceDebug ServiceDebug
	authClient   auth.AuthClientLocalSecret
}

func NewServices(r *gin.Engine) Services {
	svcs := new(services)
	svcs.authClient = auth.AuthClientLocalSecret{}
	svcs.authClient.CreateLocalSecret()
	svcs.services = make(map[string]Service)
	r.Use(svcs.authClient.GetHandler())

	svcs.r = r
	svcs.rgAPI = r.Group("/api")
	svcs.registerGinAPIHandler(svcs.rgAPI)
	svcs.load()

	return svcs
}

func (svcs *services) load() error {
	servicesDir := path.Join(GetWorkspaceRoot(), "debug", "services")

	files, err := ioutil.ReadDir(servicesDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			serviceName := file.Name()
			if serviceName == "mypi-debug" {
				svcs.addDebug()
			} else {
				svcs.AddGenericService(serviceName)
			}
		}
	}
	return nil
}

func (svcs *services) AddService(service Service) {
	name := service.Name()
	svcs.services[name] = service

	service.Subscribe("*", func(topic string, value any) {
		svcs.messageHost.Publish(name+"/"+topic, value)
	})
}

func (svcs *services) GetServices() []Service {
	if nil == svcs.services {
		return []Service{}
	}
	return maps.Values(svcs.services)
}

func (svcs *services) GetService(name string) Service {
	if nil == svcs.services {
		return nil
	}
	return svcs.services[name]
}

func (svcs *services) GetComponent(svcName string, compName string) Component {
	service := svcs.GetService(svcName)
	if service == nil {
		return nil
	}
	return service.GetComponent(compName)
}

func (svcs *services) addDebug() {
	svcs.serviceDebug = newServiceDebug(svcs, svcs.rgAPI)
}

func (svcs *services) AddGenericService(name string) {
	newGenericService(svcs, name)
}

func (svcs *services) Run() {
	go func() {
		time.Sleep(time.Second * 2)
		for _, svc := range svcs.services {
			if svc.Name() != "mypi-debug" {
				for _, comp := range svc.GetComponents() {
					comp.Start()
				}
			}
		}
	}()
	svcs.serviceDebug.Run(svcs.r)
}

func (svcs *services) registerGinAPIHandler(r *gin.RouterGroup) {

	ws := NewWS()
	ws.Run(r)
	svcs.messageHost.Subscribe("*", ws.Publish)

	r.GET("/services", svcs.getServices)
	r.GET("/services/:service", svcs.getService)
	r.GET("/services/:service/components", svcs.getComponents)
	r.POST("/services/:service/components/:component/restart", svcs.postComponentRestart) // to be deleted
	r.POST("/services/:service/components/:component/actions/:action", svcs.postComponentAction)
	r.GET("/services/:service/components/:component", svcs.getComponent)
	r.PATCH("/services/:service/components/:component", svcs.patchComponent)
	r.GET("/services/:service/components/:component/tty", svcs.getComponentTty)
}

func (svcs *services) ginGetService(c *gin.Context) Service {
	svc := svcs.GetService(c.Param("service"))
	if svc == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil
	}
	return svc
}

func (svcs *services) ginGetComponent(c *gin.Context) Component {
	comp := svcs.GetComponent(c.Param("service"), c.Param("component"))
	if comp == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil
	}
	return comp
}

func (svcs *services) getServices(c *gin.Context) {
	result := make([]any, 0)

	result = append(result, svcs.services["mypi-debug"].GetData())
	result = append(result, svcs.services["mypi-router"].GetData())

	names := make([]string, 0)

	for _, svc := range svcs.services {
		name := svc.Name()
		if name != "mypi-debug" && name != "mypi-router" {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	for _, name := range names {
		result = append(result, svcs.services[name].GetData())
	}

	c.JSON(http.StatusOK, result)
}

func (svcs *services) getService(c *gin.Context) {
	svc := svcs.ginGetService(c)
	c.JSON(http.StatusOK, svc.GetData())
}

func (svcs *services) getComponents(c *gin.Context) {
	comps := svcs.ginGetService(c).GetComponents()
	result := make([]any, 0)
	for _, comp := range comps {
		result = append(result, comp.GetData())
	}
	c.JSON(http.StatusOK, result)
}

func (svcs *services) postComponentRestart(c *gin.Context) {
	component := svcs.ginGetComponent(c)
	if component == nil {
		return
	}
	component.Stop()
	component.Start()
}

func (svcs *services) postComponentAction(c *gin.Context) {
	component := svcs.ginGetComponent(c)
	if component == nil {
		return
	}
	action := c.Param("action")
	if action == "restart" {
		component.Stop()
		component.Start()
	}
}

func (svcs *services) getComponent(c *gin.Context) {
	component := svcs.ginGetComponent(c)
	if component == nil {
		return
	}
	c.JSON(http.StatusOK, component.GetInfo())
}

func (svcs *services) patchComponent(c *gin.Context) {
	component := svcs.ginGetComponent(c)
	if component == nil {
		return
	}
	var compInfo ComponentInfo
	err := c.BindJSON(&compInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	fmt.Printf("Patching component %s: %v\n", component.Name(), compInfo)
	if len(compInfo.State) > 0 {
		component.SetState(compInfo.State)
	}
	component.SetPort(compInfo.Port)
	fmt.Printf("%v\n", component.GetInfo())
	c.JSON(http.StatusOK, component.GetInfo())
}

func (svcs *services) getComponentTty(c *gin.Context) {
	component := svcs.ginGetComponent(c)
	if component == nil {
		return
	}
	tty, _ := component.GetTTY()
	ginhandler.Handler(c, tty.GetFactory())
}
