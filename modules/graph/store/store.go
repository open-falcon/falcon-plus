package store

import (
	"container/list"
	"errors"
	"hash/crc32"
	"log"
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

var GraphItems *GraphItemMap

type GraphItemMap struct {
	sync.RWMutex
	A    []map[string]*SafeLinkedList
	Size int
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
	sl, ok := this.A[idx][key]
	if !ok {
		return nil
	}

	first := sl.Front()
	if first == nil {
		return nil
	}

	return first.Value.(*cmodel.GraphItem)
}

func (this *GraphItemMap) PushAll(key string, items []*cmodel.GraphItem) error {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	sl, ok := this.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.PushAll(items)
	return nil
}

func (this *GraphItemMap) GetFlag(key string) (uint32, error) {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	sl, ok := this.A[idx][key]
	if !ok {
		return 0, errors.New("not exist")
	}
	return sl.Flag, nil
}

func (this *GraphItemMap) SetFlag(key string, flag uint32) error {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	sl, ok := this.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.Flag = flag
	return nil
}

func (this *GraphItemMap) PopAll(key string) []*cmodel.GraphItem {
	this.Lock()
	defer this.Unlock()
	idx := hashKey(key) % uint32(this.Size)
	sl, ok := this.A[idx][key]
	if !ok {
		return []*cmodel.GraphItem{}
	}
	return sl.PopAll()
}

func (this *GraphItemMap) FetchAll(key string) ([]*cmodel.GraphItem, uint32) {
	this.RLock()
	defer this.RUnlock()
	idx := hashKey(key) % uint32(this.Size)
	sl, ok := this.A[idx][key]
	if !ok {
		return []*cmodel.GraphItem{}, 0
	}

	return sl.FetchAll()
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
			safeList.Flag = g.GRAPH_F_MISS
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

func init() {
	size := g.CACHE_TIME / g.FLUSH_DISK_STEP
	if size < 0 {
		log.Panicf("store.init, bad size %d\n", size)
	}

	GraphItems = &GraphItemMap{
		A:    make([]map[string]*SafeLinkedList, size),
		Size: size,
	}
	for i := 0; i < size; i++ {
		GraphItems.A[i] = make(map[string]*SafeLinkedList)
	}
}
