package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/starkandwayne/molten-core/config"
	"github.com/starkandwayne/molten-core/flannel"
	"github.com/starkandwayne/molten-core/units"
	"github.com/starkandwayne/molten-core/util"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type InitCommand struct {
	logger        *log.Logger
	flannelSubnet string
	zoneIndex     uint16
	dev           bool
}

func (cmd *InitCommand) register(app *kingpin.Application) {
	c := app.Command("init", "bootstrap node into MoltenCore cluster member").Action(cmd.run)
	c.Flag("zone", "Index of this node, used for BOSH availability zone").Required().Uint16Var(&cmd.zoneIndex)
	c.Flag("dev", "Base zone index of last private IP octet").BoolVar(&cmd.dev)
}

func (cmd *InitCommand) run(c *kingpin.ParseContext) error {
	cmd.logger.Printf("Generating node config")
	if cmd.dev {
		ip, _ := util.LookupIpV4Address(false)
		lastIPDiget := ip.String()[len(ip.String())-1:]
		i, _ := strconv.ParseInt(lastIPDiget, 10, 16)
		cmd.zoneIndex = uint16(i - 1)
	}
	conf, err := config.GenereateNodeConfig(cmd.zoneIndex)
	if err != nil {
		return fmt.Errorf("failed generate node config: %s", err)
	}

	cmd.logger.Printf("Writing Docker TLS certs")
	err = units.WriteDockerTLSCerts(conf.Docker)
	if err != nil {
		return fmt.Errorf("failed to write docker certs: %s", err)
	}

	cmd.logger.Printf("Writing MoltenCore managed systemd unit files")
	u := []units.Unit{
		units.Flannel(conf),
		units.DockerTLSSocket(conf.Docker),
		units.Docker,
	}
	if conf.IsSingletonZone() {
		u = append(u, units.BUCC...)
	}

	err = units.Enable(u)
	if err != nil {
		return fmt.Errorf("failed enable systemd units: %s", err)
	}

	cmd.logger.Printf("Configure Flannel subnet")
	if err = flannel.ConfigureSubnet(conf.Subnet, conf.PrivateIP); err != nil {
		return fmt.Errorf("failed to configure flannel subnet: %s", err)
	}

	return nil
}
