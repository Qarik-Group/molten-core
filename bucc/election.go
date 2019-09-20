package bucc

import (
	"context"
	"fmt"

	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

func IsBuccHost() (bool, error) {
	// - checks etcd for bucc node
	// -- starts election if not found

	// Leader gets current leader of the cluster
	kapi, err := util.NewEtcdV2Client()
	if err != nil {
		return false, err
	}
	mi := client.NewMembersAPI(kapi)
	ctx := context.Background()
	leader, err := mi.Leader(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve leader: %s", err)
	}

	fmt.Printf("%v", leader)
	//returns `&{d7fcb18db6259469 30a3b49403f2480c9fbe2be1f68b643f [http://172.17.8.101:2380] [http://172.17.8.101:2379]}`
	// https: //godoc.org/go.etcd.io/etcd/client#Member
	//return am i leader?

	return false, nil
}
