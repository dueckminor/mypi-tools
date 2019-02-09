package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

type Container struct {
	Name   string
	ID     string
	Status Status
}

func GetContainersRaw(ctx context.Context) (containers []types.Container, err error) {
	return cli.ContainerList(ctx, types.ContainerListOptions{})
}

func makeContainer(dockerContainer types.Container) *Container {
	status := StatusUnknown
	if dockerContainer.State == "running" {
		status = StatusRunning
	}

	return &Container{
		Name:   dockerContainer.Names[0],
		ID:     dockerContainer.ID,
		Status: status,
	}
}

func GetContainers(ctx context.Context) (containers []*Container, err error) {
	dockerContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	containers = make([]*Container, 0, len(dockerContainers))

	for _, dockerContainer := range dockerContainers {
		containers = append(containers, makeContainer(dockerContainer))
	}

	return containers, nil
}

func GetContainer(ctx context.Context, name string) (container *Container, err error) {
	args := filters.NewArgs()
	args.Add("name", "/"+name)
	dockerContainers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}

	return makeContainer(dockerContainers[0]), nil
}
