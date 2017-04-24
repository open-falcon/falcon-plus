package index

import (
	"fmt"
	"log"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

// 初始化索引功能模块
func Start() {
	InitCache()
	go StartIndexUpdateIncrTask()
	log.Println("index.Start ok")
}

// index收到一条新上报的监控数据,尝试用于增量更新索引
func ReceiveItem(item *cmodel.GraphItem, md5 string) {
	if item == nil {
		return
	}

	uuid := item.UUID()

	// 已上报过的数据
	if IndexedItemCache.ContainsKey(md5) {
		old := IndexedItemCache.Get(md5).(*IndexCacheItem)
		if uuid == old.UUID { // dsType+step没有发生变化,只更新缓存 TODO 存在线程安全的问题
			old.Item = item
		} else { // dsType+step变化了,当成一个新的增量来处理(甚至,不用rrd文件来过滤)
			//IndexedItemCache.Remove(md5)
			unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		}
		return
	}

	// 是否有rrdtool文件存在,如果有 认为已建立索引
	// 针对 索引缓存重建场景 做的优化, 结合索引全量更新 来保证一致性
	rrdFileName := g.RrdFileName(g.Config().RRD.Storage, md5, item.DsType, item.Step)
	if g.IsRrdFileExist(rrdFileName) {
		IndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		return
	}

	// 缓存未命中, 放入增量更新队列
	unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
}

//
func GetIndexedItemCache(endpoint string, metric string, tags map[string]string, dstype string, step int) (r *cmodel.GraphItem, rerr error) {
	itemDemo := &cmodel.GraphItem{
		Endpoint: endpoint,
		Metric:   metric,
		Tags:     tags,
		DsType:   dstype,
		Step:     step,
	}
	md5 := itemDemo.Checksum()
	uuid := itemDemo.UUID()

	cached := IndexedItemCache.Get(md5)
	if cached == nil {
		rerr = fmt.Errorf("not found")
		return
	}

	icitem := cached.(*IndexCacheItem)
	if icitem.UUID != uuid {
		rerr = fmt.Errorf("counter found, uuid not found: bad step or type")
		return
	}

	r = icitem.Item
	return
}
