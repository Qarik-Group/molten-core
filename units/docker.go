package units

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/unit"

	"github.com/starkandwayne/molten-core/config"
)

const (
	dockerSSLDir          = "/var/ssl/docker"
	BoshDockerNetworkName = "bosh"
)

var (
	Docker Unit = Unit{
		Name: "docker.service",
		DropIns: []DropIn{
			{
				Name: "60-reset-flannel-default-bridge.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "ExecStartPre",
						// 10.255.240.0/20 last network /20 from 10.0.0.0/8 range
						// bridge only support 1000 hosts so no need for big network
						// docker default 172.17.0.0/16 conflicts with virtual box (coreos-vagrant)
						"/bin/sh -c 'echo \"DOCKER_OPT_BIP=\\\\\"--bip=10.255.240.0/20\\\\\"\" > /run/flannel/flannel_docker_opts.env'"),
					unit.NewUnitOption("Service", "ExecStartPre",
						"/bin/sh -c 'echo \"DOCKER_OPT_IPMASQ=\\\\\"--ip-masq=true\\\\\"\" >> /run/flannel/flannel_docker_opts.env'"),
					unit.NewUnitOption("Service", "ExecStartPre",
						"/bin/sh -c 'echo \"DOCKER_OPT_MTU=\\\\\"--mtu=1500\\\\\"\" >> /run/flannel/flannel_docker_opts.env'"),
				},
			},
			{
				Name: "30-enable-mtls.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "Environment",
						fmt.Sprintf("DOCKER_OPTS=\"--tlsverify --tlscacert=%s/ca.pem --tlscert=%s/cert.pem --tlskey=%s/key.pem\"",
							dockerSSLDir, dockerSSLDir, dockerSSLDir)),
				},
			},
			{
				Name: "70-create-bosh-network.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "EnvironmentFile", "/run/flannel/subnet.env"),
					unit.NewUnitOption("Service", "ExecStartPost",
						fmt.Sprintf("/bin/sh -c 'docker network create -d bridge --subnet=${FLANNEL_SUBNET} --attachable --opt com.docker.network.driver.mtu=${FLANNEL_MTU} %s || true'",
							BoshDockerNetworkName)),
				},
			},
		},
	}
)

func DockerTLSSocket(conf config.Docker) Unit {
	return Unit{
		Name: "docker.socket",
		DropIns: []DropIn{
			{
				Name: "30-listen-stream.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Socket", "ListenStream", conf.Endpoint),
					unit.NewUnitOption("Socket", "FreeBind", "true"),
				},
			},
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
