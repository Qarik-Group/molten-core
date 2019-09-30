package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type UpdateBoshConfigsCommand struct {
	logger *log.Logger
}

func (cmd *UpdateBoshConfigsCommand) register(app *kingpin.Application) {
	app.Command("update-bosh-configs", "update {cloud,cpi,runtime}-configs in BOSH").Action(cmd.run)
}

func (cmd *UpdateBoshConfigsCommand) run(c *kingpin.ParseContext) error {
	cmd.logger.Printf("Loading node config")
	conf, err := config.LoadNodeConfig()
	if err != nil {
		return fmt.Errorf("failed load node config: %s", err)
	}

	cmd.logger.Printf("Assigning Zones to node configs")
	if err := config.AssignZonesToNodeConfigs(); err != nil {
		return fmt.Errorf("failed assign zones to node configs: %s", err)
	}

	cmd.logger.Printf("Loading node configs")
	confs, err := config.LoadNodeConfigs()
	if err != nil {
		return fmt.Errorf("failed load node configs: %s", err)
	}

	bc, err := bucc.NewClient(cmd.logger, conf)
	if err != nil {
		return fmt.Errorf("failed create BUCC client: %s", err)
	}

	if err = bc.UpdateCloudConfig(confs); err != nil {
		return fmt.Errorf("failed to update BUCC Cloud Config: %s", err)
	}

	if err = bc.UpdateCPIConfig(confs); err != nil {
		return fmt.Errorf("failed to update BUCC CPI Config: %s", err)
	}

	if err = bc.UpdateRuntimeConfig(confs); err != nil {
		return fmt.Errorf("failed to update BUCC Runtime Config: %s", err)
	}

	return nil
}
