package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func StartAction(ctx context.Context, service, action string, args ...string) (actionID string, err error) {
	cmd := []string{"/opt/mypi/services/" + service + "/actions/" + action}
	cmd = append(cmd, args...)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "dueckminor/aarch64-tools-docker",
		Cmd:   cmd,
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/opt/mypi",
				Target: "/opt/mypi",
			},
			mount.Mount{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}, nil, "")
	if err != nil {
		return "", err
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return resp.ID, nil
}
