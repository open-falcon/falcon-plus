package g

import (
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
)

type SafeEExpressionMap struct {
	sync.RWMutex
	// join(filtersK=V1,",") => [exp1, exp2 ...]
	// join(filtersK=V2,",") => [exp1, exp2 ...]
	M map[string][]*model.EExpression
}

type SafeEFilterMap struct {
	sync.RWMutex
	M map[string]string
}

var (
	EExpressionMap = &SafeEExpressionMap{M: make(map[string][]*model.EExpression)}
	EFilterMap     = &SafeEFilterMap{M: make(map[string]string)}
)

func (this *SafeEExpressionMap) ReInit(m map[string][]*model.EExpression) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeEExpressionMap) Get() map[string][]*model.EExpression {
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
