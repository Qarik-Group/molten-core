package units

import (
	"github.com/coreos/go-systemd/unit"
)

var (
	DockerTLSSocket Unit = Unit{
		Name:   "docker-tls-tcp.socket",
		Enable: true,
		Contents: []*unit.UnitOption{
			unit.NewUnitOption("Unit", "Description", "Docker Secured Socket for the API"),
			unit.NewUnitOption("Socket", "ListenStream", "2376"),
			unit.NewUnitOption("Socket", "BindIPv6Only", "both"),
			unit.NewUnitOption("Socket", "Service", "docker.service"),
			unit.NewUnitOption("Install", "WantedBy", "sockets.target"),
		},
	}

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
		},
	}
)
