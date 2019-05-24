package services

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/dueckminor/mypi-tools/go/docker"

	"github.com/dueckminor/mypi-tools/go/config"
)

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

type Service struct {
	Name      string
	Container string
	Status    Status
}

var (
	services []*Service
)

func init() {
	config.GetRoot()
	files, err := ioutil.ReadDir(config.GetRoot() + "/services")
	if err != nil {
		log.Fatal(err)
	}

	services = make([]*Service, 0, len(files))

	for _, f := range files {
		services = append(services, &Service{
			Name:   f.Name(),
			Status: StatusUnknown,
		})
	}

}

func GetServices(ctx context.Context) (result []*Service, err error) {
	for _, service := range services {
		container, err := docker.GetContainer(ctx, service.Name)
		if err != nil {
			continue
		}
		service.Container = container.ID
	}
	return services, nil
}
