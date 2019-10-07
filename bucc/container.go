package bucc

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/starkandwayne/molten-core/config"
)

const (
	buccImage             = "starkandwayne/mc-bucc:latest"
	buccHostStateDir      = "/var/lib/moltencore/bucc"
	buccContainerStateDir = "/bucc/state"
	credhubMoltenCorePath = "/concourse/main/moltencore"
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
	if err := c.writeStateDir(); err != nil {
		return err
	}

	return c.run([]string{
		"/bucc/bin/bucc",
		"up",
		"--recreate",
		"--cpi",
		"docker",
		"--flannel",
		"--unix-sock",
		"--concourse-lb",
	}, false)
}

func (c *Client) Shell() error {
	return c.run([]string{"/bin/bash", "-c",
		"/bin/bash --init-file <(echo 'source ~/.bashrc && bucc fly >/dev/null')"}, true)
}

func (c *Client) UpdateCloudConfig(confs *[]config.NodeConfig) error {
	data, err := renderCloudConfig(confs)
	if err != nil {
		return fmt.Errorf("failed to render Cloud Config: %s", err)
	}
	return c.updateBoshConfig("cloud", data)
}

func (c *Client) UpdateCPIConfig(confs *[]config.NodeConfig) error {
	data, err := renderCPIConfig(confs)
	if err != nil {
		return fmt.Errorf("failed to render CPI Config: %s", err)
	}
	return c.updateBoshConfig("cpi", data)
}

func (c *Client) UpdateRuntimeConfig(confs *[]config.NodeConfig) error {
	return c.updateBoshConfig("runtime", renderRuntimeConfig())
}

func (c *Client) UpdateMoltenCoreConfig(confs *[]config.NodeConfig) error {
	data, err := renderMoltenCoreConfig(confs)
	if err != nil {
		return fmt.Errorf("failed to render MoltenCore Config: %s", err)
	}
	return c.credHubSet(credhubMoltenCorePath, data)
}

func (c *Client) updateBoshConfig(t, config string) error {
	cmd := fmt.Sprintf("source <(/bucc/bin/bucc env) && bosh -n update-%s-config <(echo '%s')", t, config)
	return c.run([]string{"/bin/bash", "-c", cmd}, false)
}

func (c *Client) credHubSet(path, config string) error {
	cmd := fmt.Sprintf("source <(/bucc/bin/bucc env) && credhub set -n %s -t json -v '%s'", path, config)
	return c.run([]string{"/bin/bash", "-c", cmd}, false)
}

func (c *Client) pullImage() error {
	ctx := context.Background()
	images, err := c.dcli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list images: %s", err)
	}
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == buccImage {
				return nil
			}
		}
	}

	c.logger.Printf("Pulling MoltenCore docker image, this can take a while")
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

func (c *Client) run(entrypoint []string, tty bool) error {
	if err := c.pullImage(); err != nil {
		return err
	}

	// TODO limit ip range to reserved range
	networks := make(map[string]*network.EndpointSettings)
	networks[config.BOSHDockerNetworkName] = &network.EndpointSettings{}

	ctx := context.Background()
	resp, err := c.dcli.ContainerCreate(ctx, &container.Config{
		AttachStdin:  tty,
		AttachStdout: tty,
		AttachStderr: tty,
		Tty:          tty,
		OpenStdin:    tty,
		Image:        buccImage,
		Entrypoint:   entrypoint,
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

	if tty {
		// statusCh, errCh := c.dcli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

		// resp, err := c.dcli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		// 	Stdin:  true,
		// 	Stdout: true,
		// 	//			Stderr: true,
		// 	//			Stream: true,
		// })
		// if err != nil {
		// 	return fmt.Errorf("failed attach to docker container: %s", err)
		// }

		// defer resp.Close()

		// select {
		// case err := <-errCh:
		// 	if err != nil {
		// 		return fmt.Errorf("failed start docker container: %s", err)
		// 	}
		// case status := <-statusCh:
		// 	if status.StatusCode != 0 {
		// 		return fmt.Errorf("container process failed: %s", status.Error.Message)
		// 	}
		// }

		// stdcopy.StdCopy(os.Stdout, os.Stderr, session.Reader)

		cmd := exec.Command("docker", "attach", resp.ID)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
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
				return fmt.Errorf("container process failed: %s", status.Error)
			}
		}
	}

	return nil
}
