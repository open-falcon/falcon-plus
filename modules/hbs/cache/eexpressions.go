package cache

import (
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeEExpressionCache struct {
	sync.RWMutex
	L []*model.EExpression
}

var EExpressionCache = &SafeEExpressionCache{}

func (this *SafeEExpressionCache) Get() []*model.EExpression {
	this.RLock()
	defer this.RUnlock()
	return this.L
}

func (this *SafeEExpressionCache) Init() {
	es, err := db.QueryEExpressions()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.L = es
}
