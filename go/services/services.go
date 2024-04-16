package services

import (
	"context"
	"log"
	"os"
	"sync"

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
	initOnce sync.Once
	services []*Service
)

func initServices() {
	initOnce.Do(func() {
		config.GetRoot()
		files, err := os.ReadDir(config.GetRoot() + "/services")
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
	})
}

func GetServices(ctx context.Context) (result []*Service, err error) {
	initServices()
	for _, service := range services {
		container, err := docker.GetContainer(ctx, service.Name)
		if err != nil {
			continue
		}
		service.Container = container.ID
	}
	return services, nil
}
