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

package g

import (
	"fmt"
	"sync"
	"time"
)

type Cluster struct {
	Id          int64
	GroupId     int64
	Numerator   string
	Denominator string
	Endpoint    string
	Metric      string
	Tags        string
	DsType      string
	Step        int
	LastUpdate  time.Time
}

func (this *Cluster) String() string {
	return fmt.Sprintf(
		"<Id:%d, GroupId:%d, Numerator:%s, Denominator:%s, Endpoint:%s, Metric:%s, Tags:%s, DsType:%s, Step:%d, LastUpdate:%v>",
		this.Id,
		this.GroupId,
		this.Numerator,
		this.Denominator,
		this.Endpoint,
		this.Metric,
		this.Tags,
		this.DsType,
		this.Step,
		this.LastUpdate,
	)
}

// key: Id+LastUpdate
type SafeClusterMonitorItems struct {
	sync.RWMutex
	M map[string]*Cluster
}

func NewClusterMonitorItems() *SafeClusterMonitorItems {
	return &SafeClusterMonitorItems{M: make(map[string]*Cluster)}
}

func (this *SafeClusterMonitorItems) Init(m map[string]*Cluster) {
	this.Lock()
	defer this.Unlock()
	this.M = m
}

func (this *SafeClusterMonitorItems) Get() map[string]*Cluster {
	this.RLock()
	defer this.RUnlock()
	return this.M
}
