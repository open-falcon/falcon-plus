package g

import (
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
)

type SafeEExpMap struct {
	sync.RWMutex
	M map[string][]model.EExp
}

type SafeEFilterMap struct {
	sync.RWMutex
	M map[string]string
}

var (
	EExpMap    = &SafeEExpMap{M: make(map[string][]model.EExp)}
	EFilterMap = &SafeEFilterMap{M: make(map[string]string)}
)

func (this *SafeEExpMap) ReInit(m map[string][]model.EExp) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeEExpMap) Get() map[string][]model.EExp {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeEFilterMap) ReInit(m map[string]string) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeEFilterMap) Exists(key string) bool {
	this.RLock()
	defer this.RUnlock()
	if _, ok := this.M[key]; ok {
		return true
	}
	return false
}
