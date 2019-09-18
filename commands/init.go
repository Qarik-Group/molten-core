package commands

import (
	"fmt"
	"log"

	"github.com/coreos/go-systemd/dbus"
	"github.com/starkandwayne/molten-core/bucc"
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
	isBuccHost, err := bucc.IsBuccHost()
	if err != nil {
		return fmt.Errorf("failed to determine BUCC host: %s", err)
	}
	if isBuccHost {
		units.Update(units.BUCC)
	}
	units.Update(units.Docker)

	conn, err := dbus.New()
	if err != nil {
		return fmt.Errorf("failed to connect to systemd D-Bus: %s", err)
	}

	if conn.Reload() != nil {
		return fmt.Errorf("failed to reload systemd: %s", err)
	}

	// - make flannel subnet persistent (remove ttl)
	// - write docker certs to disk
	return nil
}
