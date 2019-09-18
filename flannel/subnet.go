package flannel

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

// etcdctl -o extended get /flannel/network/subnets/10.1.7.0-24
func PersistSubnetReservations() error {
	kapi, err := util.NewEtcdV2Client()
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, "/flannel/network/subnets", nil)
	if err != nil {
		return fmt.Errorf("failed to list flannel subnets: %s", err)
	}

	for _, n := range resp.Node.Nodes {
		_, err := kapi.Set(ctx, n.Key, n.Value, &client.SetOptions{
			TTL: 0 * time.Second})
	}
	return nil
}
