package bucc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/starkandwayne/molten-core/config"
	"net"
	"sort"
	"strconv"
)

type moltenCoreConfig struct {
	SingletonAZ string            `json:"az_singleton"`
	OtherAZs    []string          `json:"azs_other"`
	AllAZs      []string          `json:"azs"`
	PublicIPs   map[string]string `json:"public_ips"`
}

func renderMoltenCoreConfig(confs *[]config.NodeConfig) (string, error) {
	var mcconf moltenCoreConfig
	mcconf.PublicIPs = make(map[string]string)
	var publicIPs []net.IP

	for _, conf := range *confs {
		mcconf.AllAZs = append(mcconf.AllAZs, conf.Zone())

		if conf.IsSingletonZone() {
			mcconf.SingletonAZ = conf.Zone()
		} else {
			mcconf.OtherAZs = append(mcconf.OtherAZs, conf.Zone())
			publicIPs = append(publicIPs, conf.PublicIP)
		}
	}

	sort.Slice(publicIPs, func(i, j int) bool {
		return bytes.Compare(publicIPs[i], publicIPs[j]) < 0
	})

	for i, ip := range publicIPs {
		mcconf.PublicIPs[strconv.FormatInt(int64(i), 10)] = ip.String()
	}

	raw, err := json.Marshal(mcconf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MoltenCore Config: %s", err)
	}

	return string(raw), nil
}
