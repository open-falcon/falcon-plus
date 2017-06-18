package cache

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"sync"
)

// 每次心跳的时候agent把hostname汇报上来，经常要知道这个机器的hostid，把此信息缓存
// key: hostname value: hostid
type SafeHostMap struct {
	sync.RWMutex
	M map[string]int
}

var HostMap = &SafeHostMap{M: make(map[string]int)}

func (this *SafeHostMap) GetID(hostname string) (int, bool) {
	this.RLock()
	defer this.RUnlock()
	id, exists := this.M[hostname]
	return id, exists
}

func (this *SafeHostMap) Init() {
	m, err := db.QueryHosts()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

type SafeMonitoredHosts struct {
	sync.RWMutex
	M map[int]*model.Host
}

var MonitoredHosts = &SafeMonitoredHosts{M: make(map[int]*model.Host)}

func (this *SafeMonitoredHosts) Get() map[int]*model.Host {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeMonitoredHosts) Init() {
	m, err := db.QueryMonitoredHosts()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
