package bucc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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

type Client struct {
	logger *log.Logger
	config *config.NodeConfig
	dcli   *client.Client
}

func NewClient(l *log.Logger, conf *config.NodeConfig) (*Client, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %s", err)
	}
	cli.NegotiateAPIVersion(ctx)

	return &Client{logger: l, config: conf, dcli: cli}, nil
}

func (c *Client) Up() error {
	return c.run([]string{
		"/bucc/bin/bucc",
		"up",
		"--recreate",
		"--cpi",
		"docker",
		"--flannel",
		"--unix-sock",
		"--concourse-lb",
	})
}

func (c *Client) pullImage() error {
	ctx := context.Background()
	c.logger.Printf("Pulling %s docker image, this can take a while", buccImage)
	reader, err := c.dcli.ImagePull(ctx, buccImage, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %s", buccImage, err)
	}

	defer reader.Close()
	io.Copy(ioutil.Discard, reader)
	return nil
}

func (c *Client) writeStateDir() error {
	if err := os.MkdirAll(buccHostStateDir, 0775); err != nil {
		return fmt.Errorf("failed to create state dir: %s", err)
	}

	err := writeVars(filepath.Join(buccHostStateDir, "vars.yml"), c.config)
	if err != nil {
		return fmt.Errorf("failed to write vars file: %s", err)
	}
	return nil
}

func (c *Client) run(cmd []string) error {
	if err := c.pullImage(); err != nil {
		return err
	}

	if err := c.writeStateDir(); err != nil {
		return err
	}

	//	initContainerIp, err := conf.Subnet.IP(2)
	networks := make(map[string]*network.EndpointSettings)
	networks[units.BoshDockerNetworkName] = &network.EndpointSettings{
		// IPAMConfig: &network.EndpointIPAMConfig{
		// 	IPv4Address: initContainerIp.String(),
		// },
	}

	ctx := context.Background()
	resp, err := c.dcli.ContainerCreate(ctx, &container.Config{
		Image:      buccImage,
		Entrypoint: cmd,
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

	if err := c.dcli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed start docker container: %s", err)
	}

	statusCh, errCh := c.dcli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		out, err := c.dcli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true, ShowStderr: true, Follow: true})
		if err != nil {
			c.logger.Printf("[warning] Failed to tail docker container logs: %s", err)
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
