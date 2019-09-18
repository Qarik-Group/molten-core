package config

type Docker struct {
	Endpoint   string
	CACert     string
	CAKey      string
	ServerCert string
	ServerKey  string
	ClientCert string
	ClientKey  string
}

type NodeConfig struct {
	*Docker
	Subnet string
}

func LoadNodeConfig() (*NodeConfig, error) {
	// - generate node config (docker certs, host, subnet) and store in etcd
	return &NodeConfig{}, nil
}
