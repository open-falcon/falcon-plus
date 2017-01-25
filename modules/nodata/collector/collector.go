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

func AddItem(key string, val *DataItem) {
	listv, found := ItemMap.Get(key)
	if !found {
		ll := tlist.NewSafeListLimited(3) //每个采集指标,缓存最新的3个数据点
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
