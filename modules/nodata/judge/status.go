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

package judge

import (
	"fmt"
	"sync"

	"github.com/toolkits/container/nmap"
	ttime "github.com/toolkits/time"
)

// Nodata Status Cache
var (
	statusLock = sync.RWMutex{}
	StatusMap  = nmap.NewSafeMap()
)

func LastTs(key string) int64 {
	statusLock.RLock()

	var ts int64 = 0
	v, found := StatusMap.Get(key)
	if !found {
		statusLock.RUnlock()
		return ts
	}

	ts = v.(*NodataStatus).Ts

	statusLock.RUnlock()
	return ts
}

func TurnOk(key string, ts int64) {
	statusLock.Lock()

	v, found := StatusMap.Get(key)
	if !found {
		// create new status
		ns := NewNodataStatus(key, "OK", 0, ts)
		StatusMap.Put(key, ns)
		statusLock.Unlock()
		return
	}

	// update status
	ns := v.(*NodataStatus)
	ns.Status = "OK"
	ns.Cnt = 0
	ns.Ts = ts

	statusLock.Unlock()
	return
}

func TurnNodata(key string, ts int64) {
	statusLock.Lock()

	v, found := StatusMap.Get(key)
	if !found {
		// create new status
		ns := NewNodataStatus(key, "NODATA", 1, ts)
		StatusMap.Put(key, ns)
		statusLock.Unlock()
		return
	}

	// update status
	ns := v.(*NodataStatus)
	ns.Status = "NODATA"
	ns.Cnt += 1
	ns.Ts = ts

	statusLock.Unlock()
	return
}

func GetNodataStatus(key string) *NodataStatus {
	statusLock.RLock()
	defer statusLock.RUnlock()

	v, found := StatusMap.Get(key)
	if !found {
		return &NodataStatus{}
	}
	return v.(*NodataStatus)
}

func GetAllNodataStatus() []*NodataStatus {
	statusLock.RLock()
	defer statusLock.RUnlock()

	ret := make([]*NodataStatus, 0)
	keys := StatusMap.Keys()
	for _, key := range keys {
		if v, found := StatusMap.Get(key); found {
			ret = append(ret, v.(*NodataStatus))
		}
	}

	return ret
}

// Nodata Status Struct
type NodataStatus struct {
	Key    string
	Status string // OK|NODATA
	Cnt    int
	Ts     int64
}

func NewNodataStatus(key string, status string, cnt int, ts int64) *NodataStatus {
	return &NodataStatus{key, status, cnt, ts}
}

func (this *NodataStatus) String() string {
	return fmt.Sprintf("NodataStatus key=%s status=%s cnt=%d ts=%d date=%s",
		this.Key, this.Status, this.Cnt, this.Ts, ttime.FormatTs(this.Ts))
}
