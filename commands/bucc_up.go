package commands

import (
	"fmt"
	"log"

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
	_, err := config.LoadNodeConfig()
	if err != nil {
		return fmt.Errorf("failed load node config: %s", err)
	}

	return nil
}
