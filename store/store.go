package store

import (
	"container/list"
	"hash/crc32"
	"log"
	"sync"
	"time"

	cmodel "github.com/open-falcon/common/model"

	"github.com/open-falcon/graph/g"
)

var GraphItems *GraphItemMap

type fetch_rrd struct {
	md5    string
	dstype string
	step   int
	sl     *SafeLinkedList
}

type GraphItemMap struct {
	sync.RWMutex
	A          []map[string]*SafeLinkedList
	Size       int
	fetch_list map[string]*SafeLinkedList
}

func (this *GraphItemMap) Get(key string) (*SafeLinkedList, bool) {
	this.RLock()
	defer this.RUnlock()
	idx := hashKey(key) % uint32(this.Size)
	val, ok := this.A[idx][key]
	return val, ok
}

func (this *GraphItemMap) Getitems(idx int) map[string]*SafeLinkedList {
	this.RLock()
	defer this.RUnlock()
	items := this.A[idx]
	this.A[idx] = make(map[string]*SafeLinkedList)
	return items
}

func (this *GraphItemMap) Set(key string, val *SafeLinkedList) {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	this.A[idx][key] = val
}

func (this *GraphItemMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	var l int
	for i := 0; i < this.Size; i++ {
		l += len(this.A[i])
	}
	return l
}

func (this *GraphItemMap) First(key string) *cmodel.GraphItem {
	this.RLock()
	defer this.RUnlock()
	idx := hashKey(key) % uint32(this.Size)
	L, ok := this.A[idx][key]
	if !ok {
		return nil
	}

	first := L.Front()
	if first == nil {
		return nil
	}

	return first.Value.(*cmodel.GraphItem)
}

func (this *GraphItemMap) PopAll(key string) []*cmodel.GraphItem {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	L, ok := this.A[idx][key]
	if !ok || (L.Flag & GRAPH_F_MISS) {
		return []*cmodel.GraphItem{}
	}
	return L.PopAll()
}

func (this *GraphItemMap) FetchAll(key string) []*cmodel.GraphItem {
	this.RLock()
	defer this.RUnlock()
	idx := hashKey(key) % uint32(this.Size)
	L, ok := this.A[idx][key]
	if !ok {
		return []*cmodel.GraphItem{}
	}

	return L.FetchAll()
}

func hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func getWts(key string, now int64) int64 {
	interval := int64(g.CACHE_TIME)
	return now + interval - (int64(hashKey(key)) % interval)
}

func (this *GraphItemMap) PushFront(key string,
	item *cmodel.GraphItem, md5 string, cfg *g.GlobalConfig) {
	if linkedList, exists := this.Get(key); exists {
		linkedList.PushFront(item)
	} else {
		//log.Println("new key:", key)
		safeList := &SafeLinkedList{L: list.New()}
		safeList.L.PushFront(item)

		if cfg.Migrate.Enabled && !g.IsRrdFileExist(g.RrdFileName(
			cfg.RRD.Storage, md5, item.DsType, item.Step)) {
			safeList.Flag = GRAPH_F_MISS
			f := &fetch_rrd{md5, dstype, step, sl}
			this.fetch_list.PushFront(f)
		}
		this.Set(key, safeList)
	}
}

func (this *GraphItemMap) KeysByIndex(idx int) []string {
	this.RLock()
	defer this.RUnlock()

	count := len(this.A[idx])
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range this.A[idx] {
		keys = append(keys, key)
	}

	return keys
}

func fetch_rrd(queue *SafeLinkedList, node, addr string) {
	for {
		entry := queue.PopBack()
		if entry == nil {
			time.Sleep(time.Second)
			continue
		}
		//todo: call jsonrpc to fetch rrd file
		time.Sleep(time.Second)
	}
}

func (this *GraphItemMap) Migrate_start() {
	cfg = g.Config()
	if cfg.Migrate.Enabled {
		for node, addr := range cfg.Migrate.Cluster {
			this.fetch_list[node] = &SafeLinkedList{L: list.New()}
			go fetch_rrd(this.fetch_list[node], node, addr)
		}
	}
}

func init() {
	size := g.CACHE_TIME / g.FLUSH_DISK_STEP
	if size < 0 {
		log.Panicf("store.init, bad size %d\n", size)
	}

	GraphItems = &GraphItemMap{
		A:          make([]map[string]*SafeLinkedList, size),
		Size:       size,
		fetch_list: make(map[string]*SafeLinkedList),
	}
	for i := 0; i < size; i++ {
		GraphItems.A[i] = make(map[string]*SafeLinkedList)
	}

}
