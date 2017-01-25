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
