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
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
)

type JudgeItemMap struct {
	sync.RWMutex
	M map[string]*SafeLinkedList
}

func NewJudgeItemMap() *JudgeItemMap {
	return &JudgeItemMap{M: make(map[string]*SafeLinkedList)}
}

func (this *JudgeItemMap) Get(key string) (*SafeLinkedList, bool) {
	this.RLock()
	defer this.RUnlock()
	val, ok := this.M[key]
	return val, ok
}

func (this *JudgeItemMap) Set(key string, val *SafeLinkedList) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *JudgeItemMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *JudgeItemMap) Delete(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

func (this *JudgeItemMap) BatchDelete(keys []string) {
	count := len(keys)
	if count == 0 {
		return
	}

	this.Lock()
	defer this.Unlock()
	for i := 0; i < count; i++ {
		delete(this.M, keys[i])
	}
}

func (this *JudgeItemMap) CleanStale(before int64) {
	keys := []string{}

	this.RLock()
	for key, L := range this.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*model.JudgeItem).Timestamp < before {
			keys = append(keys, key)
		}
	}
	this.RUnlock()

	this.BatchDelete(keys)
}

func (this *JudgeItemMap) PushFrontAndMaintain(key string, val *model.JudgeItem, maxCount int, now int64) {
	validStrategies, needJudgeByStrategy := CheckStrategy(val)
	validExpressions, needJudgeByExpression := CheckExpression(val)

	if !needJudgeByStrategy && !needJudgeByExpression {
		return
	}

	if linkedList, exists := this.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			Judge(linkedList, val, now, validStrategies, validExpressions)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeLinkedList{L: NL}
		this.Set(key, safeList)
		Judge(safeList, val, now, validStrategies, validExpressions)
	}
}

// 这是个线程不安全的大Map，需要提前初始化好
var HistoryBigMap = make(map[string]*JudgeItemMap)

func InitHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			HistoryBigMap[arr[i]+arr[j]] = NewJudgeItemMap()
		}
	}
}
