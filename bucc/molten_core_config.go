package bucc

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/starkandwayne/molten-core/config"
)

var (
	sizeMultipliers = []int{1, 2, 4, 8, 16, 32, 64}
)

type moltenCoreConfig struct {
	SingletonAZ string            `json:"singleton_az"`
	OtherAZs    []string          `json:"other_azs"`
	AllAZs      []string          `json:"all_azs"`
	PublicIPs   map[string]string `json:"public_ips"`
	Sizes       sizes             `json:"sizes"`
}

type sizes struct {
	OtherAZs map[string]int `json:"other_azs"`
	AllAZs   map[string]int `json:"all_azs"`
}

func RenderMoltenCoreConfig(confs *[]config.NodeConfig) (string, error) {
	var mcconf moltenCoreConfig
	mcconf.PublicIPs = make(map[string]string)
	for _, conf := range *confs {
		mcconf.AllAZs = append(mcconf.AllAZs, conf.Zone())
		mcconf.PublicIPs[conf.Zone()] = conf.PublicIP.String()

		if conf.IsSingletonZone() {
			mcconf.SingletonAZ = conf.Zone()
		} else {
			mcconf.OtherAZs = append(mcconf.OtherAZs, conf.Zone())
		}
	}

	if len(*confs) == 1 {
		mcconf.OtherAZs = mcconf.AllAZs
	}

	sort.Strings(mcconf.OtherAZs)
	sort.Strings(mcconf.AllAZs)

	mcconf.Sizes.AllAZs = make(map[string]int)
	mcconf.Sizes.OtherAZs = make(map[string]int)
	for _, size := range sizeMultipliers {
		mcconf.Sizes.AllAZs[fmt.Sprintf("x%d", size)] = size * len(mcconf.AllAZs)
		mcconf.Sizes.OtherAZs[fmt.Sprintf("x%d", size)] = size * len(mcconf.OtherAZs)
	}

	raw, err := json.Marshal(mcconf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MoltenCore Config: %s", err)
	}

	return string(raw), nil
}
