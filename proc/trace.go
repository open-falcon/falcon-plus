package proc

import (
	"container/list"
	"github.com/open-falcon/common/model"
	MUtils "github.com/open-falcon/common/utils"
	"sync"
)

const (
	DefaultMaxTraceSize = 3
)

var (
	RecvDataTrace = NewDataTrace("RecvDataTrace", DefaultMaxTraceSize)
)

// DataCacheQueue
type DataTrace struct {
	sync.RWMutex
	// config
	Endpoint string
	Metric   string
	Tags     map[string]string
	PK       string
	// data
	Name    string
	MaxSize int
	L       *list.List
}

func NewDataTrace(name string, maxSize int) *DataTrace {
	return &DataTrace{L: list.New(), Name: name, MaxSize: maxSize}
}

func (this *DataTrace) SetTraceConfig(endpoint string, metric string, tags map[string]string) {
	this.Lock()
	defer this.Unlock()

	if endpoint == "" || metric == "" {
		return
	}

	this.Endpoint = endpoint
	this.Metric = metric
	this.Tags = tags
	this.PK = MUtils.PK(endpoint, metric, tags)
}

func (this *DataTrace) PushFront(v interface{}) {
	this.Lock()
	defer this.Unlock()

	this.L.PushFront(v)
	if this.L.Len() > this.MaxSize {
		this.L.Remove(this.L.Back())
	}
}

func (this *DataTrace) FilterAll() []*list.Element {
	this.RLock()
	defer this.RUnlock()

	items := make([]*list.Element, 0)
	for e := this.L.Front(); e != nil; e = e.Next() {
		item := e.Value.(*model.MetaData)
		if this.PK == item.PK() {
			items = append(items, e)
		}
	}

	return items
}
