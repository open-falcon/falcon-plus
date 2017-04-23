package cron

import (
	"time"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

func CleanCache() {

	ticker := time.NewTicker(time.Duration(g.CLEAN_CACHE) * time.Second)
	for {
		<-ticker.C
		store.DeleteInvalidItems()   // 删除无效的GraphItems
		store.DeleteInvalidHistory() // 删除无效的historyCache
	}
}

/*

  概念定义及结构体简谱:
  ckey = md5_type_step
  uuid = endpoint/metric/tags/dstype/step
  md5  = md5(endpoint/metric/tags)

  GraphItems        [idx]  [ckey] [{timestamp, value}, {timestamp, value} ...]
  MetaCache         [ckey] item
  history           [md5]  [itemFirst, itemSecond]
  IndexedItemCache  [md5]  {UUID, Item}

*/
