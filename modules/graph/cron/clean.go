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

// Copyright 2016 Xiaomi, Inc.
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

package cron

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/store"

	pfc "github.com/niean/goperfcounter"
)

func CleanCache() {

	ticker := time.NewTicker(time.Duration(g.CLEAN_CACHE) * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		DeleteInvalidItems()   // 删除无效的GraphItems
		DeleteInvalidHistory() // 删除无效的HistoryCache
	}
}

/*

  概念定义及结构体简谱:
  ckey = md5_type_step
  uuid = endpoint/metric/tags/dstype/step
  md5  = md5(endpoint/metric/tags)

  GraphItems        [idx]  [ckey] [{timestamp, value}, {timestamp, value} ...]
  HistoryCache      [md5]  [itemFirst, itemSecond]
  IndexedItemCache  [md5]  {UUID, Item}

*/

// TODO: 删除长期不更新数据(依赖index)
func DeleteInvalidItems() int {

	var currentCnt, deleteCnt int
	graphItems := store.GraphItems

	for idx := 0; idx < graphItems.Size; idx++ {
		keys := graphItems.KeysByIndex(idx)

		for _, key := range keys {
			tmp := strings.Split(key, "_") // key = md5_type_step
			if len(tmp) == 3 && !index.IndexedItemCache.ContainsKey(tmp[0]) {
				graphItems.Remove(key)
				deleteCnt++
			}
		}
	}
	currentCnt = graphItems.Len()

	pfc.Gauge("GraphItemsCacheCnt", int64(currentCnt))
	pfc.Gauge("GraphItemsCacheInvalidCnt", int64(deleteCnt))
	log.Infof("GraphItemsCache: Count=>%d, DeleteInvalid=>%d", currentCnt, deleteCnt)

	return deleteCnt
}

// TODO: 删除长期不更新数据(依赖index)
func DeleteInvalidHistory() int {

	var currentCnt, deleteCnt int
	historyCache := store.HistoryCache

	keys := historyCache.Keys()
	for _, key := range keys {
		if !index.IndexedItemCache.ContainsKey(key) {
			historyCache.Remove(key)
			deleteCnt++
		}
	}
	currentCnt = historyCache.Size()

	pfc.Gauge("HistoryCacheCnt", int64(currentCnt))
	pfc.Gauge("HistoryCacheInvalidCnt", int64(deleteCnt))
	log.Infof("HistoryCache: Count=>%d, DeleteInvalid=>%d", currentCnt, deleteCnt)

	return deleteCnt
}
