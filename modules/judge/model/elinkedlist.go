package model

import (
	"container/list"
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

type SafeELinkedList struct {
	sync.RWMutex
	L *list.List
}

func (this *SafeELinkedList) ToSlice() []*cmodel.EMetric {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*cmodel.EMetric{}
	}

	ret := make([]*cmodel.EMetric, 0, sz)
	for e := this.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*cmodel.EMetric))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeELinkedList) HistoryData(limit uint64) ([]*cmodel.EHistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*cmodel.EHistoryData{}, false
	}

	size := uint64(this.Len())
	if size == 0 {
		return []*cmodel.EHistoryData{}, false
	}

	firstElement := this.Front()
	firstItem := firstElement.Value.(*cmodel.EMetric)

	var vs []*cmodel.EHistoryData
	isEnough := true

	if size < limit {
		// 有多少获取多少
		limit = size
		isEnough = false
	}
	vs = make([]*cmodel.EHistoryData, limit)
	vs[0] = &cmodel.EHistoryData{Timestamp: firstItem.Timestamp, Filters: firstItem.Filters}
	var i uint64
	currentElement := firstElement
	for i = 1; i < limit; {
		nextElement := currentElement.Next()
		vs[i] = &cmodel.EHistoryData{
			Timestamp: nextElement.Value.(*cmodel.EMetric).Timestamp,
			Filters:   nextElement.Value.(*cmodel.EMetric).Filters,
		}
		i++
		currentElement = nextElement
	}

	return vs, isEnough
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeELinkedList) HistoryDataByTime(ago uint64, hits uint64, now uint64) ([]*cmodel.EHistoryData, bool) {
	var isEnough bool
	var vs []*cmodel.EHistoryData

	if hits < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return vs, isEnough
	}

	size := uint64(this.Len())
	if size == 0 {
		return vs, isEnough
	}

	if size < hits {
		return vs, isEnough
	}

	t := int64(now - ago)

	matched := uint64(0)
	for e := this.Front(); e != nil && matched < hits; e = e.Next() {
		item := e.Value.(*cmodel.EMetric)
		if item.Timestamp >= t {
			item := &cmodel.EHistoryData{Timestamp: item.Timestamp, Filters: item.Filters}
			vs = append(vs, item)
			matched++
		}
	}

	isEnough = matched >= hits

	return vs, isEnough
}

func (this *SafeELinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

// @return needJudge 如果是false不需要做judge，因为新上来的数据不合法
func (this *SafeELinkedList) PushFrontAndMaintain(v *cmodel.EMetric, maxCount uint32) bool {
	this.Lock()
	defer this.Unlock()

	sz := uint32(this.L.Len())
	this.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	var i uint32
	for i = 0; i < del; i++ {
		this.L.Remove(this.L.Back())
	}

	return true
}

func (this *SafeELinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

func (this *SafeELinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}
