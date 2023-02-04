package debug

import (
	"context"
	"os"
	"os/exec"
	"path"

	"github.com/dueckminor/mypi-tools/go/gotty/buffered"
	"github.com/dueckminor/mypi-tools/go/util/network"
)

type Component interface {
	MessageHost
	GetData() any
	GetInfo() ComponentInfo
	Name() string
	GetTTY() (buffered.BufferedTty, error)
	Start() error
	Stop() error
	SetPort(port int) error
	SetState(state string)
}

type ComponentInfo struct {
	Service string `json:"service"`
	Name    string `json:"name"`
	Port    int    `json:"port"`
	State   string `json:"state"`
}

//##############################################################################

type component struct {
	messageHost
	svc     *service
	info    ComponentInfo
	ctx     context.Context
	cancel  context.CancelFunc
	tty     buffered.BufferedTty
	stopped chan bool
}

func newComponent(svc *service, name string) *component {
	comp := new(component)
	comp.svc = svc
	comp.info.Name = name
	comp.info.Service = svc.name
	svc.components[name] = comp

	comp.messageHost.Subscribe("*", func(topic string, value any) {
		svc.messageHost.Publish(name+"/"+topic, value)
	})

	return comp
}

func (comp *component) GetData() any {
	return comp.info
}

func (comp *component) GetInfo() ComponentInfo {
	return comp.info
}

func (comp *component) SetPort(port int) (err error) {
	if port == 0 {
		port = comp.info.Port
	}
	if port <= 0 {
		port, err = comp.NewPort()
		if err != nil {
			return err
		}
	}
	if port != comp.info.Port {
		comp.info.Port = port
		comp.messageHost.Publish("port", port)
	}
	return nil
}

func (comp *component) SetState(state string) {
	if len(state) > 0 && state != comp.info.State {
		comp.info.State = state
		comp.messageHost.Publish("state", state)
	}
}

func (comp *component) Name() string {
	return comp.info.Name
}

func (comp *component) GetTTY() (tty buffered.BufferedTty, err error) {
	if comp.tty == nil {
		tty, err = buffered.NewBufferedTty()
		if err != nil {
			return nil, err
		}
		comp.tty = tty
	}
	return comp.tty, nil
}

func (comp *component) NewCommand(name string, arg ...string) *exec.Cmd {
	if nil == comp.ctx || comp.ctx.Err() != nil {
		comp.ctx, comp.cancel = context.WithCancel(context.Background())
	}

	tty, err := comp.GetTTY()
	if err != nil {
		panic(err)
	}

	cmd := exec.CommandContext(comp.ctx, name, arg...)

	cmd.Env = append(os.Environ(),
		"PATH="+path.Join(GetWorkspaceRoot(), ".venv/bin")+":"+os.Getenv("PATH"),
		"PYTHONPATH="+path.Join(GetWorkspaceRoot(), "scripts", "lib"),
		"MYPI_DEBUG_URL=http://localhost:8080",
		"MYPI_DEBUG_SECRET="+comp.svc.svcs.authClient.LocalSecret,
	)

	tty.AttachProcess(cmd)
	return cmd
}

func (comp *component) NewPort() (port int, err error) {
	return network.GetFreePort()
}

func (comp *component) Stop() error {

	cancel := comp.cancel
	stopped := comp.stopped

	if cancel == nil {
		return nil
	}

	cancel()
	<-stopped

	return nil
}

func (comp *component) Start() error {
	comp.stopped = make(chan bool)

	dir := GetWorkspaceRoot()
	ctrl := path.Join(dir, "debug", "services", comp.info.Service, "components", comp.info.Name, "ctrl")

	cmd := comp.NewCommand(ctrl, "run")
	cmd.Dir = dir

	err := cmd.Start()
	if err != nil {
		comp.tty.ClosePTY()
		cmd = comp.NewCommand(ctrl, "run")
		cmd.Dir = dir
		err = cmd.Start()
	}

	if err != nil {
		close(comp.stopped)
		comp.stopped = nil
		comp.cancel = nil
		return err
	}

	go func() {
		err = cmd.Wait()
		comp.tty.ClosePTY()
		close(comp.stopped)
		comp.stopped = nil
		comp.cancel = nil
	}()

	return nil
}
