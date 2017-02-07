package sender

import (
	backend "github.com/open-falcon/falcon-plus/common/backend_pool"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	nset "github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	// judge
	judgeInstances := nset.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = backend.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())

	// tsdb
	if cfg.Tsdb.Enabled {
		TsdbConnPoolHelper = backend.NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
	}

	// graph
	graphInstances := nset.NewSafeSet()
	for _, nitem := range cfg.Graph.ClusterList {
		for _, addr := range nitem.Addrs {
			graphInstances.Add(addr)
		}
	}
	GraphConnPools = backend.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())

}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	TsdbConnPoolHelper.Destroy()
}
