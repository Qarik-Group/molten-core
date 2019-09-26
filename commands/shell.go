package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ShellCommand struct {
	logger *log.Logger
}

func (cmd *ShellCommand) register(app *kingpin.Application) {
	app.Command("shell", "start interactive shell for interacting with BUCC").Action(cmd.run)
}

func (cmd *ShellCommand) run(c *kingpin.ParseContext) error {
	conf, err := config.LoadNodeConfig()
	if err != nil {
		return fmt.Errorf("failed load node config: %s", err)
	}

	bc, err := bucc.NewClient(cmd.logger, conf)
	if err != nil {
		return fmt.Errorf("failed create BUCC client: %s", err)
	}

	if err = bc.Shell(); err != nil {
		return fmt.Errorf("failed to start shell container: %s", err)
	}
	return nil
}
