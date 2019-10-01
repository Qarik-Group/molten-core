package bucc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/starkandwayne/molten-core/config"
)

const (
	dockerSocket = "/run/docker.sock"
)

type Vars struct {
	DirectorName    string `json:"director_name"`
	Alias           string `json:"alias"`
	DockerHost      string `json:"docker_host"`
	Network         string `json:"network"`
	InternalCIDR    string `json:"internal_cidr"`
	InternalGW      string `json:"internal_gw"`
	InternalIP      string `json:"internal_ip"`
	HostIP          string `json:"host_ip"`
	ConcourseDomain string `json:"concourse_domain"`
}

func writeVars(path string, c *config.NodeConfig) error {
	gw, err := c.Subnet.IP(1)
	if err != nil {
		return fmt.Errorf("failed to get gatway ip: %s", err)
	}
	buccIP, err := c.Subnet.IP(10)
	if err != nil {
		return fmt.Errorf("failed to get bucc ip: %s", err)
	}

	vars := Vars{
		DirectorName:    "molten-core",
		Alias:           "mc",
		DockerHost:      dockerSocket,
		Network:         config.BOSHDockerNetworkName,
		InternalCIDR:    c.Subnet.CIDR(),
		InternalGW:      gw.String(),
		InternalIP:      buccIP.String(),
		HostIP:          "0.0.0.0",
		ConcourseDomain: c.PublicIP.String(),
	}

	data, err := json.Marshal(vars)
	if err != nil {
		return fmt.Errorf("failed to marshal vars file: %s", err)
	}

	return ioutil.WriteFile(path, data, 0644)
}
