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
    },
    {
      "disk_size": 1024,
      "name": "1GB"
    },
    {
      "disk_size": 5120,
      "name": "5GB"
    },
    {
      "disk_size": 10240,
      "name": "10GB"
    },
    {
      "disk_size": 100240,
      "name": "100GB"
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
    },
    {
      "name": "minimal",
      "cloud_properties": {
	"RestartPolicy": {
	  "Name": "always"
	}
      }
    },
    {
      "name": "small",
      "cloud_properties": {
	"RestartPolicy": {
	  "Name": "always"
	}
      }
    },
    {
      "name": "small-highmem",
      "cloud_properties": {
	"RestartPolicy": {
	  "Name": "always"
	}
      }
    }
  ],
  "vm_extensions": [
    {
      "name": "5GB_ephemeral_disk"
    },
    {
      "name": "10GB_ephemeral_disk"
    },
    {
      "name": "50GB_ephemeral_disk"
    },
    {
      "name": "100GB_ephemeral_disk"
    },
    {
      "name": "500GB_ephemeral_disk"
    },
    {
      "name": "1TB_ephemeral_disk"
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
