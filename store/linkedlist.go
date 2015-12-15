package store

import (
	"container/list"
	"sync"

	cmodel "github.com/open-falcon/common/model"
)

const (
	GRAPH_F_MISS     = 0x01
	GRAPH_F_FETCHING = 0x02
	GRAPH_F_ERR      = 0x04
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

//return为倒叙的?
func (this *SafeLinkedList) FetchAll() []*cmodel.GraphItem {
	this.Lock()
	defer this.Unlock()
	count := this.L.Len()
	ret := make([]*cmodel.GraphItem, 0, count)

	p := this.L.Back()
	for p != nil {
		ret = append(ret, p.Value.(*cmodel.GraphItem))
		p = p.Prev()
	}

	return ret
}
