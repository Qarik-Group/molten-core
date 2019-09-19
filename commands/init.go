package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/bucc"
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
	cmd.logger.Printf("Writing MoltenCore managed systemd unit files")
	u := []units.Unit{
		units.DockerTLSSocket,
		units.Docker,
	}

	isBuccHost, err := bucc.IsBuccHost()
	if err != nil {
		return fmt.Errorf("failed to determine BUCC host: %s", err)
	}
	if isBuccHost {
		u = append(u, units.BUCC)
	}

	err = units.Enable(u)
	if err != nil {
		return fmt.Errorf("failed enable systemd units: %s", err)
	}

	cmd.logger.Printf("Persisting Flannel subnet reservations")
	if err = flannel.PersistSubnetReservations(); err != nil {
		return fmt.Errorf("failed to persist flannel subnets: %s", err)
	}

	// - write docker certs to disk
	return nil
}
