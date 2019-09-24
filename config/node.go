package config

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/starkandwayne/molten-core/certs"
	"github.com/starkandwayne/molten-core/flannel"
	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	EtcdMolenCorePath  string = "/moltencore/nodes"
	DockerCertValidFor        = time.Hour * 24 * 365
	DockerTLSPort             = 2376
)

type Docker struct {
	Endpoint string
	CA       certs.Cert
	Server   certs.Cert
	Client   certs.Cert
}

type NodeConfig struct {
	Subnet flannel.Subnet
	Docker Docker
}

func LoadNodeConfig() (*NodeConfig, error) {
	hostIP, err := util.LookupIpV4Address(false)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup node address: %s", err)
	}

	subnets, err := flannel.GetSubnets(&hostIP)
	if err != nil || len(subnets) == 0 {
		return nil, fmt.Errorf("failed to get flannel subnet: %s", err)
	}
	subnet := subnets[0]

	kapi, err := util.NewEtcdV2Client()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, nodePath(subnet), nil)
	if err != nil && !client.IsKeyNotFound(err) {
		return nil, fmt.Errorf("failed to load node config from etcd: %s", err)
	}

	if resp != nil {
		var c NodeConfig
		err = json.Unmarshal([]byte(resp.Node.Value), &c)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal node config: %s", err)
		}
		return &c, nil
	}

	docker, err := newDocker(subnet, hostIP)
	if err != nil {
		return nil, fmt.Errorf("failed to generate docker certs: %s", err)
	}

	conf := NodeConfig{Subnet: subnet, Docker: docker}

	rawConf, err := json.Marshal(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal node config: %s", err)
	}
	_, err = kapi.Set(ctx, nodePath(subnet), string(rawConf), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update node config in etcd: %s", err)
	}

	return &conf, nil
}

func nodePath(s flannel.Subnet) string {
	key, _ := s.IP(1)
	return filepath.Join(EtcdMolenCorePath, key.String())
}

func newDocker(s flannel.Subnet, hostIP net.IP) (Docker, error) {
	caCert, err := certs.Genereate(certs.GenArg{
		ValidFor: DockerCertValidFor,
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker ca cert: %s", err)
	}

	serverCert, err := certs.Genereate(certs.GenArg{
		CA:          caCert,
		ValidFor:    DockerCertValidFor,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{hostIP, net.ParseIP("127.0.0.1")},
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker server cert: %s", err)
	}

	clientCert, err := certs.Genereate(certs.GenArg{
		CA:          caCert,
		ValidFor:    DockerCertValidFor,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker client cert: %s", err)
	}

	return Docker{
		Endpoint: fmt.Sprintf("%s:%d", hostIP, DockerTLSPort),
		CA:       caCert,
		Server:   serverCert,
		Client:   clientCert,
	}, nil
}
