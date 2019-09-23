package units

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"

	"github.com/starkandwayne/molten-core/config"
)

const (
	dockerSSLDir = "/var/ssl/docker"
)

var (
	Docker Unit = Unit{
		Name: "docker.service",
		DropIns: []DropIn{
			{
				Name: "60-disable-flannel-default-bridge.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "ExecStartPre",
						"/bin/sh -c 'echo \"\" > /run/flannel/flannel_docker_opts.env'"),
				},
			},
			{
				Name: "20-clct-docker.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "Environment",
						fmt.Sprintf("DOCKER_OPTS=--tlsverify --tlscacert=%s/ca.pem --tlscert=%s/cert.pem --tlskey=%s/key.pem",
							dockerSSLDir, dockerSSLDir, dockerSSLDir)),
				},
			},
		},
	}
)

func DockerTLSSocket(conf config.Docker) Unit {
	return Unit{
		Name:   "docker-tls-tcp.socket",
		Enable: true,
		Contents: []*unit.UnitOption{
			unit.NewUnitOption("Unit", "Description", "Docker Secured Socket for the API"),
			unit.NewUnitOption("Socket", "ListenStream", conf.Endpoint),
			unit.NewUnitOption("Socket", "FreeBind", "true"),
			unit.NewUnitOption("Socket", "Service", "docker.service"),
			unit.NewUnitOption("Install", "WantedBy", "sockets.target"),
		},
	}
}

func WriteDockerTLSCerts(d config.Docker) error {
	if err := os.MkdirAll(dockerSSLDir, 0777); err != nil {
		return err
	}
	if err := writeFile("ca.pem", d.CA.Cert); err != nil {
		return err
	}
	if err := writeFile("cert.pem", d.Server.Cert); err != nil {
		return err
	}
	if err := writeFile("key.pem", d.Server.Key); err != nil {
		return err
	}
	return nil
}

func writeFile(name string, data []byte) error {
	path := filepath.Join(dockerSSLDir, name)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
