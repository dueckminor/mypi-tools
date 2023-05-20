package debug

import (
	"io/ioutil"
	"path"

	"golang.org/x/exp/maps"
)

type Service interface {
	MessageHost
	GetData() any
	Name() string
	AddComponent(comp Component)
	GetComponents() []Component
	GetComponent(name string) Component
}

//##############################################################################

type service struct {
	messageHost
	name       string
	svcs       *services
	components map[string]Component
}

type serviceData struct {
	Name       string          `json:"name"`
	Components []ComponentInfo `json:"components"`
}

func (sc *service) GetData() any {
	components := []ComponentInfo{}
	for _, component := range sc.components {
		components = append(components, component.GetInfo())
	}
	return serviceData{
		Name:       sc.name,
		Components: components,
	}
}

func (sc *service) Name() string {
	return sc.name
}
func (sc *service) AddComponent(comp Component) {
	if sc.components == nil {
		sc.components = make(map[string]Component)
	}

	name := comp.Name()
	sc.components[name] = comp

	comp.Subscribe("*", func(topic string, value any) {
		sc.messageHost.Publish(name+"/"+topic, value)
	})

}

func (sc *service) GetComponents() []Component {
	if nil == sc.components {
		return []Component{}
	}
	return maps.Values(sc.components)
}
func (sc *service) GetComponent(name string) Component {
	if nil == sc.components {
		return nil
	}
	return sc.components[name]
}

func newEmptyService(svcs *services, name string) *service {
	svc := new(service)
	svc.svcs = svcs
	svc.name = name
	svc.components = make(map[string]Component)
	svcs.AddService(svc)
	return svc
}

func newGenericService(svcs *services, name string) Service {
	svc := newEmptyService(svcs, name)

	servicesDir := path.Join(GetWorkspaceRoot(), "debug", "services", name, "components")

	files, err := ioutil.ReadDir(servicesDir)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			newComponent(svc, file.Name())
		}
	}

	return svc
}
