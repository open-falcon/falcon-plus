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

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

type SafeLinkedList struct {
	sync.RWMutex
	Flag uint32
	L    *list.List
}

// 新创建SafeLinkedList容器
func NewSafeLinkedList() *SafeLinkedList {
	return &SafeLinkedList{L: list.New()}
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

func (this *SafeLinkedList) PopBack() *list.Element {
	this.Lock()
	defer this.Unlock()

	back := this.L.Back()
	if back != nil {
		this.L.Remove(back)
	}

	return back
}

func (this *SafeLinkedList) Back() *list.Element {
	this.Lock()
	defer this.Unlock()

	return this.L.Back()
}

func (this *SafeLinkedList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.L.Len()
}

// remain参数表示要给linkedlist中留几个元素
// 在cron中刷磁盘的时候要留一个，用于创建数据库索引
// 在程序退出的时候要一个不留的全部刷到磁盘
func (this *SafeLinkedList) PopAll() []*cmodel.GraphItem {
	this.Lock()
	defer this.Unlock()

	size := this.L.Len()
	if size <= 0 {
		return []*cmodel.GraphItem{}
	}

	ret := make([]*cmodel.GraphItem, 0, size)

	for i := 0; i < size; i++ {
		item := this.L.Back()
		ret = append(ret, item.Value.(*cmodel.GraphItem))
		this.L.Remove(item)
	}

	return ret
}

//restore PushAll
func (this *SafeLinkedList) PushAll(items []*cmodel.GraphItem) {
	this.Lock()
	defer this.Unlock()

	size := len(items)
	if size > 0 {
		for i := size - 1; i >= 0; i-- {
			this.L.PushBack(items[i])
		}
	}
}

//return为倒叙的?
func (this *SafeLinkedList) FetchAll() ([]*cmodel.GraphItem, uint32) {
	this.Lock()
	defer this.Unlock()
	count := this.L.Len()
	ret := make([]*cmodel.GraphItem, 0, count)

	p := this.L.Back()
	for p != nil {
		ret = append(ret, p.Value.(*cmodel.GraphItem))
		p = p.Prev()
	}

	return ret, this.Flag
}
