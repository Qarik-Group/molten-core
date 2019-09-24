package bucc

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/starkandwayne/molten-core/config"
)

const (
	buccImage             = "rkoster/bucc"
	buccHostStateDir      = "/var/lib/moltencore/bucc"
	buccContainerStateDir = "/bucc/state"
)

func WriteStateDir(conf *config.NodeConfig) error {
	if err := os.MkdirAll(buccHostStateDir, 0775); err != nil {
		return fmt.Errorf("failed to create state dir: %s", err)
	}

	err := writeVars(filepath.Join(buccHostStateDir, "vars.yml"), conf)
	if err != nil {
		return fmt.Errorf("failed to write vars file: %s", err)
	}
	return nil
}

func CreateContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	reader, err := cli.ImagePull(ctx, buccImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: buccImage,
		Entrypoint: []string{
			"/bucc/bin/bucc",
			"up",
			"--recreate",
			"--cpi",
			"docker",
			"--flannel",
			"--unix-sock",
			"--concourse-lb",
		},
	}, &container.HostConfig{
		DNS: []string{"8.8.4.4", "8.8.8.8"},
		Mounts: []mount.Mount{
			{
				Type:        mount.TypeBind,
				Source:      buccHostStateDir,
				Target:      buccContainerStateDir,
				Consistency: mount.ConsistencyFull,
				ReadOnly:    false,
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
