package cache

import (
	"github.com/open-falcon/hbs/db"
	"sort"
	"sync"
)

// 一个HostGroup可以绑定多个Plugin
type SafeGroupPlugins struct {
	sync.RWMutex
	M map[int][]string
}

var GroupPlugins = &SafeGroupPlugins{M: make(map[int][]string)}

func (this *SafeGroupPlugins) GetPlugins(gid int) ([]string, bool) {
	this.RLock()
	defer this.RUnlock()
	plugins, exists := this.M[gid]
	return plugins, exists
}

func (this *SafeGroupPlugins) Init() {
	m, err := db.QueryPlugins()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

// 根据hostname获取关联的插件
func GetPlugins(hostname string) []string {
	hid, exists := HostMap.GetID(hostname)
	if !exists {
		return []string{}
	}

	gids, exists := HostGroupsMap.GetGroupIds(hid)
	if !exists {
		return []string{}
	}

	// 因为机器关联了多个Group，每个Group可能关联多个Plugin，故而一个机器关联的Plugin可能重复
	pluginDirs := make(map[string]struct{})
	for _, gid := range gids {
		plugins, exists := GroupPlugins.GetPlugins(gid)
		if !exists {
			continue
		}

		for _, plugin := range plugins {
			pluginDirs[plugin] = struct{}{}
		}
	}

	size := len(pluginDirs)
	if size == 0 {
		return []string{}
	}

	dirs := make([]string, size)
	i := 0
	for dir := range pluginDirs {
		dirs[i] = dir
		i++
	}

	sort.Strings(dirs)
	return dirs
}
