package flannel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/subosito/gotenv"

	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	EtcdSubnetsPath      string = "/coreos.com/network/subnets"
	flannelSubnetFile           = "/run/flannel/subnet.env"
	flannelSubnetEnvName        = "FLANNEL_SUBNET"
)

type Subnet struct {
	cidr *net.IPNet
}

func RemoveSubnetTTL(s Subnet) error {
	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, s.etcdKey(), nil)
	if err != nil {
		return fmt.Errorf("failed to get flannel subnet: %s", err)
	}

	_, err = kapi.Set(ctx, resp.Node.Key, resp.Node.Value, &client.SetOptions{
		TTL: 0 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to update subnet TTL for %s got: %s", resp.Node.Key, err)
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

func (s Subnet) CIDR() string {
	return s.cidr.String()
}

func (s Subnet) Equals(subnet Subnet) bool {
	return s.cidr.String() == subnet.cidr.String()
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

func GetSubnetStatus() (bool, string, error) {
	members, err := getEtcdMemberCount()
	if err != nil {
		return false, "", err
	}
	subnets, err := getSubnets()
	if err != nil {
		return false, "", err
	}
	return len(subnets) == members, fmt.Sprintf("%d/%d",
		members, len(subnets)), nil
}

func IsFirstSubnet(subnet Subnet) (bool, error) {
	subnets, err := getSubnets()
	if err != nil {
		return false, err
	}
	return subnet.Equals(subnets[0]), nil
}

func GetNodeSubnet() (Subnet, error) {
	gotenv.Load(flannelSubnetFile)
	v := os.Getenv(flannelSubnetEnvName)
	_, ipv4Net, err := net.ParseCIDR(v)
	if err != nil {
		return Subnet{}, fmt.Errorf("failed to parse subnet CIDR: %s", err)
	}
	return Subnet{cidr: ipv4Net}, nil
}

func getEtcdMemberCount() (int, error) {
	c, err := util.NewEtcdV2MembersAPI()
	if err != nil {
		return 0, err
	}

	members, err := c.List(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to list etcd members: %s", err)
	}

	return len(members), nil
}

func getSubnets() ([]Subnet, error) {
	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return []Subnet{}, err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, EtcdSubnetsPath, &client.GetOptions{Recursive: true})
	if err != nil {
		return []Subnet{}, fmt.Errorf("failed to list flannel subnets: %s", err)
	}

	subnets := make(map[string]Subnet)
	var ips []net.IP
	for _, n := range resp.Node.Nodes {
		res := struct {
			PublicIP net.IP
		}{}
		err = json.Unmarshal([]byte(n.Value), &res)
		if err != nil {
			return []Subnet{}, fmt.Errorf("failed to unmarshal flannel subnet reservation: %s", err)
		}
		subnet, err := parseFlannelSubnet(n.Key)
		if err != nil {
			return []Subnet{}, fmt.Errorf("failed to parse flannel subnet: %s", err)
		}

		subnets[res.PublicIP.String()] = subnet
		ips = append(ips, res.PublicIP)
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	out := []Subnet{}
	for _, ip := range ips {
		out = append(out, subnets[ip.String()])
	}
	return out, nil
}

func parseFlannelSubnet(s string) (Subnet, error) {
	s = strings.Replace(filepath.Base(s), "-", "/", -1)
	_, ipv4Net, err := net.ParseCIDR(s)
	if err != nil {
		return Subnet{}, fmt.Errorf("failed to parse subnet CIDR: %s", err)
	}
	return Subnet{cidr: ipv4Net}, nil
}
