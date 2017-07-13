package rings

import (
	consistent "github.com/toolkits/consistent"
)

// 一致性哈希环,用于管理服务器节点.
type ConsistentHashNodeRing struct {
	ring *consistent.Consistent
}

func NewConsistentHashNodesRing(numberOfReplicas int32, nodes []string) *ConsistentHashNodeRing {
	ret := &ConsistentHashNodeRing{ring: consistent.New()}
	ret.SetNumberOfReplicas(numberOfReplicas)
	ret.SetNodes(nodes)
	return ret
}

// 根据pk,获取node节点. chash(pk) -> node
func (this *ConsistentHashNodeRing) GetNode(pk string) (string, error) {
	return this.ring.Get(pk)
}

func (this *ConsistentHashNodeRing) SetNodes(nodes []string) {
	for _, node := range nodes {
		this.ring.Add(node)
	}
}

func (this *ConsistentHashNodeRing) SetNumberOfReplicas(num int32) {
	this.ring.NumberOfReplicas = int(num)
}
