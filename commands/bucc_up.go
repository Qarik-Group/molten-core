package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BuccUpCommand struct {
	logger *log.Logger
}

func (cmd *BuccUpCommand) register(app *kingpin.Application) {
	app.Command("bucc-up", "create or update BUCC").Action(cmd.run)
}

func (cmd *BuccUpCommand) run(c *kingpin.ParseContext) error {
	cmd.logger.Printf("Loading node config")
	conf, err := config.LoadNodeConfig()
	if err != nil {
		return fmt.Errorf("failed load node config: %s", err)
	}

	if err = bucc.WriteStateDir(conf); err != nil {
		return fmt.Errorf("failed to write state dir: %s", err)
	}
	bucc.CreateContainer()
	// docker create container starkandwayne/mc-bucc
	// load bucc creds from etcds
	// bucc int to generate creds
	// save creds and vars to etcd
	// bucc up

	return nil
}
