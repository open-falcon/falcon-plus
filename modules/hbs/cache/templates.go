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

package cache

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"sync"
)

// 一个HostGroup对应多个Template
type SafeGroupTemplates struct {
	sync.RWMutex
	M map[int][]int
}

var GroupTemplates = &SafeGroupTemplates{M: make(map[int][]int)}

func (this *SafeGroupTemplates) GetTemplateIds(gid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	templateIds, exists := this.M[gid]
	return templateIds, exists
}

func (this *SafeGroupTemplates) Init() {
	m, err := db.QueryGroupTemplates()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

type SafeTemplateCache struct {
	sync.RWMutex
	M map[int]*model.Template
}

var TemplateCache = &SafeTemplateCache{M: make(map[int]*model.Template)}

func (this *SafeTemplateCache) GetMap() map[int]*model.Template {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeTemplateCache) Init() {
	ts, err := db.QueryTemplates()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = ts
}

type SafeHostTemplateIds struct {
	sync.RWMutex
	M map[int][]int
}

var HostTemplateIds = &SafeHostTemplateIds{M: make(map[int][]int)}

func (this *SafeHostTemplateIds) GetMap() map[int][]int {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeHostTemplateIds) Init() {
	m, err := db.QueryHostTemplateIds()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
