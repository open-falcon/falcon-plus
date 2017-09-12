// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
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
