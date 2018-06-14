package model

import (
	"container/list"
	"sync"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

type EJudgeItemMap struct {
	sync.RWMutex
	M map[string]*SafeELinkedList
}

func NewEJudgeItemMap() *EJudgeItemMap {
	return &EJudgeItemMap{M: make(map[string]*SafeELinkedList)}
}

func (this *EJudgeItemMap) Get(key string) (*SafeELinkedList, bool) {
	this.RLock()
	defer this.RUnlock()
	val, ok := this.M[key]
	return val, ok
}

func (this *EJudgeItemMap) Set(key string, val *SafeELinkedList) {
	this.Lock()
	defer this.Unlock()
	this.M[key] = val
}

func (this *EJudgeItemMap) Len() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.M)
}

func (this *EJudgeItemMap) Delete(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.M, key)
}

func (this *EJudgeItemMap) BatchDelete(keys []string) {
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

func (this *EJudgeItemMap) CleanStale(before int64) {
	keys := []string{}

	this.RLock()
	for key, L := range this.M {
		front := L.Front()
		if front == nil {
			continue
		}

		if front.Value.(*cmodel.EMetric).Timestamp < before {
			keys = append(keys, key)
		}
	}
	this.RUnlock()

	this.BatchDelete(keys)
}

func (this *EJudgeItemMap) PushFrontAndMaintain(key string, val *cmodel.EMetric, maxCount int, now int64) {
	if linkedList, exists := this.Get(key); exists {
		needJudge := linkedList.PushFrontAndMaintain(val, maxCount)
		if needJudge {
			EJudge(linkedList, val, now)
		}
	} else {
		NL := list.New()
		NL.PushFront(val)
		safeList := &SafeELinkedList{L: NL}
		this.Set(key, safeList)
		EJudge(safeList, val, now)
	}
}

// 这是个线程不安全的大Map，需要提前初始化好
var EHistoryBigMap = make(map[string]*EJudgeItemMap)

func InitEHistoryBigMap() {
	arr := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			EHistoryBigMap[arr[i]+arr[j]] = NewEJudgeItemMap()
		}
	}
}
