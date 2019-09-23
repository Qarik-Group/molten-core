package flannel

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	EtcdSubnetsPath string = "/coreos.com/network/subnets"
)

type Subnet struct {
	cidr *net.IPNet
}

func LookupNodeSubnet(nodeIP net.IP) (Subnet, error) {
	kapi, err := util.NewEtcdV2Client()
	if err != nil {
		return Subnet{}, err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, EtcdSubnetsPath, &client.GetOptions{Recursive: true})
	if err != nil {
		return Subnet{}, fmt.Errorf("failed to list flannel subnets: %s", err)
	}

	for _, n := range resp.Node.Nodes {
		reservation := struct {
			PublicIP net.IP
		}{}
		err = json.Unmarshal([]byte(n.Value), &reservation)
		if err != nil {
			return Subnet{}, fmt.Errorf("failed to unmarshal flannel subnet reservation: %s", err)
		}
		if reservation.PublicIP.Equal(nodeIP) {
			return parseFlannelSubnet(n.Key)
		}
	}
	return Subnet{}, fmt.Errorf("failed to find subnet for: %s", nodeIP.String())
}

func PersistSubnetReservations() error {
	kapi, err := util.NewEtcdV2Client()
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, EtcdSubnetsPath, nil)
	if err != nil {
		return fmt.Errorf("failed to list flannel subnets: %s", err)
	}

	for _, n := range resp.Node.Nodes {
		_, err := kapi.Set(ctx, n.Key, n.Value, &client.SetOptions{
			TTL: 0 * time.Second})
		if err != nil {
			return fmt.Errorf("failed to update subnet TTL for %s got: %s", n.Key, err)
		}
	}
	return nil
}

func (s Subnet) IP(i uint8) (net.IP, error) {
	ip := net.ParseIP(strings.Replace(s.cidr.String(), "0/24", strconv.Itoa(int(i)), -1))
	if !s.cidr.Contains(ip) {
		return net.IP{}, fmt.Errorf("ip: %s, out of range: %s", ip, s.cidr)
	}
	return ip, nil
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

func parseFlannelSubnet(s string) (Subnet, error) {
	s = strings.Replace(filepath.Base(s), "-", "/", -1)
	_, ipv4Net, err := net.ParseCIDR(s)
	if err != nil {
		return Subnet{}, fmt.Errorf("failed to parse subnet CIDR: %s", err)
	}
	return Subnet{cidr: ipv4Net}, nil
}
