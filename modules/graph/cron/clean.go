package cron

import (
	"log"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/index"
	"github.com/open-falcon/falcon-plus/modules/graph/store"

	pfc "github.com/niean/goperfcounter"
)

func CleanCache() {

	ticker := time.NewTicker(time.Duration(g.CLEAN_CACHE) * time.Second)
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

	graphItems := store.GraphItems

	deleteCnt := 0
	for idx := 0; idx < graphItems.Size; idx++ {
		keys := graphItems.KeysByIndex(idx)

		deleteKeys := make([]string, 0)
		for _, key := range keys {
			tmp := strings.Split(key, "_") // key = md5_type_step
			if len(tmp) == 3 && !index.IndexedItemCache.ContainsKey(tmp[0]) {
				deleteKeys = append(deleteKeys, key)
			}
		}

		graphItems.Lock()
		for _, key := range deleteKeys {
			delete(graphItems.A[idx], key)
		}
		graphItems.Unlock()
		deleteCnt += len(deleteKeys)
	}

	pfc.Gauge("GraphItemsCacheCnt", int64(graphItems.Len()))
	pfc.Gauge("GraphItemsCacheInvalidCnt", int64(deleteCnt))
	log.Printf("GraphItemsCache: Count=>%d, DeleteInvalid=>%d", graphItems.Len(), deleteCnt)

	return deleteCnt
}

// TODO: 删除长期不更新数据(依赖index)
func DeleteInvalidHistory() int {

	historyCache := store.HistoryCache

	deleteKeys := make([]string, 0)
	historyCache.RLock()
	for key, _ := range historyCache.M {
		if !index.IndexedItemCache.ContainsKey(key) {
			deleteKeys = append(deleteKeys, key)
		}
	}
	historyCache.RUnlock()

	historyCache.Lock()
	for _, key := range deleteKeys {
		delete(historyCache.M, key)
	}

	pfc.Gauge("HistoryCacheCnt", int64(len(historyCache.M)))
	pfc.Gauge("HistoryCacheInvalidCnt", int64(len(deleteKeys)))
	log.Printf("HistoryCache: Count=>%d, DeleteInvalid=>%d", len(historyCache.M), len(deleteKeys))

	historyCache.Unlock()
	return len(deleteKeys)
}
