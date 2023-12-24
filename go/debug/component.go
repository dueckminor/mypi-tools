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
	Debug() error
	SetPort(port int) error
	SetDist(dist string)
	SetState(state string)
}

type ActionInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Selected bool   `json:"selected"`
	Disabled bool   `json:"disabled"`
}

type ComponentInfo struct {
	Service string       `json:"service"`
	Name    string       `json:"name"`
	Port    int          `json:"port"`
	Dist    string       `json:"dist,omitempty"`
	State   string       `json:"state"`
	Actions []ActionInfo `json:"actions"`
}

type startFunc func(ctx context.Context) (result chan error, err error)

//##############################################################################

type component struct {
	messageHost
	svc       *service
	info      ComponentInfo
	ctx       context.Context
	cancel    context.CancelFunc
	tty       buffered.BufferedTty
	stopped   chan bool
	startFunc startFunc
}

func newComponent(svc *service, name string) *component {
	comp := new(component)
	comp.svc = svc
	comp.info.Name = name
	comp.info.Service = svc.name
	comp.info.State = "stopped"
	svc.components[name] = comp

	if name == "go" {
		comp.info.Actions = []ActionInfo{
			{
				Name:     "restart",
				Type:     "button",
				Disabled: false,
			},
			{
				Name:     "debug",
				Type:     "button",
				Disabled: false,
			},
		}
	} else {
		comp.info.Actions = []ActionInfo{
			{
				Name:     "restart",
				Type:     "button",
				Disabled: false,
			},
		}
	}

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

func (comp *component) SetDist(dist string) {
	if len(dist) > 0 && dist != comp.info.Dist {
		comp.info.Dist = dist
		comp.messageHost.Publish("dist", dist)
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

func (comp *component) createContext() context.Context {
	if nil == comp.ctx || comp.ctx.Err() != nil {
		comp.ctx, comp.cancel = context.WithCancel(context.Background())
	}
	return comp.ctx
}

func (comp *component) NewCommand(ctx context.Context, name string, arg ...string) *exec.Cmd {
	tty, err := comp.GetTTY()
	if err != nil {
		panic(err)
	}

	cmd := exec.CommandContext(ctx, name, arg...)

	cmd.Env = append(os.Environ(),
		"PATH="+path.Join(GetWorkspaceRoot(), ".venv/bin")+":"+os.Getenv("PATH"),
		"PYTHONPATH="+path.Join(GetWorkspaceRoot(), "scripts", "lib"),
		"MYPI_DEBUG_URL=http://localhost:8080",
		"MYPI_DEBUG_SECRET="+comp.svc.svcs.authClient.LocalSecret,
	)

	err = tty.AttachProcess(cmd)
	if err != nil {
		panic(err)
	}
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
	comp.SetState("stopped")

	return nil
}

func (comp *component) Start() error {
	f := comp.startFunc
	if nil == f {
		f = comp.startGeneric
	}
	return comp.start(f)
}

func (comp *component) start(f startFunc) error {
	ctx := comp.createContext()

	comp.stopped = make(chan bool)

	comp.SetState("starting")
	done, err := f(ctx)

	if err != nil {
		close(comp.stopped)
		comp.stopped = nil
		comp.cancel = nil
		return err
	}

	go func() {
		<-done
		comp.tty.ClosePTY()
		comp.SetState("stopped")
		close(comp.stopped)
		comp.stopped = nil
		comp.cancel = nil
	}()

	return nil
}

func (comp *component) Debug() error {
	return comp.start(func(ctx context.Context) (result chan error, err error) {
		return comp.startCtrl(ctx, "debug")
	})
}

func (comp *component) startGeneric(ctx context.Context) (result chan error, err error) {
	return comp.startCtrl(ctx, "run")
}

func (comp *component) startCtrl(ctx context.Context, args ...string) (result chan error, err error) {
	dir := GetWorkspaceRoot()
	ctrl := path.Join(dir, "debug", "services", comp.info.Service, "components", comp.info.Name, "ctrl")

	cmd := comp.NewCommand(ctx, ctrl, args...)
	cmd.Dir = dir

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	result = make(chan error)

	go func() {
		result <- cmd.Wait()
	}()

	return result, nil
}
