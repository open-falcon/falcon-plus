package sender

import (
	"github.com/open-falcon/common/pool"
	"github.com/open-falcon/transfer/g"
	"github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	JudgeConnPools = pool.CreateSafeRpcConnPools(
		cfg.Judge.PingMethod,
		cfg.Judge.MaxConns,
		cfg.Judge.MaxIdle,
		cfg.Judge.ConnTimeout,
		cfg.Judge.CallTimeout,
	)
	judgeInstances := set.NewStringSet()
	for _, instance := range cfg.Judge.Cluster {
		judgeInstances.Add(instance)
	}
	JudgeConnPools.Init(judgeInstances.ToSlice())

	GraphConnPools = pool.CreateSafeRpcConnPools(
		cfg.Graph.PingMethod,
		cfg.Graph.MaxConns,
		cfg.Graph.MaxIdle,
		cfg.Graph.ConnTimeout,
		cfg.Graph.CallTimeout,
	)
	graphInstances := set.NewStringSet()
	for _, instance := range cfg.Graph.Cluster {
		graphInstances.Add(instance)
	}
	GraphConnPools.Init(graphInstances.ToSlice())

	if cfg.Graph.Migrating && cfg.Graph.ClusterMigrating != nil {
		GraphMigratingConnPools = pool.CreateSafeRpcConnPools(
			cfg.Graph.PingMethod,
			cfg.Graph.MaxConns,
			cfg.Graph.MaxIdle,
			cfg.Graph.ConnTimeout,
			cfg.Graph.CallTimeout,
		)
		graphMigratingInstances := set.NewStringSet()
		for _, instance := range cfg.Graph.ClusterMigrating {
			graphMigratingInstances.Add(instance)
		}
		GraphMigratingConnPools.Init(graphMigratingInstances.ToSlice())
	}
}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	GraphMigratingConnPools.Destroy()
}
