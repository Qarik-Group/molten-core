package flannel

import (
	"context"
	"fmt"
	"time"

	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	EtcdSubnetsPath string = "/coreos.com/network/subnets"
)

func PersistSubnetReservations() error {
	etcdClient, err := util.NewEtcdV2Client()
	kapi, err := util.NewAPIClient(etcdClient)
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
