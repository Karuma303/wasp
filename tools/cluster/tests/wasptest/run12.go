// +build ignore

package wasptest

import (
	"fmt"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	waspapi "github.com/iotaledger/wasp/packages/apilib"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/tools/cluster"
	"time"
)

func Activate1SC(clu *cluster.Cluster, sc *cluster.Chain) error {
	if err := activate(sc, clu); err != nil {
		return fmt.Errorf("activate %s: %v\n", sc.Address, err)
	}
	fmt.Printf("[cluster] wait for %v to connect the committee\n", 3*time.Second)
	time.Sleep(3 * time.Second)
	return nil
}

func Deactivate1SC(clu *cluster.Cluster, sc *cluster.Chain) error {
	if err := deactivate(sc, clu); err != nil {
		return fmt.Errorf("deactivate %s: %v\n", sc.Address, err)
	}
	return nil
}

func ActivateAllSC(clu *cluster.Cluster) error {
	for _, sc := range clu.SmartContractConfig {
		if err := activate(&sc, clu); err != nil {
			return fmt.Errorf("activate %s: %v\n", sc.Address, err)
		}
	}

	return nil
}

func DeactivateAllSC(clu *cluster.Cluster) error {
	for _, sc := range clu.SmartContractConfig {
		if err := deactivate(&sc, clu); err != nil {
			return fmt.Errorf("deactivate %s: %v\n", sc.Address, err)
		}
	}

	return nil
}

func activate(sc *cluster.Chain, clu *cluster.Cluster) error {
	allNodesApi := clu.WaspHosts(sc.AllNodes(), (*cluster.WaspNodeConfig).ApiHost)

	addr, err := address.FromBase58(sc.Address)
	if err != nil {
		return err
	}
	return waspapi.ActivateChain(waspapi.ActivateChainParams{
		ChainID:  (coretypes.ChainID)(addr),
		ApiHosts: allNodesApi,
	})
}

func deactivate(sc *cluster.Chain, clu *cluster.Cluster) error {
	allNodesApi := clu.WaspHosts(sc.AllNodes(), (*cluster.WaspNodeConfig).ApiHost)

	addr, err := address.FromBase58(sc.Address)
	if err != nil {
		return err
	}
	return waspapi.DeactivateChain(waspapi.ActivateChainParams{
		ChainID:  (coretypes.ChainID)(addr),
		ApiHosts: allNodesApi,
	})
}
