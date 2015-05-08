package sender

import (
	"github.com/open-falcon/transfer/g"
	"github.com/toolkits/container/list"
)

func initSendQueues() {
	cfg := g.Config()
	for node, _ := range cfg.Judge.Cluster {
		Q := list.NewSafeLinkedListLimited(DefaultSendQueueMaxSize)
		JudgeQueues[node] = Q
	}

	for node, _ := range cfg.Graph.Cluster {
		Q := list.NewSafeLinkedListLimited(DefaultSendQueueMaxSize)
		GraphQueues[node] = Q
	}

	if cfg.Graph.Migrating && cfg.Graph.ClusterMigrating != nil {
		for node, _ := range cfg.Graph.ClusterMigrating {
			Q := list.NewSafeLinkedListLimited(DefaultSendQueueMaxSize)
			GraphMigratingQueues[node] = Q
		}
	}
}
