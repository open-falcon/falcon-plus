package store

import (
	"container/list"
	"hash/crc32"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/graph/g"
	//"log"
	"sync"
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
	idx := int(hashKey(key)) % this.Size
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
	idx := int(hashKey(key)) % this.Size
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

func (this *GraphItemMap) LenOf(key string) int {
	this.RLock()
	defer this.RUnlock()

	idx := int(hashKey(key)) % this.Size
	L, ok := this.A[idx][key]
	if !ok {
		return 0
	}
	return L.Len()
}

func (this *GraphItemMap) First(key string) *cmodel.GraphItem {
	this.RLock()
	defer this.RUnlock()
	idx := int(hashKey(key)) % this.Size
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
	idx := int(hashKey(key)) % this.Size
	L, ok := this.A[idx][key]
	if !ok {
		return []*cmodel.GraphItem{}
	}
	return L.PopAll()
}

func (this *GraphItemMap) FetchAll(key string) []*cmodel.GraphItem {
	this.RLock()
	defer this.RUnlock()
	idx := int(hashKey(key)) % this.Size
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

func (this *GraphItemMap) PushFront(key string, val *cmodel.GraphItem) {
	if linkedList, exists := this.Get(key); exists {
		linkedList.PushFront(val)
	} else {
		//log.Println("new key:", key)
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
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

/*
func (this *GraphItemMap) Keys() []string {
	this.RLock()
	defer this.RUnlock()

	count := len(this.M)
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range this.M {
		keys = append(keys, key)
	}
	return keys
}
*/

func init() {
	size := g.CACHE_TIME / g.FLUSH_DISK_STEP
	GraphItems = &GraphItemMap{
		A:    make([]map[string]*SafeLinkedList, size),
		Size: size,
	}
	for i := 0; i < size; i++ {
		GraphItems.A[i] = make(map[string]*SafeLinkedList)
	}
}
