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

package sender

import (
	"container/list"
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	L *list.List
}

func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
}

func (this *SafeLinkedList) PopBack(limit int) []*model.JsonMetaData {
	this.RLock()
	defer this.RUnlock()
	sz := this.L.Len()
	if sz == 0 {
		return []*model.JsonMetaData{}
	}

	if sz < limit {
		limit = sz
	}

	ret := make([]*model.JsonMetaData, 0, limit)
	for i := 0; i < limit; i++ {
		e := this.L.Back()
		ret = append(ret, e.Value.(*model.JsonMetaData))
		this.L.Remove(e)
	}

	return ret
}

func (this *SafeLinkedList) PushFront(v interface{}) *list.Element {
	this.Lock()
	defer this.Unlock()
	return this.L.PushFront(v)
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
