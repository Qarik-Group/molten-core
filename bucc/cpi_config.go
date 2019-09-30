package bucc

import (
	"encoding/json"
	"fmt"
	"github.com/starkandwayne/molten-core/config"
)

type cpi struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Properties dockerProperties
}

type dockerProperties struct {
	Host    string    `json:"host"`
	Gateway string    `json:"gateway"`
	TLS     dockerTLS `json:"tls"`
}

type dockerTLS struct {
	CA          []byte `json:"ca"`
	Certificate []byte `json:"certificate"`
	PrivateKey  []byte `json:"private_key"`
}

func renderCPIConfig(confs *[]config.NodeConfig) (string, error) {
	var cpis []cpi

	for _, conf := range *confs {
		endpoint := fmt.Sprintf("tcp://%s", conf.Docker.Endpoint)

		cpis = append(cpis, cpi{
			Name: conf.CPIName(),
			Type: "docker",
			Properties: dockerProperties{
				Host: endpoint,
				TLS: dockerTLS{
					CA:          conf.Docker.CA.Cert,
					Certificate: conf.Docker.Client.Cert,
					PrivateKey:  conf.Docker.Client.Key,
				},
			},
		})
	}

	cpisRaw, err := json.Marshal(cpis)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cpis: %s", err)
	}

	raw := fmt.Sprintf(`{"cpis":%s}`, cpisRaw)
	return raw, nil
}
