package bucc

import (
	"encoding/json"
	"fmt"
	"github.com/starkandwayne/molten-core/config"
)

type cpi struct {
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Properties dockerProperties `json:"properties"`
}

type dockerProperties struct {
	Host       string    `json:"host"`
	APIVersion string    `json:"api_version"`
	TLS        dockerTLS `json:"tls"`
}

type dockerTLS struct {
	CA          string `json:"ca"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

func renderCPIConfig(confs *[]config.NodeConfig) (string, error) {
	var cpis []cpi

	for _, conf := range *confs {
		endpoint := fmt.Sprintf("tcp://%s", conf.Docker.Endpoint)

		cpis = append(cpis, cpi{
			Name: conf.CPIName(),
			Type: "docker",
			Properties: dockerProperties{
				Host:       endpoint,
				APIVersion: "1.38",
				TLS: dockerTLS{
					CA:          string(conf.Docker.CA.Cert),
					Certificate: string(conf.Docker.Client.Cert),
					PrivateKey:  string(conf.Docker.Client.Key),
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
