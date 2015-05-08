package proc

import (
	"sync"
	"time"
)

const (
	DefaultOtherMaxSize      = 100 // 默认的 自定义统计字段的最大值
	DefaultSCounterQpsPeriod = 1   // QPS计算周期, 默认值为1s
)

// 基本的计数器
type SCounterBase struct {
	sync.RWMutex
	Name string
	Cnt  int64
	Time string
	ts   int64
	// 自定义统计指标
	Other map[string]interface{}
}

func NewSCounterBase(name string) *SCounterBase {
	uts := time.Now().Unix()
	return &SCounterBase{Name: name, Cnt: 0, Time: formatTs(uts), ts: uts,
		Other: make(map[string]interface{})}
}

func (this *SCounterBase) Get() *SCounterBase {
	this.RLock()
	defer this.RUnlock()
	return this
}

// to be abandoned
func (this *SCounterBase) Set(cnt int64) {
	this.Lock()
	defer this.Unlock()

	this.Cnt = cnt
	this.ts = time.Now().Unix()
	this.Time = formatTs(this.ts)
}

func (this *SCounterBase) SetCnt(cnt int64) {
	this.Lock()
	defer this.Unlock()

	this.Cnt = cnt
	this.ts = time.Now().Unix()
	this.Time = formatTs(this.ts)
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

func formatTs(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

// QPS统计,只支持 增加计数操作
type SCounterQps struct {
	sync.RWMutex
	Name string
	Cnt  int64
	Qps  int64
	Time string
	ts   int64
	// for qps
	lastTs  int64
	lastCnt int64
	// 自定义统计指标
	Other map[string]interface{}
}

func NewSCounterQps(name string) *SCounterQps {
	uts := time.Now().Unix()
	return &SCounterQps{Name: name, Cnt: 0, Time: formatTs(uts), ts: uts,
		Qps: 0, lastCnt: 0, lastTs: uts,
		Other: make(map[string]interface{})}
}

func (this *SCounterQps) Get() *SCounterQps {
	this.RLock()
	defer this.RUnlock()

	this.incrByDirty(0) // update qps
	return this
}

func (this *SCounterQps) Incr() {
	this.IncrBy(int64(1))
}

func (this *SCounterQps) IncrBy(incr int64) {
	this.Lock()
	defer this.Unlock()

	this.incrByDirty(incr)
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

// 操作我的时候,请加写锁
func (this *SCounterQps) incrByDirty(incr int64) {
	this.Cnt += incr
	this.ts = time.Now().Unix()
	this.Time = formatTs(this.ts)

	// qps
	if this.ts-this.lastTs > DefaultSCounterQpsPeriod {
		this.Qps = int64((this.Cnt - this.lastCnt) / (this.ts - this.lastTs))
		this.lastTs = this.ts
		this.lastCnt = this.Cnt
	}
}
