package cache

import (
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
)

type SafeEExpCache struct {
	sync.RWMutex
	L []*model.EExp
}

var EExpCache = &SafeEExpCache{}

func (this *SafeEExpCache) Get() []*model.EExp {
	this.RLock()
	defer this.RUnlock()
	return this.L
}

func (this *SafeEExpCache) Init() {
	es, err := db.QueryEExps()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.L = es
}
