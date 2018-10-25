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

package config

import (
	log "github.com/Sirupsen/logrus"
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/container/nmap"

	"github.com/open-falcon/falcon-plus/modules/nodata/config/service"
	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

// nodata配置(mockcfg)的缓存, 这些数据来自配置中心
var (
	rwlock      = sync.RWMutex{}
	NdConfigMap = nmap.NewSafeMap()
)

func Start() {
	if !g.Config().Config.Enabled {
		log.Println("config.Start warning, not enabled")
		return
	}

	service.InitDB()
	StartNdConfigCron()
	log.Println("config.Start ok")
}

// Interfaces Of StrategyMap
func SetNdConfigMap(val *nmap.SafeMap) {
	rwlock.Lock()
	defer rwlock.Unlock()

	NdConfigMap = val
}

func Keys() []string {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return NdConfigMap.Keys()
}

func Size() int {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return NdConfigMap.Size()
}

func GetNdConfig(key string) (*cmodel.NodataConfig, bool) {
	rwlock.RLock()
	defer rwlock.RUnlock()

	val, found := NdConfigMap.Get(key)
	if found && val != nil {
		return val.(*cmodel.NodataConfig), true
	}
	return &cmodel.NodataConfig{}, false
}
