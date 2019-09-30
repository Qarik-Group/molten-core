package commands

import (
	"fmt"
	"log"

	"github.com/starkandwayne/molten-core/bucc"
	"github.com/starkandwayne/molten-core/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type UpdateBUCCConfigsCommand struct {
	logger *log.Logger
}

func (cmd *UpdateBUCCConfigsCommand) register(app *kingpin.Application) {
	app.Command("update-bucc-configs", "update configs in BOSH and Credhub ").Action(cmd.run)
}

func (cmd *UpdateBUCCConfigsCommand) run(c *kingpin.ParseContext) error {
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

	cmd.logger.Printf("Updating BOSH Cloud Config")
	if err = bc.UpdateCloudConfig(confs); err != nil {
		return fmt.Errorf("failed to update BOSH Cloud Config: %s", err)
	}

	cmd.logger.Printf("Updating BOSH CPI Config")
	if err = bc.UpdateCPIConfig(confs); err != nil {
		return fmt.Errorf("failed to update BOSH CPI Config: %s", err)
	}

	cmd.logger.Printf("Updating BOSH Runtime Config")
	if err = bc.UpdateRuntimeConfig(confs); err != nil {
		return fmt.Errorf("failed to update BOSH Runtime Config: %s", err)
	}

	cmd.logger.Printf("Updating Credhub MoltenCore Config (for consumption via Concourse)")
	if err = bc.UpdateMoltenCoreConfig(confs); err != nil {
		return fmt.Errorf("failed to update Credhub MoltenCore Config: %s", err)
	}

	return nil
}
