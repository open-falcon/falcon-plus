package proc

import (
	"container/list"
	"fmt"
	"math"
	"sync"
)

type DataFilter struct {
	sync.RWMutex
	MaxSize   int
	Name      string
	PK        string
	Opt       string
	Threshold float64
	L         *list.List
}

func NewDataFilter(name string, maxSize int) *DataFilter {
	return &DataFilter{L: list.New(), Name: name, MaxSize: maxSize}
}

func (this *DataFilter) SetFilter(pk string, opt string, threshhold float64) error {
	this.Lock()
	defer this.Unlock()

	if !legalOpt(opt) {
		return fmt.Errorf("bad opt: %s", opt)
	}

	// rm old caches when filter's pk changed
	if this.PK != pk {
		this.L = list.New()
	}
	this.PK = pk
	this.Opt = opt
	this.Threshold = threshhold

	return nil
}

// proposed that there were few traced items
func (this *DataFilter) Filter(pk string, val float64, v interface{}) {
	this.RLock()
	if this.PK != pk {
		this.RUnlock()
		return
	}
	this.RUnlock()

	// we could almost not step here, so we get few wlock
	this.Lock()
	defer this.Unlock()
	if compute(this.Opt, val, this.Threshold) {
		this.L.PushFront(v)
		if this.L.Len() > this.MaxSize {
			this.L.Remove(this.L.Back())
		}
	}
}

func (this *DataFilter) GetAllFiltered() []interface{} {
	this.RLock()
	defer this.RUnlock()

	items := make([]interface{}, 0)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e)
	}

	return items
}

// internal
const (
	MinPositiveFloat64 = 0.000001
	MaxNegativeFloat64 = -0.000001
)

func compute(opt string, left float64, right float64) bool {
	switch opt {
	case "eq":
		return math.Abs(left-right) < MinPositiveFloat64
	case "ne":
		return math.Abs(left-right) >= MinPositiveFloat64
	case "gt":
		return (left - right) > MinPositiveFloat64
	case "lt":
		return (left - right) < MaxNegativeFloat64
	default:
		return false
	}
}

func legalOpt(opt string) bool {
	switch opt {
	case "eq", "ne", "gt", "lt":
		return true
	default:
		return false
	}
}
