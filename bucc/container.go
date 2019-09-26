package bucc

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/starkandwayne/molten-core/config"
	"github.com/starkandwayne/molten-core/units"
)

const (
	buccImage             = "starkandwayne/mc-bucc"
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

func Up(l *log.Logger, conf *config.NodeConfig) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	l.Printf("Pulling %s docker image, this can take a while", buccImage)
	_, err = cli.ImagePull(ctx, buccImage, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed pull image %s: %s", buccImage, err)
	}
	// TODO find a way to display image pull progress without spamming systemd journal
	// io.Copy(os.Stdout, reader)

	initContainerIp, err := conf.Subnet.IP(2)
	networks := make(map[string]*network.EndpointSettings)
	networks[units.BoshDockerNetworkName] = &network.EndpointSettings{
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address: initContainerIp.String(),
		},
	}

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
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:        mount.TypeBind,
				Source:      buccHostStateDir,
				Target:      buccContainerStateDir,
				Consistency: mount.ConsistencyFull,
				ReadOnly:    false,
			},
			{
				Type:        mount.TypeBind,
				Source:      dockerSocket,
				Target:      dockerSocket,
				Consistency: mount.ConsistencyFull,
				ReadOnly:    false,
			},
		},
	}, &network.NetworkingConfig{
		EndpointsConfig: networks,
	}, "")
	if err != nil {
		return fmt.Errorf("failed create docker container: %s", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed start docker container: %s", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true, ShowStderr: true, Follow: true})
		if err != nil {
			l.Printf("[warning] Failed to tail docker container logs: %s", err)
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			return fmt.Errorf("failed start docker container: %s", err)
		}
	case status := <-statusCh:
		cancel()
		if status.StatusCode != 0 {
			return fmt.Errorf("container process failed: %s", status.Error.Message)
		}
	}
	return nil
}
