package bucc

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/starkandwayne/molten-core/config"
)

var (
	allSizes = []int{1, 2, 4, 8, 16, 32, 64}
)

type moltenCoreConfig struct {
	PublicIPs map[string]string `json:"public_ips"`
	Scaling   scaling           `json:"scaling"`
}

type scaling struct {
	Odd3 map[string]azsAndInstances `json:"odd3"`
	Odd5 map[string]azsAndInstances `json:"odd5"`
	Max1 map[string]azsAndInstances `json:"max1"`
	Max2 map[string]azsAndInstances `json:"max2"`
	Max3 map[string]azsAndInstances `json:"max3"`
	All  map[string]azsAndInstances `json:"all"`
}

type azsAndInstances struct {
	Azs       []string `json:"azs"`
	Instances int      `json:"instances"`
}

func RenderMoltenCoreConfig(confs *[]config.NodeConfig) (string, error) {
	var mcconf moltenCoreConfig

	azs := make([]string, 0)
	mcconf.PublicIPs = make(map[string]string)
	for _, conf := range *confs {
		mcconf.PublicIPs[conf.Zone()] = conf.PublicIP.String()
		azs = append(azs, conf.Zone())
	}
	sort.Strings(azs)

	mcconf.Scaling.All = make(map[string]azsAndInstances)
	for _, size := range allSizes {
		mcconf.Scaling.All[fmt.Sprintf("x%d",
			size)] = azsAndInstances{azs, size * len(azs)}
	}

	mcconf.Scaling.Max1 = slice(azs, 1, 1)

	switch size := len(*confs); {
	case size == 1:
		mcconf.Scaling.Max2 = slice(azs, size, 5)
		mcconf.Scaling.Max3 = slice(azs, size, 3)
		mcconf.Scaling.Odd3 = slice(azs, 1, 3)
		mcconf.Scaling.Odd5 = slice(azs, 1, 2)
	case size == 2:
		mcconf.Scaling.Max2 = slice(azs, 2, 5)
		mcconf.Scaling.Max3 = slice(azs, size, 3)
		mcconf.Scaling.Odd3 = slice(azs, 1, 3)
		mcconf.Scaling.Odd5 = slice(azs, 1, 2)
	case size < 5:
		mcconf.Scaling.Max2 = slice(azs, 2, 5)
		mcconf.Scaling.Max3 = slice(azs, 3, 3)
		mcconf.Scaling.Odd3 = slice(azs, 3, 3)
		mcconf.Scaling.Odd5 = slice(azs, 3, 2)
	case size >= 5:
		mcconf.Scaling.Max2 = slice(azs, 2, 5)
		mcconf.Scaling.Max3 = slice(azs, 3, 3)
		mcconf.Scaling.Odd3 = slice(azs, 3, 3)
		mcconf.Scaling.Odd5 = slice(azs, 5, 2)
	}

	raw, err := json.Marshal(mcconf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MoltenCore Config: %s", err)
	}

	return string(raw), nil
}

func slice(azs []string, size int, slices int) map[string]azsAndInstances {
	out := make(map[string]azsAndInstances)

	for i := 0; i < slices; i++ {
		if size > len(azs) {
			size = len(azs)
		}
		j := i * size
		k := (i + 1) * size
		if k > len(azs) {
			j = 0
			k = size
		}

		out[fmt.Sprintf("slice%d", i+1)] = azsAndInstances{azs[j:k], size}
	}

	return out
}
