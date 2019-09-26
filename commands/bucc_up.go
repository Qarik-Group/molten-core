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

	if err = bucc.Up(cmd.logger, conf); err != nil {
		return fmt.Errorf("failed to create BUCC: %s", err)
	}

	// TODO store creds and state in etcd currently just in file
	// load bucc creds from etcds
	// bucc int to generate creds
	// save creds and vars to etcd
	return nil
}
