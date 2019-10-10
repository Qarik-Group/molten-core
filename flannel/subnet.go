package flannel

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	EtcdSubnetsPath string = "/coreos.com/network/subnets"
)

var (
	_, FlannelNetwork, _ = net.ParseCIDR("10.1.0.0/16")
)

type Subnet struct {
	cidr *net.IPNet
}

func GetSubnetByIndex(i uint16) (Subnet, error) {
	// increment index by 1 since first flannel subnet does not work
	s, err := cidr.Subnet(FlannelNetwork, 8, int(i)+1)
	if err != nil {
		return Subnet{}, fmt.Errorf("failed get flannel subnet by index: %d got: %s", i, err)
	}
	return Subnet{cidr: s}, nil
}

func ConfigureSubnet(s Subnet, publicIP net.IP) error {
	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return err
	}

	data := struct{ PublicIP net.IP }{PublicIP: publicIP}
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to generate flannel subnet config: %s", err)
	}

	ctx := context.Background()
	_, err = kapi.Set(ctx, s.etcdKey(), string(value), &client.SetOptions{
		TTL: 0 * time.Second})
	if err != nil {
		return fmt.Errorf("failed write flannel subnet config to etcd: %s", err)
	}

	return nil
}

func (s Subnet) Host(num int) (net.IP, error) {
	return cidr.Host(s.cidr, num)
}

func (s Subnet) String() string {
	return s.cidr.String()
}

func (s Subnet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.cidr.String())
}

func (s *Subnet) UnmarshalJSON(data []byte) error {
	var parsed string
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}
	_, cidr, err := net.ParseCIDR(parsed)
	s.cidr = cidr
	return err
}

func (s Subnet) etcdKey() string {
	return filepath.Join(EtcdSubnetsPath,
		strings.Replace(s.cidr.String(), "/", "-", -1))
}
