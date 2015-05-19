package sender

import (
	"github.com/open-falcon/transfer/g"
	cpool "github.com/open-falcon/transfer/sender/conn_pool"
	"github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	judgeInstances := set.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools = cpool.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())

	graphInstances := set.NewStringSet()
	for _, instance := range cfg.Graph.Cluster {
		graphInstances.Add(instance)
	}
	GraphConnPools = cpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())

	if cfg.Graph.Migrating && cfg.Graph.ClusterMigrating != nil {
		graphMigratingInstances := set.NewStringSet()
		for _, instance := range cfg.Graph.ClusterMigrating {
			graphMigratingInstances.Add(instance)
		}
		GraphMigratingConnPools = cpool.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
			cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphMigratingInstances.ToSlice())
	}
}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	GraphMigratingConnPools.Destroy()
}
