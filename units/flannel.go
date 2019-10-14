package units

import (
	"fmt"

	"github.com/coreos/go-systemd/unit"
	"github.com/starkandwayne/molten-core/config"
	"github.com/starkandwayne/molten-core/flannel"
)

const (
	confNetworkCMDTmpl = `/usr/bin/etcdctl set /coreos.com/network/config '{"Network": "%s", "Backend": {"Type": "vxlan"}}'`
	flannelOPTSTmpl    = `FLANNEL_OPTS="--iface=%s --public-ip=%s"`
)

func Flannel(conf *config.NodeConfig) Unit {
	return Unit{
		Name: "flanneld.service",
		DropIns: []DropIn{
			{
				Name: "30-mc-flannel.conf",
				Contents: []*unit.UnitOption{
					unit.NewUnitOption("Service", "ExecStartPre",
						fmt.Sprintf(confNetworkCMDTmpl, flannel.FlannelNetwork.String())),
					unit.NewUnitOption("Service", "Environment",
						fmt.Sprintf(flannelOPTSTmpl, conf.PrivateIP, conf.PrivateIP)),
				},
			},
		},
	}
}
