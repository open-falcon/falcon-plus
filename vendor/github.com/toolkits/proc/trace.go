package proc

import (
	"container/list"
	"sync"
)

type DataTrace struct {
	sync.RWMutex
	MaxSize int
	Name    string
	PK      string
	L       *list.List
}

func NewDataTrace(name string, maxSize int) *DataTrace {
	return &DataTrace{L: list.New(), Name: name, MaxSize: maxSize}
}

func (this *DataTrace) SetPK(pk string) {
	this.Lock()
	defer this.Unlock()

	// rm old caches when trace's pk changed
	if this.PK != pk {
		this.L = list.New()
	}
	this.PK = pk
}

// proposed that there were few traced items
func (this *DataTrace) Trace(pk string, v interface{}) {
	this.RLock()
	if this.PK != pk {
		this.RUnlock()
		return
	}
	this.RUnlock()

	// we could almost not step here, so we get few wlock
	this.Lock()
	defer this.Unlock()
	this.L.PushFront(v)
	if this.L.Len() > this.MaxSize {
		this.L.Remove(this.L.Back())
	}
}

func (this *DataTrace) GetAllTraced() []interface{} {
	this.RLock()
	defer this.RUnlock()

	items := make([]interface{}, 0)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e)
	}

	return items
}
