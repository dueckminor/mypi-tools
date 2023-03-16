package debug

import (
	"context"

	"github.com/gin-gonic/gin"
)

type componentDebug struct {
	component
	connector SSHConnector
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

func newServiceDebug(svcs *services, rgAPI *gin.RouterGroup) Service {
	svc := newEmptyService(svcs, "mypi-debug")

	comp := &componentDebug{}
	comp.startFunc = comp.startSSH
	comp.info.Name = "ssh"
	comp.info.Service = svc.Name()
	comp.info.Actions = []ActionInfo{
		{Name: "restart"},
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

	ccNodejs := newComponent(svc, "web")

	ccNodejs.Start()
	svc.GetComponent("web")

	return svc
}
