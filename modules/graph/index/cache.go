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

package index

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"

	tcache "github.com/toolkits/cache/localcache/timedcache"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/proc"
)

const (
	DefaultMaxCacheSize                     = 5000000 // 默认 最多500w个,太大了内存会耗尽
	DefaultCacheProcUpdateTaskSleepInterval = time.Duration(1) * time.Second
)

// item缓存
var (
	IndexedItemCache   = NewIndexCacheBase(DefaultMaxCacheSize)
	unIndexedItemCache = NewIndexCacheBase(DefaultMaxCacheSize)
)

// db本地缓存
var (
	// endpoint表的内存缓存, key:endpoint(string) / value:id(int64)
	dbEndpointCache = tcache.New(600*time.Second, 60*time.Second)
	// endpoint_counter表的内存缓存, key:endpoint_id-counter(string) / val:dstype-step(string)
	dbEndpointCounterCache = tcache.New(600*time.Second, 60*time.Second)
)

// 初始化cache
func InitCache() {
	go startCacheProcUpdateTask()
}

// USED WHEN QUERY
func GetTypeAndStep(endpoint string, counter string) (dsType string, step int, found bool) {
	// get it from index cache
	pk := cutils.Md5(fmt.Sprintf("%s/%s", endpoint, counter))
	if icitem := IndexedItemCache.Get(pk); icitem != nil {
		if item := icitem.(*IndexCacheItem).Item; item != nil {
			dsType = item.DsType
			step = item.Step
			found = true
			return
		}
	}

	// statistics
	proc.GraphLoadDbCnt.Incr()

	// get it from db, this should rarely happen
	var endpointId int64 = -1
	if endpointId, found = GetEndpointFromCache(endpoint); found {
		if dsType, step, found = GetCounterFromCache(endpointId, counter); found {
			//found = true
			return
		}
	}

	// do not find it, this must be a bad request
	found = false
	return
}

// Return EndpointId if Found
func GetEndpointFromCache(endpoint string) (int64, bool) {
	// get from cache
	endpointId, found := dbEndpointCache.Get(endpoint)
	if found {
		return endpointId.(int64), true
	}

	// get from db
	var id int64 = -1
	err := g.DB.QueryRow("SELECT id FROM endpoint WHERE endpoint = ?", endpoint).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Println("query endpoint id fail,", err)
		return -1, false
	}

	if err == sql.ErrNoRows || id < 0 {
		return -1, false
	}

	// update cache
	dbEndpointCache.Set(endpoint, id, 0)

	return id, true
}

// Return DsType Step if Found
func GetCounterFromCache(endpointId int64, counter string) (dsType string, step int, found bool) {
	var err error
	// get from cache
	key := fmt.Sprintf("%d-%s", endpointId, counter)
	dsTypeStep, found := dbEndpointCounterCache.Get(key)
	if found {
		arr := strings.Split(dsTypeStep.(string), "_")
		step, err = strconv.Atoi(arr[1])
		if err != nil {
			found = false
			return
		}
		dsType = arr[0]
		return
	}

	// get from db
	err = g.DB.QueryRow("SELECT type, step FROM endpoint_counter WHERE endpoint_id = ? and counter = ?",
		endpointId, counter).Scan(&dsType, &step)
	if err != nil && err != sql.ErrNoRows {
		log.Println("query type and step fail", err)
		return
	}

	if err == sql.ErrNoRows {
		return
	}

	// update cache
	dbEndpointCounterCache.Set(key, fmt.Sprintf("%s_%d", dsType, step), 0)

	found = true
	return
}

// 更新 cache的统计信息
func startCacheProcUpdateTask() {
	for {
		time.Sleep(DefaultCacheProcUpdateTaskSleepInterval)
		proc.IndexedItemCacheCnt.SetCnt(int64(IndexedItemCache.Size()))
		proc.UnIndexedItemCacheCnt.SetCnt(int64(unIndexedItemCache.Size()))
		proc.EndpointCacheCnt.SetCnt(int64(dbEndpointCache.Size()))
		proc.CounterCacheCnt.SetCnt(int64(dbEndpointCounterCache.Size()))
	}
}

// INDEX CACHE
// 索引缓存的元素数据结构
type IndexCacheItem struct {
	UUID string
	Item *cmodel.GraphItem
}

func NewIndexCacheItem(uuid string, item *cmodel.GraphItem) *IndexCacheItem {
	return &IndexCacheItem{UUID: uuid, Item: item}
}

// 索引缓存-基本缓存容器
type IndexCacheBase struct {
	sync.RWMutex
	maxSize int
	data    map[string]interface{}
}

func NewIndexCacheBase(max int) *IndexCacheBase {
	return &IndexCacheBase{maxSize: max, data: make(map[string]interface{})}
}

func (this *IndexCacheBase) GetMaxSize() int {
	return this.maxSize
}

func (this *IndexCacheBase) Put(key string, item interface{}) {
	this.Lock()
	defer this.Unlock()
	this.data[key] = item
}

func (this *IndexCacheBase) Remove(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.data, key)
}

func (this *IndexCacheBase) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()
	return this.data[key]
}

func (this *IndexCacheBase) ContainsKey(key string) bool {
	this.RLock()
	defer this.RUnlock()
	return this.data[key] != nil
}

func (this *IndexCacheBase) Size() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.data)
}

func (this *IndexCacheBase) Keys() []string {
	this.RLock()
	defer this.RUnlock()

	count := len(this.data)
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range this.data {
		keys = append(keys, key)
	}

	return keys
}
