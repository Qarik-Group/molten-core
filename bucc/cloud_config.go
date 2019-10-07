package bucc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/starkandwayne/molten-core/config"
)

const (
	numberOfReservedIPs = 19
	ccTmpl              = `
{
  "azs": %s,
  "networks": [
    {
      "name": "default",
      "subnets": %s,
      "type": "manual"
    }
  ],
  "compilation": {
    "az": "z0",
    "network": "default",
    "reuse_compilation_vms": true,
    "vm_type": "default",
    "workers": 5
  },
  "disk_types": [
    {
      "disk_size": 1024,
      "name": "default"
    }
  ],
  "vm_types": [
    {
      "name": "default",
      "cloud_properties": {
        "RestartPolicy": {
          "Name": "always"
        }
      }
    }
  ]
}
`
)

type az struct {
	Name string `json:"name"`
	CPI  string `json:"cpi"`
}

type subnet struct {
	AZ              string            `json:"az"`
	Range           string            `json:"range"`
	Gateway         string            `json:"gateway"`
	Reserved        []string          `json:"reserved"`
	CloudProperties map[string]string `json:"cloud_properties"`
}

func renderCloudConfig(confs *[]config.NodeConfig) (string, error) {
	var azs []az
	var subnets []subnet

	for _, conf := range *confs {
		azs = append(azs, az{Name: conf.Zone(), CPI: conf.CPIName()})

		gw, err := conf.Subnet.IP(1)
		if err != nil {
			return "", fmt.Errorf("failed to determine cloud config gatway: %s", err)
		}
		resMax, err := conf.Subnet.IP(numberOfReservedIPs + 1)
		if err != nil {
			return "", fmt.Errorf("failed to determine cloud config reserved range: %s", err)
		}

		reserved := fmt.Sprintf("%s-%s", gw, resMax)

		subnets = append(subnets, subnet{
			AZ:       conf.Zone(),
			Range:    conf.Subnet.CIDR(),
			Gateway:  gw.String(),
			Reserved: []string{reserved},
			CloudProperties: map[string]string{
				"name": config.BOSHDockerNetworkName,
			},
		})
	}

	azsRaw, err := json.Marshal(azs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal azs: %s", err)
	}
	subnetsRaw, err := json.Marshal(subnets)
	if err != nil {
		return "", fmt.Errorf("failed to marshal subnets: %s", err)
	}

	raw := fmt.Sprintf(ccTmpl, azsRaw, subnetsRaw)
	raw = strings.ReplaceAll(raw, "\n", "")
	return raw, nil
}
