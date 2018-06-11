package proc

import (
	"sync"
	"time"

	ntime "github.com/toolkits/time"
)

const (
	DefaultOtherMaxSize      = 100
	DefaultSCounterQpsPeriod = 1
)

// basic counter
type SCounterBase struct {
	sync.RWMutex
	Name  string
	Cnt   int64
	Time  string
	ts    int64
	Other map[string]interface{}
}

func NewSCounterBase(name string) *SCounterBase {
	uts := time.Now().Unix()
	return &SCounterBase{Name: name, Cnt: 0, Time: ntime.FormatTs(uts),
		ts: uts, Other: make(map[string]interface{})}
}

func (this *SCounterBase) Get() *SCounterBase {
	this.RLock()
	defer this.RUnlock()

	return &SCounterBase{
		Name:  this.Name,
		Cnt:   this.Cnt,
		Time:  this.Time,
		ts:    this.ts,
		Other: deepCopyMap(this.Other),
	}
}

func (this *SCounterBase) SetCnt(cnt int64) {
	this.Lock()
	this.Cnt = cnt
	this.ts = time.Now().Unix()
	this.Time = ntime.FormatTs(this.ts)
	this.Unlock()
}

func (this *SCounterBase) Incr() {
	this.IncrBy(int64(1))
}

func (this *SCounterBase) IncrBy(incr int64) {
	this.Lock()
	this.Cnt += incr
	this.Unlock()
}

func (this *SCounterBase) PutOther(key string, value interface{}) bool {
	this.Lock()
	defer this.Unlock()

	ret := false
	_, exist := this.Other[key]
	if exist {
		this.Other[key] = value
		ret = true
	} else {
		if len(this.Other) < DefaultOtherMaxSize {
			this.Other[key] = value
			ret = true
		}
	}

	return ret
}

// counter with qps
type SCounterQps struct {
	sync.RWMutex
	Name    string
	Cnt     int64
	Qps     int64
	Time    string
	ts      int64
	lastTs  int64
	lastCnt int64
	Other   map[string]interface{}
}

func NewSCounterQps(name string) *SCounterQps {
	uts := time.Now().Unix()
	return &SCounterQps{Name: name, Cnt: 0, Time: ntime.FormatTs(uts), ts: uts,
		Qps: 0, lastCnt: 0, lastTs: uts, Other: make(map[string]interface{})}
}

func (this *SCounterQps) Get() *SCounterQps {
	this.Lock()
	defer this.Unlock()

	this.ts = time.Now().Unix()
	this.Time = ntime.FormatTs(this.ts)
	// get smooth qps value
	if this.ts-this.lastTs > DefaultSCounterQpsPeriod {
		this.Qps = int64((this.Cnt - this.lastCnt) / (this.ts - this.lastTs))
		this.lastTs = this.ts
		this.lastCnt = this.Cnt
	}

	return &SCounterQps{
		Name:    this.Name,
		Cnt:     this.Cnt,
		Qps:     this.Qps,
		Time:    this.Time,
		ts:      this.ts,
		lastTs:  this.lastTs,
		lastCnt: this.lastCnt,
		Other:   deepCopyMap(this.Other),
	}
}

func (this *SCounterQps) Incr() {
	this.IncrBy(int64(1))
}

func (this *SCounterQps) IncrBy(incr int64) {
	this.Lock()
	this.incrBy(incr)
	this.Unlock()
}

func (this *SCounterQps) PutOther(key string, value interface{}) bool {
	this.Lock()
	defer this.Unlock()

	ret := false
	_, exist := this.Other[key]
	if exist {
		this.Other[key] = value
		ret = true
	} else {
		if len(this.Other) < DefaultOtherMaxSize {
			this.Other[key] = value
			ret = true
		}
	}

	return ret
}

func (this *SCounterQps) incrBy(incr int64) {
	this.Cnt += incr
}

func deepCopyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for key, val := range src {
		dst[key] = val
	}
	return dst
}
