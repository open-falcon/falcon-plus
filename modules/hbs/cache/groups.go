package cache

import (
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"sync"
)

// 一个机器可能在多个group下，做一个map缓存hostid与groupid的对应关系
type SafeHostGroupsMap struct {
	sync.RWMutex
	M map[int][]int
}

var HostGroupsMap = &SafeHostGroupsMap{M: make(map[int][]int)}

func (this *SafeHostGroupsMap) GetGroupIds(hid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	gids, exists := this.M[hid]
	return gids, exists
}

func (this *SafeHostGroupsMap) Init() {
	m, err := db.QueryHostGroups()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
