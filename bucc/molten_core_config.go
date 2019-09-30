package bucc

import (
	"encoding/json"
	"fmt"
	"github.com/starkandwayne/molten-core/config"
)

type moltenCoreConfig struct {
	SingletonAZ string   `json:"az_singleton"`
	OtherAZs    []string `json:"azs_other"`
	AllAZs      []string `json:"azs"`
	PublicIPs   []string `json:"public_ips"`
}

func renderMoltenCoreConfig(confs *[]config.NodeConfig) (string, error) {
	var mcconf moltenCoreConfig

	for _, conf := range *confs {
		mcconf.AllAZs = append(mcconf.AllAZs, conf.Zone())

		if conf.IsSingletonZone() {
			mcconf.SingletonAZ = conf.Zone()
		} else {
			mcconf.OtherAZs = append(mcconf.OtherAZs, conf.Zone())
			mcconf.PublicIPs = append(mcconf.PublicIPs, conf.PublicIP.String())
		}
	}

	raw, err := json.Marshal(mcconf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MoltenCore Config: %s", err)
	}

	return string(raw), nil
}
