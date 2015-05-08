package store

import (
	"sync"
)

// endpoint表的内存缓存
// key: endpoint
// val: id
type Endpoint2Id struct {
	sync.RWMutex
	M map[string]int64
}

func (this *Endpoint2Id) EndpointId(endpoint string) (int64, bool) {
	this.RLock()
	defer this.RUnlock()
	id, exists := this.M[endpoint]
	return id, exists
}

func (this *Endpoint2Id) Set(endpoint string, id int64) {
	this.Lock()
	defer this.Unlock()
	this.M[endpoint] = id
}

// tag_endpoint表的内存缓存
// key:
// val:
type TagEndpointMap struct {
	sync.RWMutex
	M map[string]struct{}
}

func (this *TagEndpointMap) Exists(key string) bool {
	this.RLock()
	defer this.RUnlock()

	_, exists := this.M[key]
	return exists
}

func (this *TagEndpointMap) Set(key string) {
	this.Lock()
	defer this.Unlock()

	this.M[key] = struct{}{}
}

// 这是为endpoint_counter表准备的内存缓存
// key: endpointId-counter
// val: step-dstype
type CounterMap struct {
	sync.RWMutex
	M map[string]string
}

func (this *CounterMap) ExistsInMemory(key, nv string) bool {
	this.RLock()
	defer this.RUnlock()
	ov, exists := this.M[key]
	if !exists {
		return false
	}

	return ov == nv
}

func (this *CounterMap) Get(key string) (string, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[key]
	return val, exists
}

func (this *CounterMap) Set(key, val string) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

var (
	Endpoint2Ids   = &Endpoint2Id{M: make(map[string]int64)}
	TagEndpointSet = &TagEndpointMap{M: make(map[string]struct{})}
	Counters       = &CounterMap{M: make(map[string]string)}
)
