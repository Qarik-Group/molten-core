package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/config"
	"github.com/starkandwayne/molten-core/flannel"
	"github.com/starkandwayne/molten-core/units"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type InitCommand struct {
	logger *log.Logger
}

func (cmd *InitCommand) register(app *kingpin.Application) {
	app.Command("init", "bootstrap node into MoltenCore cluster member").Action(cmd.run)
}

func (cmd *InitCommand) run(c *kingpin.ParseContext) error {
	cmd.logger.Printf("Loading node config")
	conf, err := config.LoadNodeConfig()
	if err != nil {
		return fmt.Errorf("failed load node config: %s", err)
	}

	cmd.logger.Printf("Writing Docker TLS certs")
	err = units.WriteDockerTLSCerts(conf.Docker)
	if err != nil {
		return fmt.Errorf("failed to write docker certs: %s", err)
	}

	cmd.logger.Printf("Writing MoltenCore managed systemd unit files")
	u := []units.Unit{
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

	cmd.logger.Printf("Removing Flannel subnet TTL")
	if err = flannel.RemoveSubnetTTL(conf.Subnet); err != nil {
		return fmt.Errorf("failed to persist flannel subnets: %s", err)
	}

	return nil
}
