// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"container/list"
	"github.com/open-falcon/falcon-plus/common/model"

	"sync"
	"time"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func (this *SafeLinkedList) ToSlice() []*model.JudgeItem {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*model.JudgeItem{}
	}

	ret := make([]*model.JudgeItem, 0, sz)
	for e := this.L.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*model.JudgeItem))
	}
	return ret
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeLinkedList) HistoryData(limit int) ([]*model.HistoryData, bool) {
	if limit < 1 {
		// 其实limit不合法，此处也返回false吧，上层代码要注意
		// 因为false通常使上层代码进入异常分支，这样就统一了
		return []*model.HistoryData{}, false
	}

	size := this.Len()
	if size == 0 {
		return []*model.HistoryData{}, false
	}

	firstElement := this.Front()
	firstItem := firstElement.Value.(*model.JudgeItem)

	var vs []*model.HistoryData
	isEnough := true

	judgeType := firstItem.JudgeType[0]
	if judgeType == 'G' || judgeType == 'g' {
		if size < limit {
			// 有多少获取多少
			limit = size
			isEnough = false
		}
		vs = make([]*model.HistoryData, limit)
		vs[0] = &model.HistoryData{Timestamp: firstItem.Timestamp, Value: firstItem.Value}
		i := 1
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			vs[i] = &model.HistoryData{
				Timestamp: nextElement.Value.(*model.JudgeItem).Timestamp,
				Value:     nextElement.Value.(*model.JudgeItem).Value,
			}
			i++
			currentElement = nextElement
		}
	} else {
		if size < limit+1 {
			isEnough = false
			limit = size - 1
		}

		vs = make([]*model.HistoryData, limit)

		i := 0
		currentElement := firstElement
		for i < limit {
			nextElement := currentElement.Next()
			diffVal := currentElement.Value.(*model.JudgeItem).Value - nextElement.Value.(*model.JudgeItem).Value
			diffTs := currentElement.Value.(*model.JudgeItem).Timestamp - nextElement.Value.(*model.JudgeItem).Timestamp
			vs[i] = &model.HistoryData{
				Timestamp: currentElement.Value.(*model.JudgeItem).Timestamp,
				Value:     diffVal / float64(diffTs),
			}
			i++
			currentElement = nextElement
		}
	}

	return vs, isEnough
}

// @param limit 至多返回这些，如果不够，有多少返回多少
// @return bool isEnough
func (this *SafeLinkedList) HistoryDataString(pattern string, period int) ([]*model.HistoryData, bool) {
	size := this.Len()
	if size == 0 {
		return []*model.HistoryData{}, false
	}

	now := time.Now().Unix()
	then := now - int64(period)

	maxItems := 512
	var vs []*model.HistoryData
	hits := 0

	for e := this.Front(); e != nil && hits < maxItems; e = e.Next() {
		item := e.Value.(*model.JudgeItem)

		if item.Timestamp >= then {
			vs = append(vs, &model.HistoryData{Timestamp: item.Timestamp, Value: item.Value, ValueRaw: item.ValueRaw})
			hits += 1
		}
	}
	return vs, hits > 0
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
}

// @return needJudge 如果是false不需要做judge，因为新上来的数据不合法
func (this *SafeLinkedList) PushFrontAndMaintain(v *model.JudgeItem, maxCount int) bool {
	this.Lock()
	defer this.Unlock()

	sz := this.L.Len()
	if sz > 0 {
		// 新push上来的数据有可能重复了，或者timestamp不对，这种数据要丢掉
		if v.Timestamp <= this.L.Front().Value.(*model.JudgeItem).Timestamp || v.Timestamp <= 0 {
			return false
		}
	}

	this.L.PushFront(v)

	sz++
	if sz <= maxCount {
		return true
	}

	del := sz - maxCount
	for i := 0; i < del; i++ {
		this.L.Remove(this.L.Back())
	}

	return true
}

func (this *SafeLinkedList) Front() *list.Element {
	this.RLock()
	defer this.RUnlock()
	return this.L.Front()
}

func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}
