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

package collector

import (
	"fmt"
	"log"

	tlist "github.com/toolkits/container/list"
	"github.com/toolkits/container/nmap"
	ttime "github.com/toolkits/time"

	"github.com/open-falcon/falcon-plus/modules/nodata/g"
)

// 主动收集到的监控数据 的缓存
var (
	// map - list
	ItemMap = nmap.NewSafeMap()
)

func Start() {
	if !g.Config().Collector.Enabled {
		log.Println("collector.Start warning, not enabled")
		return
	}

	StartCollectorCron()
	log.Println("collector.Start ok")
}

// Interfaces Of ItemMap
func GetFirstItem(key string) (*DataItem, bool) {
	listv, found := ItemMap.Get(key)
	if !found || listv == nil {
		return &DataItem{}, false
	}

	first := listv.(*tlist.SafeListLimited).Front()
	if first == nil {
		return &DataItem{}, false
	}

	return first.(*DataItem), true
}

func GetItemByIndex(key string, index int) (*DataItem, bool) {
	listv, found := ItemMap.Get(key)
	if !found || listv == nil {
		return &DataItem{}, false
	}

	all := listv.(*tlist.SafeListLimited).FrontAll()
	if all == nil || len(all) <= index {
		return &DataItem{}, false
	}

	return all[index].(*DataItem), true
}

func AddItem(key string, val *DataItem) {
	listv, found := ItemMap.Get(key)
	if !found {
		ll := tlist.NewSafeListLimited(10) //每个采集指标,缓存最新的3个数据点
		ll.PushFrontViolently(val)
		ItemMap.Put(key, ll)
		return
	}

	listv.(*tlist.SafeListLimited).PushFrontViolently(val)
}

func RemoveItem(key string) {
	ItemMap.Remove(key)
}

// NoData Data Item Struct
type DataItem struct {
	Ts      int64
	Value   float64
	FStatus string // OK|ERR
	FTs     int64
}

func NewDataItem(ts int64, val float64, fstatus string, fts int64) *DataItem {
	return &DataItem{Ts: ts, Value: val, FStatus: fstatus, FTs: fts}
}

func (this *DataItem) String() string {
	return fmt.Sprintf("ts:%s, value:%f, fts:%s, fstatus:%s",
		ttime.FormatTs(this.Ts), this.Value, ttime.FormatTs(this.FTs), this.FStatus)
}
