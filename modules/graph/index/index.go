package index

import (
	log "github.com/Sirupsen/logrus"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/toolkits/file"
)

// 初始化索引功能模块
func Start() {
	InitCache()
	go StartIndexUpdateIncrTask()
	log.Debug("index.Start ok")
}

// index收到一条新上报的监控数据,尝试用于增量更新索引
func ReceiveItem(item *cmodel.GraphItem, md5 string) {
	if item == nil {
		return
	}

	uuid := item.UUID()

	// 已上报过的数据
	if indexedItemCache.ContainsKey(md5) {
		old := indexedItemCache.Get(md5).(*IndexCacheItem)
		if uuid == old.UUID { // dsType+step没有发生变化,只更新缓存
			indexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		} else { // dsType+step变化了,当成一个新的增量来处理
			unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		}
		return
	}

	// 针对 索引缓存重建场景 做的优化, 结合索引全量更新 来保证一致性
	// 是否有rrdtool文件存在,如果有 认为已建立索引
	rrdFileName := g.RrdFileName(g.Config().RRD.Storage, md5, item.DsType, item.Step)
	if g.IsRrdFileExist(rrdFileName) {
		indexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
		return
	}

	// 缓存未命中, 放入增量更新队列
	unIndexedItemCache.Put(md5, NewIndexCacheItem(uuid, item))
}

//从graph cache中删除掉某个item, 并删除指定的counter对应的rrd文件
func RemoveItem(item *cmodel.GraphItem) {
	md5 := item.Checksum()
	indexedItemCache.Remove(md5)

	rrdFileName := g.RrdFileName(g.Config().RRD.Storage, md5, item.DsType, item.Step)
	file.Remove(rrdFileName)
}
